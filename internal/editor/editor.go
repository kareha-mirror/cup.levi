package editor

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"tea.kareha.org/cup/termi"
)

type Mode int

const (
	ModeCommand Mode = iota
	ModeInsert
	ModeSearch
	ModePrompt
)

type Stamp struct {
	Time time.Time
	Size int64
}

type KillMode int

const (
	KillNone = iota
	KillRunes
	KillLines
)

type KillBuf struct {
	mode  KillMode
	runes []rune
	lines []string
}

func (k *KillBuf) SetRunes(runes []rune) {
	k.mode = KillRunes
	k.runes = append([]rune{}, runes...)
}

func (k *KillBuf) SetLines(lines []string) {
	k.mode = KillLines
	k.lines = append([]string{}, lines...)
}

type Editor struct {
	cfg      *Config
	col, row int // 0-based
	vrow     int // 0-based
	virtCol  int // 0-based
	w, h     int
	x, y     int // 0-based
	lines    []string
	inp      *Input
	inpRow   int // 0-based
	mode     Mode
	path     string
	stamp    Stamp
	alive    bool
	message  string
	ring     string
	parser   *Parser
	prompt   termi.RuneBuf
	modified bool
	killed   KillBuf
	redraw   bool
	view     []string
	listener termi.EscapeListener
	esc      bool
}

func (ed *Editor) Load(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	stamp := Stamp{
		Time: info.ModTime(),
		Size: info.Size(),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(data) < 1 {
		ed.lines = []string{}
		ed.modified = false
		return nil
	}
	// TODO should also support CRLF or not?
	if data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}
	ed.lines = strings.Split(string(data), "\n")
	ed.stamp = stamp
	ed.modified = false
	return nil
}

func Init(args []string) *Editor {
	var path string
	if len(args) > 1 {
		path = args[1]
	}

	w, h := termi.Size()
	ed := &Editor{
		cfg:      DefaultConfig(),
		col:      0,
		row:      0,
		vrow:     0,
		virtCol:  0,
		w:        w,
		h:        h,
		x:        0,
		y:        0,
		lines:    []string{},
		inp:      NewInput(),
		inpRow:   0,
		mode:     ModeCommand,
		path:     path,
		stamp:    Stamp{},
		alive:    true,
		message:  "",
		ring:     "",
		parser:   NewParser(),
		prompt:   termi.RuneBuf{},
		modified: false,
		killed:   KillBuf{},
		redraw:   true,
		view:     []string{},
		listener: nil,
		esc:      false,
	}

	if path != "" {
		_, err := os.Stat(path)
		if err == nil { // file exists
			ed.Load(path)
		}
	}

	termi.TabWidth = ed.cfg.TabWidth
	termi.Raw()
	termi.Init()

	listener := func(esc bool) {
		ed.esc = esc
		ed.DrawStatus()
		ed.PlaceCursor()
	}
	ed.listener = termi.EscapeListener(&listener)
	termi.AddEscapeListener(ed.listener)

	return ed
}

func (ed *Editor) SaveAs(path string, force bool) error {
	if path == "" {
		ed.Ring("No filename specified")
		return fmt.Errorf("no filename specified")
	}
	info, err := os.Stat(path)
	newFile := ""
	stamp := Stamp{}
	if err != nil {
		newFile = " new file:"
	} else {
		stamp = Stamp{
			Time: info.ModTime(),
			Size: info.Size(),
		}
	}
	if !force && path == ed.path && stamp != ed.stamp {
		ed.Ring(
			"%s: file modified more recently than this copy; use ! to override.",
			path,
		)
		return fmt.Errorf("file modified more recently")
	}

	text := ""
	if len(ed.lines) > 0 {
		text = strings.Join(ed.lines, "\n") + "\n"
	}
	err = os.WriteFile(path, []byte(text), 0666)
	if err != nil {
		return err
	}
	info, err = os.Stat(path)
	if err != nil {
		return err
	}
	stamp = Stamp{
		Time: info.ModTime(),
		Size: info.Size(),
	}

	ed.Message(
		"%s:%s %d lines, %d bytes, %d runes.",
		path, newFile, len(ed.lines), len(text), utf8.RuneCountInString(text),
	)

	if ed.path == "" {
		ed.path = path
	}
	if path == ed.path {
		ed.stamp = stamp
	}
	ed.modified = false
	return nil
}

func (ed *Editor) Save(force bool) error {
	return ed.SaveAs(ed.path, force)
}

func (ed *Editor) Finish() {
	termi.RemoveEscapeListener(ed.listener)

	termi.Finish()

	fmt.Print(termi.Clear)
	fmt.Print(termi.HomeCursor)
	termi.Cooked()
	fmt.Print(termi.ShowCursor)
}

func (ed *Editor) Line(row int) string {
	if ed.mode == ModeInsert {
		if row < ed.inpRow {
			return ed.lines[row]
		} else if row < ed.inpRow+ed.inp.LineLen() {
			return ed.inp.Line(row - ed.inpRow)
		} else {
			return ed.lines[row-ed.inp.LineLen()+1]
		}
	}

	if len(ed.lines) < 1 {
		return ""
	}
	return ed.lines[row]
}

func (ed *Editor) CurrentLine() string {
	return ed.Line(ed.row)
}

func (ed *Editor) RuneCount() int {
	return utf8.RuneCountInString(ed.CurrentLine())
}

func isBlankLine(s string) bool {
	for _, r := range s {
		if !isBlank(r) {
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
		lines := []string{}
		if ed.inpRow > 0 {
			lines = append(lines, ed.lines[:ed.inpRow]...)
		}
		inputLines := ed.inp.Lines()
		if ed.cfg.AutoIndent {
			for i := 0; i < len(inputLines); i++ {
				if isBlankLine(inputLines[i]) {
					inputLines[i] = ""
				}
			}
		}
		lines = append(lines, inputLines...)
		if ed.inpRow+1 <= len(ed.lines)-1 {
			lines = append(lines, ed.lines[ed.inpRow+1:]...)
		}
		ed.lines = lines
		ed.inp.Reset()
		ed.mode = ModeCommand
		ed.MoveLeft(1)

		ed.modified = true
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
	ed.ring = fmt.Sprintf("Error: "+format, a...)
}

func (ed *Editor) Unimplemented(name string) {
	ed.Ring("not implemented (" + name + ")")
}
