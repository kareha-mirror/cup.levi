package editor

import (
	"unicode/utf8"

	"tea.kareha.org/lab/termi"
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

		key := termi.ReadKey()
		switch ed.mode {
		case modeCommand:
			switch key.Kind {
			case termi.KeyRune:
				switch key.Rune {
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
			case termi.KeyUp:
				ed.moveUp(1)
			case termi.KeyDown:
				ed.moveDown(1)
			case termi.KeyRight:
				ed.moveRight(1)
			case termi.KeyLeft:
				ed.moveLeft(1)
			default:
				ed.ring()
			}
		case modeInsert:
			switch key.Kind {
			case termi.KeyRune:
				switch key.Rune {
				case termi.RuneEscape:
					ed.exitInsert()
				case termi.RuneEnter:
					ed.insertNewline()
				case termi.RuneBackspace:
					ed.deleteBefore()
				case termi.RuneDelete:
					ed.deleteBefore()
				default:
					ed.insertRune(key.Rune)
				}
			case termi.KeyUp:
				ed.exitInsert()
				ed.moveUp(1)
			case termi.KeyDown:
				ed.exitInsert()
				ed.moveDown(1)
			case termi.KeyRight:
				ed.exitInsert()
				ed.moveRight(1)
			case termi.KeyLeft:
				ed.exitInsert()
				ed.moveLeft(1)
			default:
				ed.ring()
			}
		}
	}
}
