package editor

import (
	"fmt"
	"time"
	"unicode/utf8"

	"tea.kareha.org/cup/termi"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/colors"
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
	dir    string
	cfg    *Config
	bufs   []*buf.Buf
	bufIdx int
	mode   Mode
	alive  bool
	msg    *Msg

	parser   Parser
	inp      Input
	inpRow   int // 0-based
	inserted []string
	prompt   termi.RuneBuf
	search   Search
	find     Find
	regs     Regs
	lastCmd  Cmd

	w, h     int
	redraw   bool
	view     []string
	viewMeta []ViewMeta
	colors   *colors.Colors

	listener termi.EscapeListener
	esc      bool
}

func Init(dir string, args []string) (*Editor, error) {
	msg := new(Msg)

	cfg, err := PrepareConfig(dir)
	if err != nil {
		msg.Error("%v", err)
	}

	var clrs *colors.Colors
	list, err := colors.LoadList(dir)
	if err != nil {
		msg.Error("%v", err)
	} else {
		clrs, err = list.Load(cfg.Colors)
		if err != nil {
			msg.Error("%v", err)
		}
	}

	ed := &Editor{
		dir:    dir,
		cfg:    cfg,
		bufs:   []*buf.Buf{},
		bufIdx: 0,
		mode:   ModeCommand,
		alive:  true,
		msg:    msg,

		parser:   Parser{},
		inp:      Input{},
		inpRow:   0,
		inserted: []string{},
		prompt:   termi.RuneBuf{},
		search:   Search{},
		find:     Find{},
		regs:     Regs{},
		lastCmd:  Cmd{Kind: CmdInvalid},

		w:        80,
		h:        24,
		redraw:   true,
		view:     nil,
		viewMeta: nil,
		colors:   clrs,

		listener: nil,
		esc:      false,
	}

	ed.regs.SyncWithConfig(ed.cfg)

	termi.EscapeTimeout = time.Duration(ed.cfg.EscTimeout) * time.Millisecond
	termi.TabWidth = ed.cfg.TabStop

	termi.Raw()
	fmt.Print(termi.SetAlternate)
	err = termi.StartKey()
	if err != nil {
		fmt.Print(termi.ResetAlternate)
		termi.Cooked()
		return nil, err
	}
	termi.StartSig()

	listener := func(esc bool) {
		ed.esc = esc
		ed.DrawStatus()
		ed.PlaceCursor()
	}
	ed.listener = termi.EscapeListener(&listener)

	if len(args) < 1 {
		ed.NewBuf()
		ed.Load("", true)
	} else {
		for _, path := range args {
			ed.NewBuf()
			ed.Load(path, true)
			ed.bufIdx++
		}
		ed.bufIdx = 0
	}
	ed.InitialInfo()

	termi.SetEscapeListener(ed.listener)
	return ed, nil
}

func (ed *Editor) Finish() error {
	termi.SetEscapeListener(nil)
	termi.StopSig()
	err := termi.StopKey()

	fmt.Print(termi.Clear)
	fmt.Print(termi.HomeCursor)
	fmt.Print(termi.ResetAlternate)
	termi.Cooked()
	fmt.Print(termi.ShowCursor)

	return err
}

func (ed *Editor) Line(row int) string {
	b := ed.Buf()

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
		return
	case ModeInsert:
		b := ed.Buf()
		lines := append([]string{}, b.Lines[:ed.inpRow]...)
		inputLines := ed.inp.Lines()
		if ed.cfg.AutoIndent {
			for i := 0; i < len(inputLines); i++ {
				if rkind.IsBlankLine(inputLines[i]) {
					inputLines[i] = ""
				}
			}
		}
		lines = append(lines, inputLines...)
		if ed.inpRow+1 < b.NumLines() {
			lines = append(lines, b.Lines[ed.inpRow+1:]...)
		}
		b.Lines = lines
		ed.inserted = ed.inp.Inserted()
		ed.inp.Reset()
		ed.mode = ModeCommand

		if _, ok := MultiInsertCmds[ed.lastCmd.Kind]; ok && ed.lastCmd.Num > 1 {
			cmd := ed.lastCmd
			cmd.Num--
			ed.Run(cmd, true) // replay
		} else {
			b.Loc.Col--
			b.Loc = b.ConfineInclusive(b.Loc)
			b.VirtCol = b.Loc.Col
			b.Modified = true
		}

		ed.EndMemory()
		ed.Reset()

		if !ed.cfg.Silent {
			numLines := len(ed.inserted)
			numBytes := numLines - 1
			numRunes := numLines - 1
			for _, line := range ed.inserted {
				numBytes += len(line)
				numRunes += utf8.RuneCountInString(line)
			}
			if numBytes > 0 {
				if numLines < 2 {
					ed.Notice(
						"%d bytes, %d runes inserted",
						numBytes, numRunes,
					)
				} else {
					ed.Notice(
						"%d lines, %d bytes, %d runes inserted",
						numLines, numBytes, numRunes,
					)
				}
			}
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
