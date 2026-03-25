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
	scr        *screen
	kb         *keyboard
	col, row   int
	x, y       int
	lines      []string
	head, tail string
	insert     *strings.Builder
	mode       mode
}

func Init() *Editor {
	console.Raw()

	scr := newScreen()
	kb := newKeyboard()

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

func (ed *Editor) lineHeight(line string) int {
	w, _ := ed.scr.size()
	rc := utf8.RuneCountInString(line)
	width := util.StringWidth(line, rc)
	return 1 + max(width-1, 0)/w
}

func (ed *Editor) drawBuffer() {
	_, h := ed.scr.size()

	y := 0
	for i := 0; i < len(ed.lines); i++ {
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
	_, h := ed.scr.size()

	console.MoveCursor(0, h-1)
	switch ed.mode {
	case modeCommand:
		console.Print("c")
	case modeInsert:
		console.Print("i")
	}
}

func (ed *Editor) updateCursor() {
	w, _ := ed.scr.size()

	var dy int
	switch ed.mode {
	case modeCommand:
		ed.row = min(max(ed.row, 0), max(len(ed.lines)-1, 0))
		len := ed.runeCount()
		ed.col = min(ed.col, max(len-1, 0))

		// XXX approximation
		width := util.StringWidth(ed.lines[ed.row], ed.col)
		ed.x = width % w
		dy = width / w
	case modeInsert:
		// XXX approximation
		width := util.StringWidth(ed.head+ed.insert.String(), ed.col)
		ed.x = width % w
		dy = width / w
	}

	y := 0
	for i := 0; i < ed.row; i++ {
		var line string
		if ed.mode == modeInsert && i == ed.row {
			line = ed.head + ed.insert.String() + ed.tail
		} else {
			line = ed.lines[i]
		}

		y += ed.lineHeight(line)
	}
	ed.y = y + dy
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
	lines := make([]string, 0, len(ed.lines)+1)
	lines = append(lines, ed.lines[:ed.row+1]...)
	lines = append(lines, "")
	if ed.row+1 < len(ed.lines) {
		lines = append(lines, ed.lines[ed.row+1:]...)
	}

	lines[ed.row] = ed.head + ed.insert.String()
	lines[ed.row+1] = ed.tail
	ed.lines = lines
	ed.row++

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

		k, r := ed.kb.readKey()
		switch ed.mode {
		case modeCommand:
			switch k {
			case keyNormal:
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
			case keyUp:
				ed.moveUp(1)
			case keyDown:
				ed.moveDown(1)
			case keyRight:
				ed.moveRight(1)
			case keyLeft:
				ed.moveLeft(1)
			default:
				// TODO ring
			}
		case modeInsert:
			switch k {
			case keyNormal:
				switch r {
				case runeEscape:
					ed.exitInsert()
				case runeEnter:
					ed.insertNewline()
				default:
					ed.insertRune(r)
				}
			case keyUp:
				ed.exitInsert()
				ed.moveUp(1)
			case keyDown:
				ed.exitInsert()
				ed.moveDown(1)
			case keyRight:
				ed.exitInsert()
				ed.moveRight(1)
			case keyLeft:
				ed.exitInsert()
				ed.moveLeft(1)
			default:
				// TODO ring
			}
		}
	}
}
