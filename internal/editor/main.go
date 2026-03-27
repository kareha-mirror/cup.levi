package editor

import (
	"tea.kareha.org/lab/termi"
)

func (ed *Editor) ExitInsert() {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	ed.lines[ed.row] = ed.ins.Line()
	ed.ins.Reset()
	ed.mode = ModeCommand
	ed.MoveLeft(1)
}

func (ed *Editor) InsertNewline() {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	before := make([]string, 0, len(ed.lines)+1)
	before = append(before, ed.lines[:ed.row]...)
	var after []string
	if ed.row+1 < len(ed.lines) {
		after = ed.lines[ed.row+1:]
	} else {
		after = []string{}
	}
	lines := ed.ins.Newline()
	ed.lines = append(append(before, lines...), after...)
	ed.row++
	ed.col = 0
	// row and col are confined automatically
}

func (ed *Editor) Backspace() {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	if !ed.ins.Backspace() {
		ed.Ring()
		return
	}
	ed.col--
	// col is confined automatically
}

func (ed *Editor) Main() {
	for {
		ed.Draw()

		key := termi.ReadKey()
		switch ed.mode {
		case ModeCommand:
			switch key.Kind {
			case termi.KeyRune:
				switch key.Rune {
				case 'q':
					return
				case 'i':
					ed.Insert()
				case 'a':
					ed.InsertAfter()
				case 'h':
					ed.MoveLeft(1)
				case 'l':
					ed.MoveRight(1)
				case 'j':
					ed.MoveDown(1)
				case 'k':
					ed.MoveUp(1)
				case 'x':
					ed.DeleteRune(1)
				default:
					ed.Ring()
				}
			case termi.KeyUp:
				ed.MoveUp(1)
			case termi.KeyDown:
				ed.MoveDown(1)
			case termi.KeyRight:
				ed.MoveRight(1)
			case termi.KeyLeft:
				ed.MoveLeft(1)
			default:
				ed.Ring()
			}
		case ModeInsert:
			switch key.Kind {
			case termi.KeyRune:
				switch key.Rune {
				case termi.RuneEscape:
					ed.ExitInsert()
				case termi.RuneEnter:
					ed.InsertNewline()
				case termi.RuneBackspace:
					ed.Backspace()
				case termi.RuneDelete:
					ed.Backspace()
				default:
					ed.InsertRune(key.Rune)
				}
			case termi.KeyUp:
				ed.ExitInsert()
				ed.MoveUp(1)
			case termi.KeyDown:
				ed.ExitInsert()
				ed.MoveDown(1)
			case termi.KeyRight:
				ed.ExitInsert()
				ed.MoveRight(1)
			case termi.KeyLeft:
				ed.ExitInsert()
				ed.MoveLeft(1)
			default:
				ed.Ring()
			}
		}
	}
}
