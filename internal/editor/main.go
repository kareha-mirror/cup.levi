package editor

import (
	"unicode/utf8"

	"tea.kareha.org/lab/levi/internal/console"
)

func (ed *Editor) exitInsert() {
	ed.lines[ed.row] = ed.head + ed.insert.String() + ed.tail
	ed.tail = ""
	ed.insert.Reset()
	ed.mode = modeCommand
	ed.moveLeft(1)
}

func (ed *Editor) insertNewline() {
	before := make([]string, 0, len(ed.lines)+1)
	before = append(before, ed.lines[:ed.row]...)
	var after []string
	if ed.row+1 < len(ed.lines) {
		after = ed.lines[ed.row+1:]
	} else {
		after = []string{}
	}
	newLines := []string{
		ed.head + ed.insert.String(),
		ed.tail,
	}
	ed.lines = append(append(before, newLines...), after...)

	ed.row++

	ed.col = 0
	ed.head = ""
	ed.insert.Reset()
}

func (ed *Editor) deleteBefore() {
	if ed.insert.Len() < 1 {
		ed.ring()
		return
	}
	insert := ed.insert.String()
	_, size := utf8.DecodeLastRuneInString(insert)
	insert = insert[:len(insert)-size]
	ed.insert.Reset()
	ed.insert.WriteString(insert)
	ed.col--
}

func (ed *Editor) Main() {
	for {
		ed.repaint()

		k, r := console.ReadKey()
		switch ed.mode {
		case modeCommand:
			switch k {
			case console.KeyNormal:
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
				case 'x':
					ed.deleteRune(1)
				default:
					ed.ring()
				}
			case console.KeyUp:
				ed.moveUp(1)
			case console.KeyDown:
				ed.moveDown(1)
			case console.KeyRight:
				ed.moveRight(1)
			case console.KeyLeft:
				ed.moveLeft(1)
			default:
				ed.ring()
			}
		case modeInsert:
			switch k {
			case console.KeyNormal:
				switch r {
				case console.RuneEscape:
					ed.exitInsert()
				case console.RuneEnter:
					ed.insertNewline()
				case console.RuneBackspace:
					ed.deleteBefore()
				case console.RuneDelete:
					ed.deleteBefore()
				default:
					ed.insertRune(r)
				}
			case console.KeyUp:
				ed.exitInsert()
				ed.moveUp(1)
			case console.KeyDown:
				ed.exitInsert()
				ed.moveDown(1)
			case console.KeyRight:
				ed.exitInsert()
				ed.moveRight(1)
			case console.KeyLeft:
				ed.exitInsert()
				ed.moveLeft(1)
			default:
				ed.ring()
			}
		}
	}
}
