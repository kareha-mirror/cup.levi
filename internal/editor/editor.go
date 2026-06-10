package editor

import (
	"fmt"
	"os"
	"strings"
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

type Editor struct {
	col, row int
	vrow     int
	virtCol  int
	w, h     int
	x, y     int
	lines    []string
	inp      *Input
	inpRow   int
	mode     Mode
	path     string
	message  string
	parser   *Parser
	prompt   termi.RuneBuf
	save     bool
	quit     bool
	listener termi.EscapeListener
	esc      bool
}

func (ed *Editor) Load(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if len(data) < 1 {
		ed.lines = make([]string, 1)
	}
	// TODO CRLF
	if data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}
	ed.lines = strings.Split(string(data), "\n")
	ed.path = path
}

func Init(args []string) *Editor {
	var path string
	if len(args) > 1 {
		path = args[1]
	}

	w, h := termi.Size()
	ed := &Editor{
		col:      0,
		row:      0,
		vrow:     0,
		virtCol:  0,
		w:        w,
		h:        h,
		x:        0,
		y:        0,
		lines:    make([]string, 1),
		inp:      NewInput(),
		inpRow:   0,
		mode:     ModeCommand,
		path:     path,
		message:  "",
		parser:   NewParser(),
		prompt:   termi.RuneBuf{},
		save:     false,
		quit:     false,
		listener: nil,
		esc:      false,
	}

	if path != "" {
		_, err := os.Stat(path)
		if err == nil { // file exists
			ed.Load(path)
		}
	}

	termi.Raw()
	termi.Init()

	listener := func(esc bool) {
		ed.esc = esc
		ed.DrawStatus()
	}
	ed.listener = termi.EscapeListener(&listener)
	termi.AddEscapeListener(ed.listener)

	return ed
}

func (ed *Editor) Save(path string) {
	text := strings.Join(ed.lines, "\n") + "\n"
	err := os.WriteFile(path, []byte(text), 0666)
	if err != nil {
		panic(err)
	}
	ed.path = path
}

func (ed *Editor) Finish() {
	termi.RemoveEscapeListener(ed.listener)

	termi.Finish()

	fmt.Print(termi.Clear)
	fmt.Print(termi.HomeCursor)
	termi.Cooked()
	fmt.Print(termi.ShowCursor)

	if ed.path != "" && ed.save {
		ed.Save(ed.path)
	}
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

	return ed.lines[row]
}

func (ed *Editor) CurrentLine() string {
	return ed.Line(ed.row)
}

func (ed *Editor) RuneCount() int {
	return utf8.RuneCountInString(ed.CurrentLine())
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
		lines = append(lines, inputLines...)
		if ed.inpRow+1 <= len(ed.lines)-1 {
			lines = append(lines, ed.lines[ed.inpRow+1:]...)
		}
		ed.lines = lines
		ed.inp.Reset()
		ed.mode = ModeCommand
		ed.MoveLeft(1)
		return
	case ModeSearch:
		ed.mode = ModeCommand
		return
	case ModePrompt:
		ed.mode = ModeCommand
		return
	}
}

func (ed *Editor) Ring(format string, a ...any) {
	ed.message = fmt.Sprintf(format, a...)
}

func (ed *Editor) Unimplemented(name string) {
	ed.Ring("not implemented (" + name + ")")
}
