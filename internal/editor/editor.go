package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"

	"tea.kareha.org/cup/termi"

	"tea.kareha.org/cup/levi/internal/buffer"
	"tea.kareha.org/cup/levi/internal/colors"
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

type KillMode int

const (
	KillNone = iota
	KillRunes
	KillLines
)

type KillBuf struct {
	mode  KillMode
	lines []string
}

func (k *KillBuf) SetRunes(lines []string) {
	k.mode = KillRunes
	k.lines = append([]string{}, lines...)
}

func (k *KillBuf) SetLines(lines []string) {
	k.mode = KillLines
	k.lines = append([]string{}, lines...)
}

type ViewMeta struct {
	Loc buffer.Loc
}

type Editor struct {
	dir      string
	cfg      *Config
	w, h     int
	buffers  []*buffer.Buffer
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
	killed   KillBuf
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
	if ed.bIndex < len(ed.buffers) {
		ed.buffers[ed.bIndex] = new(buffer.Buffer)
	} else {
		ed.buffers = append(ed.buffers, new(buffer.Buffer))
	}
	ed.mode = ModeCommand
	ed.redraw = true
}

func (ed *Editor) Buffer() *buffer.Buffer {
	return ed.buffers[ed.bIndex]
}

func (ed *Editor) Close(force bool) {
	b := ed.Buffer()
	if !force && b.Modified {
		ed.Ring("File modified since last complete write; write or use ! to override.")
		return
	}
	buffers := []*buffer.Buffer{}
	if ed.bIndex-1 > 0 {
		buffers = append(buffers, ed.buffers[:ed.bIndex-1]...)
	}
	if ed.bIndex+1 <= len(ed.buffers)-1 {
		buffers = append(buffers, ed.buffers[ed.bIndex+1:]...)
	}
	ed.buffers = buffers
	if len(ed.buffers) < 1 {
		ed.alive = false
	}
}

func (ed *Editor) Load(path string, force bool) error {
	b := ed.Buffer()
	if !force && b.Modified {
		ed.Ring("File modified since last complete write; write or use ! to override.")
		return fmt.Errorf("file modified")
	}
	ed.Clear()
	b = ed.Buffer()
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
	stamp := buffer.Stamp{
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
	b := ed.Buffer()
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
		buffers:  []*buffer.Buffer{},
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
		killed:   KillBuf{},
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
	stamp := buffer.Stamp{}
	if err != nil {
		newFile = " new file:"
	} else {
		stamp = buffer.Stamp{
			Time: info.ModTime(),
			Size: info.Size(),
		}
	}
	b := ed.Buffer()
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
	stamp = buffer.Stamp{
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
	b := ed.Buffer()
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
	b := ed.Buffer()

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
	b := ed.Buffer()
	return ed.Line(b.Loc.Row)
}

func (ed *Editor) RuneCount() int {
	return utf8.RuneCountInString(ed.CurrentLine())
}

func isBlankLine(s string) bool {
	for _, r := range s {
		if !isBlankRune(r) {
			return false
		}
	}
	return true
}

func (ed *Editor) EnsureCommand() {
	switch ed.mode {
	case ModeCommand:
		return
	case ModeInsert:
		b := ed.Buffer()
		lines := append([]string{}, b.Lines[:ed.inpRow]...)
		inputLines := ed.inp.Lines()
		if ed.cfg.AutoIndent {
			for i := 0; i < len(inputLines); i++ {
				if isBlankLine(inputLines[i]) {
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
		ed.MoveLeft(1)
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
