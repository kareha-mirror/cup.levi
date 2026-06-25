package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"

	"tea.kareha.org/cup/termi"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/colors"
	"tea.kareha.org/cup/levi/internal/rkind"
)

func getConfigPath(dir string) string {
	return filepath.Join(dir, "editor.yaml")
}

type Mode int

const (
	ModeCommand Mode = iota
	ModeInsert
	ModeSearch
	ModePrompt
)

type ViewMeta struct {
	Loc buf.Loc
}

type Editor struct {
	dir      string
	cfg      *Config
	w, h     int
	bufs     []*buf.Buf
	bIndex   int
	inp      *Input
	inpRow   int // 0-based
	inserted []string
	mode     Mode
	alive    bool
	message  string
	ring     string
	parser   *Parser
	prompt   termi.RuneBuf
	regs     Regs
	backward bool
	pattern  termi.RuneBuf
	regexp   *regexp.Regexp
	lastCmd  Cmd
	redraw   bool
	view     []string
	vMeta    []ViewMeta
	listener termi.EscapeListener
	esc      bool
	colors   *colors.Colors
}

func (ed *Editor) Clear() {
	if ed.bIndex < len(ed.bufs) {
		ed.bufs[ed.bIndex] = new(buf.Buf)
	} else {
		ed.bufs = append(ed.bufs, new(buf.Buf))
	}
	ed.mode = ModeCommand
	ed.redraw = true
}

func (ed *Editor) Buf() *buf.Buf {
	return ed.bufs[ed.bIndex]
}

func (ed *Editor) Close(force bool) {
	b := ed.Buf()
	if !force && b.Modified {
		ed.Ring("File modified since last complete write; write or use ! to override.")
		return
	}
	bufs := append([]*buf.Buf{}, ed.bufs[:ed.bIndex]...)
	if ed.bIndex+1 < len(ed.bufs) {
		bufs = append(bufs, ed.bufs[ed.bIndex+1:]...)
	}
	ed.bufs = bufs
	n := len(ed.bufs)
	if ed.bIndex >= n {
		ed.bIndex = max(n-1, 0)
	}
	if n < 1 {
		ed.alive = false
	}
}

func (ed *Editor) Load(path string, force bool) error {
	b := ed.Buf()
	if !force && b.Modified {
		ed.Ring("File modified since last complete write; write or use ! to override.")
		return fmt.Errorf("file modified")
	}
	ed.Clear()
	b = ed.Buf()
	b.Path = path
	if path == "" {
		ed.Message("(memory): new file: line 1")
		return nil
	}
	info, err := os.Stat(path)
	if err != nil {
		ed.Message("%s: new file: line 1", path)
		return nil
	}
	stamp := buf.Stamp{
		Time: info.ModTime(),
		Size: info.Size(),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	b.SetText(string(data))
	b.Stamp = stamp
	b.Modified = false
	return nil
}

func (ed *Editor) InitialInfo() {
	b := ed.Buf()
	path := b.Path
	if path == "" {
		path = "(memory)"
	}
	modified := "unmodified"
	if b.Modified {
		modified = "modified"
	}
	info := "empty file"
	numLines := b.NumLines()
	if numLines > 0 {
		info = fmt.Sprintf("line %d", b.Loc.Row+1)
	}
	ed.Message("%s: %s: %s", path, modified, info)
}

func Init(dir string, args []string) (*Editor, error) {
	var cfg *Config
	cfgPath := getConfigPath(dir)
	_, err := os.Stat(cfgPath)
	if err != nil {
		cfg = DefaultConfig()
		SaveConfig(cfgPath, cfg)
	} else {
		cfg = LoadConfig(cfgPath)
	}

	list := colors.LoadList(dir)
	colors, _ := list.Load(cfg.Colors)

	w, h := termi.Size()
	ed := &Editor{
		dir:      dir,
		cfg:      cfg,
		w:        w,
		h:        h,
		bufs:     []*buf.Buf{},
		bIndex:   0,
		inp:      NewInput(),
		inpRow:   0,
		inserted: []string{},
		mode:     ModeCommand,
		alive:    true,
		message:  "",
		ring:     "",
		parser:   NewParser(),
		prompt:   termi.RuneBuf{},
		regs:     Regs{},
		backward: false,
		pattern:  termi.RuneBuf{},
		regexp:   nil,
		lastCmd:  Cmd{Kind: CmdInvalid},
		redraw:   true,
		view:     []string{},
		listener: nil,
		esc:      false,
		colors:   colors,
	}

	ed.regs.LoadConfig(ed.cfg)

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
		ed.Clear()
		ed.Load("", true)
	} else {
		for _, path := range args {
			ed.Clear()
			ed.Load(path, true)
			ed.bIndex++
		}
		ed.bIndex = 0
	}
	ed.InitialInfo()

	termi.SetEscapeListener(ed.listener)
	return ed, nil
}

