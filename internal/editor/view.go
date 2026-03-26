package editor

import (
	"unicode/utf8"

	"tea.kareha.org/lab/termi"
)

func (ed *Editor) lineHeight(line string) int {
	w, _ := termi.Size()
	rc := utf8.RuneCountInString(line)
	width := termi.StringWidth(line, rc)
	return 1 + max(width-1, 0)/w
}

func (ed *Editor) drawBuffer() {
	_, h := termi.Size()

	y := 0
	for i := ed.vrow; i < len(ed.lines); i++ {
		var line string
		if ed.mode == modeInsert && i == ed.row {
			line = ed.head + ed.insert.String() + ed.tail
		} else {
			line = ed.lines[i]
		}

		termi.MoveCursor(0, y)
		termi.Draw(line)

		y += ed.lineHeight(line)
		if y >= h-1 {
			break
		}
	}

	for ; y < h-1; y++ {
		termi.MoveCursor(0, y)
		termi.Draw("~")
	}
}

func (ed *Editor) drawStatus() {
	var m string
	switch ed.mode {
	case modeCommand:
		m = "c"
	case modeInsert:
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

func (ed *Editor) updateCursor() {
	w, h := termi.Size()

	var dy int
	switch ed.mode {
	case modeCommand:
		ed.row = min(max(ed.row, 0), max(len(ed.lines)-1, 0))
		len := ed.runeCount()
		ed.col = min(ed.col, max(len-1, 0))

		// XXX approximation
		width := termi.StringWidth(ed.lines[ed.row], ed.col)
		ed.x = width % w
		dy = width / w
	case modeInsert:
		// XXX approximation
		width := termi.StringWidth(ed.head+ed.insert.String(), ed.col)
		ed.x = width % w
		dy = width / w
	}

	if ed.row < ed.vrow {
		ed.vrow = ed.row
	}

	y := 0
	for i := ed.vrow; i < ed.row; i++ {
		y += ed.lineHeight(ed.lines[i])
	}
	ed.y = y + dy

	for ed.y >= h-1 {
		ed.vrow++

		y := 0
		for i := ed.vrow; i < ed.row; i++ {
			y += ed.lineHeight(ed.lines[i])
		}
		ed.y = y + dy
	}
}

func (ed *Editor) repaint() {
	termi.HideCursor()

	termi.Clear()
	termi.HomeCursor()

	ed.updateCursor()

	ed.drawBuffer()
	ed.drawStatus()

	termi.MoveCursor(ed.x, ed.y)

	termi.ShowCursor()
}
