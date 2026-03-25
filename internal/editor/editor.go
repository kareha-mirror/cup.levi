package editor

import (
	"strings"
	"unicode/utf8"

	"tea.kareha.org/lab/levi/internal/console"
	"tea.kareha.org/lab/levi/internal/util"
)

type mode int

const (
	modeCommand = iota
	modeInsert
)

type Editor struct {
	scr        *Screen
	kb         *Keyboard
	col, row   int
	x, y       int
	head, tail string
	insert     *strings.Builder
	mode       mode
}

func New() *Editor {
	scr := NewScreen()
	kb := NewKeyboard()

	return &Editor{
		scr:    &scr,
		kb:     &kb,
		col:    0,
		row:    0,
		x:      0,
		y:      0,
		head:   "",
		tail:   "",
		insert: new(strings.Builder),
		mode:   modeCommand,
	}
}

func (ed *Editor) addRune(r rune) {
	ed.insert.WriteRune(r)
}

func (ed *Editor) drawBuffer() {
	switch ed.mode {
	case modeCommand:
		console.Print(ed.head)
	case modeInsert:
		console.Print(ed.head)
		console.Print(ed.insert.String())
		console.Print(ed.tail)
	}
}

func (ed *Editor) drawStatus() {
	_, h := console.Size()

	console.MoveCursor(0, h-2)
	switch ed.mode {
	case modeCommand:
		console.Print("-- [command] q: quit, i: insert, a: insert after --")
	case modeInsert:
		console.Print("-- [insert] Esc: command mode --")
	}
}

func (ed *Editor) updateCursor() {
	switch ed.mode {
	case modeCommand:
		ed.x = util.StringWidth(ed.head, ed.col)
	case modeInsert:
		ed.x = util.StringWidth(ed.head+ed.insert.String(), ed.col)
	}
	ed.y = ed.row
}

func (ed *Editor) repaint() {
	console.HideCursor()

	console.Clear()
	console.HomeCursor()

	ed.drawBuffer()
	ed.drawStatus()

	ed.updateCursor()
	console.MoveCursor(ed.x, ed.y)

	console.ShowCursor()
}

func (ed *Editor) enterInsert() {
	rs := []rune(ed.head)
	ed.head = string(rs[:ed.col])
	ed.tail = string(rs[ed.col:])
	ed.mode = modeInsert
}

func (ed *Editor) enterInsertAfter() {
	len := utf8.RuneCountInString(ed.head)
	if ed.col < len-1 {
		ed.moveRight(1)
		ed.enterInsert()
		return
	} else {
		ed.col++
		ed.mode = modeInsert
	}
}

func (ed *Editor) moveLeft(n int) {
	ed.col = max(ed.col-n, 0)
}

func (ed *Editor) moveRight(n int) {
	ed.col = min(ed.col+n, max(utf8.RuneCountInString(ed.head)-1, 0))
}

func (ed *Editor) Main() {
	for {
		ed.repaint()

		r := ed.kb.ReadRune()
		switch ed.mode {
		case modeCommand:
			switch r {
			case 'q':
				return
			case 'i':
				ed.enterInsert()
			case 'a':
				ed.enterInsertAfter()
			case 'h':
				ed.moveLeft(1)
			case 'l':
				ed.moveRight(1)
			}
		case modeInsert:
			switch r {
			case Esc:
				ed.head = ed.head + ed.insert.String() + ed.tail
				ed.tail = ""
				ed.insert = new(strings.Builder)
				ed.mode = modeCommand
				ed.moveLeft(1)
			default:
				ed.addRune(r)
				ed.col++
			}
		}
	}
}
