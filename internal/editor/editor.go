package editor

import (
	"os"
	"strings"
	"unicode/utf8"

	"tea.kareha.org/lab/levi/internal/console"
)

type mode int

const (
	modeCommand mode = iota
	modeInsert
)

type Editor struct {
	col, row   int
	x, y       int
	vrow       int
	lines      []string
	head, tail string
	insert     *strings.Builder
	mode       mode
	path       string
	bell       bool
}

func (ed *Editor) load() {
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
		col:    0,
		row:    0,
		x:      0,
		y:      0,
		vrow:   0,
		lines:  make([]string, 1),
		head:   "",
		tail:   "",
		insert: new(strings.Builder),
		mode:   modeCommand,
		path:   path,
		bell:   false,
	}

	ed.load()

	console.Raw()
	return ed
}

func (ed *Editor) save() {
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
	console.Clear()
	console.HomeCursor()
	console.Cooked()
	console.ShowCursor()

	ed.save()
}

func (ed *Editor) runeCount() int {
	return utf8.RuneCountInString(ed.lines[ed.row])
}

func (ed *Editor) insertRune(r rune) {
	ed.insert.WriteRune(r)
	ed.col++
}

func (ed *Editor) ring() {
	ed.bell = true
}
