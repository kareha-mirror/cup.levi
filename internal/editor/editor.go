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
	lines      []string
	head, tail string
	insert     *strings.Builder
	mode       mode
}

func Init() *Editor {
	console.Raw()

	scr := NewScreen()
	kb := NewKeyboard()

	return &Editor{
		scr:    &scr,
		kb:     &kb,
		col:    0,
		row:    0,
		x:      0,
		y:      0,
		lines:  make([]string, 1),
		head:   "",
		tail:   "",
		insert: new(strings.Builder),
		mode:   modeCommand,
	}
}

func (ed *Editor) Finish() {
	console.Clear()
	console.HomeCursor()
	console.Cooked()
	console.ShowCursor()
}

func (ed *Editor) runeCount() int {
	return utf8.RuneCountInString(ed.lines[ed.row])
}

func (ed *Editor) drawBuffer() {
	for i := 0; i < len(ed.lines); i++ {
		console.MoveCursor(0, i)
		if ed.mode == modeInsert && i == ed.row {
			console.Print(ed.head)
			console.Print(ed.insert.String())
			console.Print(ed.tail)
		} else {
			console.Print(ed.lines[i])
		}
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
		ed.row = min(max(ed.row, 0), max(len(ed.lines)-1, 0))
		len := ed.runeCount()
		ed.col = min(ed.col, max(len-1, 0))
		ed.x = util.StringWidth(ed.lines[ed.row], ed.col)
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

func (ed *Editor) exitInsert() {
	ed.lines[ed.row] = ed.head + ed.insert.String() + ed.tail
	ed.tail = ""
	ed.insert = new(strings.Builder)
	ed.mode = modeCommand
	ed.moveLeft(1)
}

func (ed *Editor) insertNewline() {
	ed.lines[ed.row] = ed.head + ed.insert.String()
	ed.lines = append(ed.lines, "")
	copy(ed.lines[ed.row+1:], ed.lines[ed.row:])
	ed.row++
	ed.lines[ed.row] = ed.tail
	ed.col = 0
	ed.head = ""
	ed.insert = new(strings.Builder)
}

func (ed *Editor) insertRune(r rune) {
	ed.insert.WriteRune(r)
	ed.col++
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
			case 'j':
				ed.moveDown(1)
			case 'k':
				ed.moveUp(1)
			}
		case modeInsert:
			switch r {
			case Escape:
				ed.exitInsert()
			case Enter:
				ed.insertNewline()
			default:
				ed.insertRune(r)
			}
		}
	}
}
