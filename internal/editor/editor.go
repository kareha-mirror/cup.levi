package editor

import (
	"strings"

	"tea.kareha.org/lab/levi/internal/console"
)

type Editor struct {
	scr *Screen
	kb *Keyboard
	x, y int
	line *strings.Builder
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
	}
}

func (ed *Editor) addRune(r rune) {
	ed.line.WriteRune(r)
}

func (ed *Editor) draw() {
	console.Clear()
	console.HomeCursor()

	console.Print("Hit Esc to Exit")

	console.MoveCursor(ed.x, ed.y)
	console.Print(ed.line.String())
}

func (ed *Editor) Main() {
	for {
		console.HideCursor()
		ed.draw()
		console.ShowCursor()

		r := ed.kb.ReadRune()
		if r == Esc {
			break
		}
		ed.addRune(r)
	}
}
