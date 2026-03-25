package editor

import (
	"strings"

	"tea.kareha.org/lab/levi/internal/console"
	"tea.kareha.org/lab/levi/internal/util"
)

type mode int

const (
	modeCommand = iota
	modeInsert
)

type Editor struct {
	scr *Screen
	kb *Keyboard
	x, y int
	line *strings.Builder
	mode mode
}

func New() *Editor {
	scr := NewScreen()
	kb := NewKeyboard()

	return &Editor{
		scr: &scr,
		kb: &kb,
		x: 0,
		y: 0,
		line: new(strings.Builder),
		mode: modeCommand,
	}
}

func (ed *Editor) addRune(r rune) {
	ed.line.WriteRune(r)
}

func (ed *Editor) drawStatus() {
	_, h := console.Size()

	console.MoveCursor(0, h - 2)
	switch ed.mode {
	case modeCommand:
		console.Print("-- [command] q: quit, i: insert --")
	case modeInsert:
		console.Print("-- [insert] Esc: command mode --")
	}
}

func (ed *Editor) drawBuffer() {
	console.Print(ed.line.String())
}

func (ed *Editor) repaint() {
	console.HideCursor()

	console.Clear()
	console.HomeCursor()

	ed.drawBuffer()
	ed.drawStatus()

	console.MoveCursor(ed.x, ed.y)

	console.ShowCursor()
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
				ed.mode = modeInsert
			}
		case modeInsert:
			switch r {
			case Esc:
				ed.mode = modeCommand
			default:
				ed.addRune(r)
				ed.x += util.RuneWidth(r)
			}
		}
	}
}
