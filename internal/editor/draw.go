package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/termi"
)

func (ed *Editor) LineHeight(line string) int {
	rc := utf8.RuneCountInString(line)
	width := termi.StringWidth(line, rc)
	return 1 + max(width-1, 0)/ed.w
}

func (ed *Editor) DrawBuffer() {
	y := 0
	for i := ed.vrow; i < len(ed.lines); i++ {
		line := ed.Line(i)

		termi.MoveCursor(0, y)
		termi.Draw(line)

		y += ed.LineHeight(line)
		if y >= ed.h-1 {
			break
		}
	}

	for ; y < ed.h-1; y++ {
		termi.MoveCursor(0, y)
		termi.Draw("~")
	}
}

func (ed *Editor) DrawStatus() {
	var m string
	switch ed.mode {
	case ModeCommand:
		m = "vi command"
	case ModeInsert:
		m = "vi insert"
	case ModeSearch:
		m = "vi search"
	case ModePrompt:
		m = "vi prompt"
	default:
		panic("invalid mode")
	}

	termi.MoveCursor(0, ed.h-1)
	if ed.message != "" {
		termi.EnableInvert()
		termi.Print(ed.message)
		termi.DisableInvert()
		ed.message = ""
	} else {
		termi.Printf("[%s] %s %d,%d %s", ed.parser.Cache(), m, ed.row, ed.col, ed.path)
	}
}

func (ed *Editor) UpdateCursor() {
	// XXX approximation
	width := termi.StringWidth(ed.CurrentLine(), ed.col)
	ed.x = width % ed.w
	dy := width / ed.w

	if ed.row < ed.vrow {
		ed.vrow = ed.row
	}

	y := 0
	for i := ed.vrow; i < ed.row; i++ {
		y += ed.LineHeight(ed.lines[i])
	}
	ed.y = y + dy

	for ed.y >= ed.h-1 {
		ed.vrow++

		y := 0
		for i := ed.vrow; i < ed.row; i++ {
			y += ed.LineHeight(ed.lines[i])
		}
		ed.y = y + dy
	}
}

func (ed *Editor) Repaint() {
	w, h := termi.Size()
	ed.w = w
	ed.h = h

	termi.HideCursor()

	termi.Clear()
	termi.HomeCursor()

	ed.UpdateCursor()

	ed.DrawBuffer()
	ed.DrawStatus()

	termi.MoveCursor(ed.x, ed.y)

	termi.ShowCursor()
}

func (ed *Editor) Draw() {
	ed.Repaint()
}
