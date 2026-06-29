package editor

import (
	"fmt"
	"time"

	"tea.kareha.org/cup/termi"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/color"
	"tea.kareha.org/cup/levi/internal/rkind"
)

type Mode int

const (
	ModeCommand Mode = iota
	ModeInsert
	ModePrompt
	ModeSearch
)

type Editor struct {
	// config and state
	cfgDir string
	cfg    *Config
	bufs   []*buf.Buf
	bufIdx int
	mode   Mode
	alive  bool
	msg    *Msg

	// parser and input
	parser   Parser
	inp      Input
	inpRow   int // 0-based
	inserted []string
	prompt   termi.RuneBuf
	search   Search
	find     Find
	regs     Regs
	clipUsed bool
	lastCmd  CmdPair
	undo     bool

	// screen
	w, h     int
	redraw   bool
	view     []string
	viewMeta []ViewMeta
	colors   *color.Scheme

	// escape key indicator
	listener termi.EscapeListener
	esc      bool
}

func Init(cfgDir string, paths []string) (*Editor, error) {
	// for storing bootup errors
	msg := new(Msg)

	// load config file
	cfg, err := PrepareConfig(cfgDir)
	if err != nil {
		msg.Error("%v", err)
	}

	// load colorscheme
	cs, err := color.LoadScheme(cfgDir, cfg.Colors)
	if err != nil {
		msg.Error("%v", err)
	}

	// create and init editor struct
	ed := &Editor{
		cfgDir: cfgDir,
		cfg:    cfg,
		alive:  true,
		msg:    msg,

		redraw: true,
		colors: cs,
	}

	// render bootup errors
	ed.w, ed.h = termi.Size() // must be filled before rendered
	ed.RenderMsg(false)       // messages
	ed.RenderMsg(true)        // errors

	// setup shared registers
	ed.regs.SyncWithConfig(ed.cfg)

	// preferences
	termi.EscapeTimeout =
		time.Duration(ed.cfg.EscapeTimeout) * time.Millisecond
	termi.TabWidth = ed.cfg.TabStop

	// init terminal framework and terminal state
	termi.Raw()
	fmt.Print(termi.SetAlternate)
	err = termi.StartKey() // setup key handler
	if err != nil {
		fmt.Print(termi.ResetAlternate)
		termi.Cooked()
		return nil, err
	}
	termi.StartSig() // setup signal handler

	// init escape key indicator
	listener := func(esc bool) {
		ed.esc = esc
		ed.DrawStatus()
		ed.PlaceCursor()
	}
	ed.listener = termi.EscapeListener(&listener)
	termi.SetEscapeListener(ed.listener)

	// load files if supplied
	for _, path := range paths {
		ed.NewBuf()
		if !ed.Load(path, true) {
			ed.Close(true)
			continue
		}
		ed.bufIdx++
	}
	// select first buffer
	ed.bufIdx = 0
	// setup empty buffer if no files are loaded
	if ed.NumBufs() < 1 {
		ed.NewBuf()
		ed.Load("", true)
	}
	ed.ShowFileInfo()

	return ed, nil
}

func (ed *Editor) Finish() error {
	// shutdown terminal framework
	termi.SetEscapeListener(nil)
	termi.StopSig()
	err := termi.StopKey()

	// reset terminal mode
	fmt.Print(termi.Clear)
	fmt.Print(termi.HomeCursor)
	fmt.Print(termi.ResetAlternate)
	termi.Cooked()
	fmt.Print(termi.ShowCursor)

	return err
}

func (ed *Editor) Line(row int) string {
	b := ed.Buf()

	// input possibly has multiple lines
	if ed.mode == ModeInsert {
		if row < ed.inpRow {
			return b.Line(row)
		} else if row < ed.inpRow+ed.inp.NumLines() {
			return ed.inp.Line(row - ed.inpRow)
		} else {
			return b.Line(row - ed.inp.NumLines() + 1)
		}
	}

	return b.Line(row)
}

func (ed *Editor) Reset() {
	ed.parser.ResetAll()
}

func (ed *Editor) Commit() {
	switch ed.mode {

	case ModeCommand:
		// do nothing
		return

	// commit input text and ensure command mode
	case ModeInsert:
		b := ed.Buf()
		// head
		lines := append([]string{}, b.Lines[:ed.inpRow]...)

		// body
		inputLines := ed.inp.Lines()
		if ed.cfg.AutoIndent {
			for i := 0; i < len(inputLines); i++ {
				if rkind.IsBlankLine(inputLines[i]) {
					inputLines[i] = ""
				}
			}
		}
		lines = append(lines, inputLines...)

		// tail
		if ed.inpRow+1 < b.NumLines() {
			lines = append(lines, b.Lines[ed.inpRow+1:]...)
		}

		// replace slice
		b.Lines = lines

		ed.inserted = ed.inp.Inserted()
		ed.inp.Reset()

		ed.mode = ModeCommand

		// if number prefix is supplied, repeat insertion
		if _, ok :=
			MultiInsertCmds[ed.lastCmd.Main.Kind]; ok && ed.lastCmd.Main.Num > 1 {
			cmd := ed.lastCmd
			cmd.Main.Num--
			ed.Run(cmd, true) // replay
		} else {
			// or finish insertion
			b.Loc.Col--
			b.Loc = b.ConfineInclusive(b.Loc)
			b.VirtCol = b.Loc.Col
			b.Modified = true
		}

		ed.EndRecordForUndo()
		ed.Reset()

		if !ed.cfg.Silent {
			ed.ShowStatOfInserted()
		}
		return

	case ModeSearch:
		ed.mode = ModeCommand
		return

	case ModePrompt:
		ed.mode = ModeCommand
		return
	}
}
