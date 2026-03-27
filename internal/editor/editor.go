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
)

type Editor struct {
	col, row int
	x, y     int
	vrow     int
	lines    []string
	ins      *Insert
	mode     Mode
	path     string
	bell     bool
}

func (ed *Editor) Load() {
	if ed.path == "" {
		return
	}
	_, err := os.Stat(ed.path)
	if err != nil { // file not exists
		return
	}
	data, err := os.ReadFile(ed.path)
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
}

func Init(args []string) *Editor {
	var path string
	if len(args) > 1 {
		path = args[1]
	}

	ed := &Editor{
		col:   0,
		row:   0,
		x:     0,
		y:     0,
		vrow:  0,
		lines: make([]string, 1),
		ins:   NewInsert(),
		mode:  ModeCommand,
		path:  path,
		bell:  false,
	}

	ed.Load()

	termi.Raw()
	return ed
}

func (ed *Editor) Save() {
	if ed.path == "" {
		return
	}
	text := strings.Join(ed.lines, "\n") + "\n"
	err := os.WriteFile(ed.path, []byte(text), 0644)
	if err != nil {
		panic(err)
	}
}

func (ed *Editor) Finish() {
	termi.Clear()
	termi.HomeCursor()
	termi.Cooked()
	termi.ShowCursor()

	ed.Save()
}

func (ed *Editor) RuneCount() int {
	return utf8.RuneCountInString(ed.lines[ed.row])
}

func (ed *Editor) InsertRune(r rune) {
	ed.ins.Write(r)
	ed.col++
}

func (ed *Editor) Ring() {
	ed.bell = true
}