func (ed *Editor) SaveAs(path string, force bool) error {
	if path == "" {
		ed.Ring("No filename specified")
		return fmt.Errorf("no filename specified")
	}
	info, err := os.Stat(path)
	newFile := ""
	stamp := buf.Stamp{}
	if err != nil {
		newFile = " new file:"
	} else {
		stamp = buf.Stamp{
			Time: info.ModTime(),
			Size: info.Size(),
		}
	}
	b := ed.Buf()
	if !force && path == b.Path && stamp != b.Stamp {
		ed.Ring(
			"%s: file modified more recently than this copy; use ! to override.",
			path,
		)
		return fmt.Errorf("file modified more recently")
	}

	text := b.Text()
	err = os.WriteFile(path, []byte(text), 0666)
	if err != nil {
		return err
	}
	info, err = os.Stat(path)
	if err != nil {
		return err
	}
	stamp = buf.Stamp{
		Time: info.ModTime(),
		Size: info.Size(),
	}

	ed.Message(
		"%s:%s %d lines, %d bytes, %d runes.",
		path, newFile, b.NumLines(), len(text), utf8.RuneCountInString(text),
	)

	if b.Path == "" {
		b.Path = path
	}
	if path == b.Path {
		b.Stamp = stamp
	}
	b.Modified = false
	return nil
}

func (ed *Editor) Save(force bool) error {
	b := ed.Buf()
	return ed.SaveAs(b.Path, force)
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
		} else if row < ed.inpRow+ed.inp.LineLen() {
			return ed.inp.Line(row - ed.inpRow)
		} else {
			return b.Line(row - ed.inp.LineLen() + 1)
		}
	}

	return b.Line(row)
}

func (ed *Editor) CurrentLine() string {
	b := ed.Buf()
	return ed.Line(b.Loc.Row)
}

func (ed *Editor) RuneCount() int {
	return utf8.RuneCountInString(ed.CurrentLine())
}

func (ed *Editor) EnsureCommand() {
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
		if ed.inpRow+1 <= b.NumLines()-1 {
			lines = append(lines, b.Lines[ed.inpRow+1:]...)
		}
		b.Lines = lines
		ed.inserted = ed.inp.Inserted()
		ed.inp.Reset()
		ed.mode = ModeCommand
		loc, ok := ed.MoveLeft(1)
		if ok {
			b.Loc = loc
		}
		b.Modified = true

		if MultiInsertCmds[ed.lastCmd.Kind] && ed.lastCmd.Num > 1 {
			cmd := ed.lastCmd
			cmd.Num--
			ed.Run(cmd, true)
		}

		ed.EndMemory()

		ed.parser.ClearAll()
		return
	case ModeSearch:
		ed.mode = ModeCommand
		return
	case ModePrompt:
		ed.mode = ModeCommand
		return
	}
}

func (ed *Editor) Message(format string, a ...any) {
	ed.message = fmt.Sprintf(format, a...)
}

func (ed *Editor) Ring(format string, a ...any) {
	ed.ring = fmt.Sprintf(format, a...)
}

func (ed *Editor) Error(format string, a ...any) {
	ed.Ring("Error: "+format, a...)
}

func (ed *Editor) Notice(format string, a ...any) {
	if ed.cfg.Silent {
		return
	}
	ed.Message("Notice: "+format, a...)
}

func (ed *Editor) Unimplemented(name string) {
	ed.Ring("not implemented (" + name + ")")
}
