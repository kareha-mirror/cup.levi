package editor

import (
	"os"
	"strings"
	"unicode/utf8"

	"tea.kareha.org/lab/termi"
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
	w, h     int
	x, y     int
	lines    []string
	inp      *Input
	mode     Mode
	path     string
	message  string
	parser   *Parser
	quit     bool
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
		col:     0,
		row:     0,
		vrow:    0,
		w:       w,
		h:       h,
		x:       0,
		y:       0,
		lines:   make([]string, 1),
		inp:     NewInput(),
		mode:    ModeCommand,
		path:    path,
		message: "",
		parser:  NewParser(),
		quit:    false,
	}

	if path != "" {
		_, err := os.Stat(path)
		if err == nil { // file exists
			ed.Load(path)
		}
	}

	termi.Raw()
	return ed
}

func (ed *Editor) Save(path string) {
	text := strings.Join(ed.lines, "\n") + "\n"
	err := os.WriteFile(path, []byte(text), 0644)
	if err != nil {
		panic(err)
	}
	ed.path = path
}

func (ed *Editor) Finish() {
	termi.Clear()
	termi.HomeCursor()
	termi.Cooked()
	termi.ShowCursor()

	if ed.path != "" {
		ed.Save(ed.path)
	}
}

func (ed *Editor) Line(row int) string {
	if ed.mode == ModeInsert && row == ed.row {
		return ed.inp.Line()
	} else {
		return ed.lines[row]
	}
}

func (ed *Editor) CurrentLine() string {
	return ed.Line(ed.row)
}

func (ed *Editor) RuneCount() int {
	return utf8.RuneCountInString(ed.CurrentLine())
}

func (ed *Editor) Confine() {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}

	n := len(ed.lines)
	if ed.row < 0 {
		ed.row = 0
	} else if ed.row >= n {
		ed.row = max(n-1, 0)
	}

	rc := ed.RuneCount()
	if ed.col < 0 {
		ed.col = 0
	} else if ed.col >= rc {
		ed.col = max(rc-1, 0)
	}
}

func (ed *Editor) InsertRune(r rune) {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	ed.inp.WriteRune(r)
	ed.col = ed.inp.Column()
}

func (ed *Editor) Ring(message string) {
	ed.message = message
}
