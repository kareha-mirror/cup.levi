package editor

import (
	"unicode/utf8"

	"tea.kareha.org/lab/levi/internal/console"
)

func (ed *Editor) lineHeight(line string) int {
	w, _ := console.Size()
	rc := utf8.RuneCountInString(line)
	width := console.StringWidth(line, rc)
	return 1 + max(width-1, 0)/w
}

func (ed *Editor) drawBuffer() {
	_, h := console.Size()

	y := 0
	for i := ed.vrow; i < len(ed.lines); i++ {
		var line string
		if ed.mode == modeInsert && i == ed.row {
			line = ed.head + ed.insert.String() + ed.tail
		} else {
			line = ed.lines[i]
		}

		console.MoveCursor(0, y)
		console.Print(line)

		y += ed.lineHeight(line)
		if y >= h-1 {
			break
		}
	}

	for ; y < h-1; y++ {
		console.MoveCursor(0, y)
		console.Print("~")
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

	_, h := console.Size()
	console.MoveCursor(0, h-1)
	if ed.bell {
		console.EnableInvert()
	}
	console.Printf("%s %d,%d %s", m, ed.row, ed.col, ed.path)
	if ed.bell {
		console.DisableInvert()
	}
	ed.bell = false
}

func (ed *Editor) updateCursor() {
	w, h := console.Size()

	var dy int
	switch ed.mode {
	case modeCommand:
		ed.row = min(max(ed.row, 0), max(len(ed.lines)-1, 0))
		len := ed.runeCount()
		ed.col = min(ed.col, max(len-1, 0))

		// XXX approximation
		width := console.StringWidth(ed.lines[ed.row], ed.col)
		ed.x = width % w
		dy = width / w
	case modeInsert:
		// XXX approximation
		width := console.StringWidth(ed.head+ed.insert.String(), ed.col)
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
	console.HideCursor()

	console.Clear()
	console.HomeCursor()

	ed.updateCursor()

	ed.drawBuffer()
	ed.drawStatus()

	console.MoveCursor(ed.x, ed.y)

	console.ShowCursor()
}
