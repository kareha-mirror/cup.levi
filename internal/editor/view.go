package editor

import (
	"unicode/utf8"

	"tea.kareha.org/lab/termi"
)

func (ed *Editor) LineHeight(line string) int {
	w, _ := termi.Size()
	rc := utf8.RuneCountInString(line)
	width := termi.StringWidth(line, rc)
	return 1 + max(width-1, 0)/w
}

func (ed *Editor) DrawBuffer() {
	_, h := termi.Size()

	y := 0
	for i := ed.vrow; i < len(ed.lines); i++ {
		var line string
		if ed.mode == ModeInsert && i == ed.row {
			line = ed.ins.Line()
		} else {
			line = ed.lines[i]
		}

		termi.MoveCursor(0, y)
		termi.Draw(line)

		y += ed.LineHeight(line)
		if y >= h-1 {
			break
		}
	}

	for ; y < h-1; y++ {
		termi.MoveCursor(0, y)
		termi.Draw("~")
	}
}

func (ed *Editor) DrawStatus() {
	var m string
	switch ed.mode {
	case ModeCommand:
		m = "c"
	case ModeInsert:
		m = "i"
	}

	_, h := termi.Size()
	termi.MoveCursor(0, h-1)
	if ed.bell {
		termi.EnableInvert()
	}
	termi.Printf("%s %d,%d %s", m, ed.row, ed.col, ed.path)
	if ed.bell {
		termi.DisableInvert()
	}
	ed.bell = false
}

func (ed *Editor) UpdateCursor() {
	w, h := termi.Size()

	var dy int
	switch ed.mode {
	case ModeCommand:
		ed.row = min(max(ed.row, 0), max(len(ed.lines)-1, 0))
		len := ed.RuneCount()
		ed.col = min(ed.col, max(len-1, 0))

		// XXX approximation
		width := termi.StringWidth(ed.lines[ed.row], ed.col)
		ed.x = width % w
		dy = width / w
	case ModeInsert:
		// XXX approximation
		width := ed.ins.Width()
		ed.x = width % w
		dy = width / w
	}

	if ed.row < ed.vrow {
		ed.vrow = ed.row
	}

	y := 0
	for i := ed.vrow; i < ed.row; i++ {
		y += ed.LineHeight(ed.lines[i])
	}
	ed.y = y + dy

	for ed.y >= h-1 {
		ed.vrow++

		y := 0
		for i := ed.vrow; i < ed.row; i++ {
			y += ed.LineHeight(ed.lines[i])
		}
		ed.y = y + dy
	}
}

func (ed *Editor) Repaint() {
	termi.HideCursor()

	termi.Clear()
	termi.HomeCursor()

	ed.UpdateCursor()

	ed.DrawBuffer()
	ed.DrawStatus()

	termi.MoveCursor(ed.x, ed.y)

	termi.ShowCursor()
}
