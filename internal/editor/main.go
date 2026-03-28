package editor

import (
	"tea.kareha.org/lab/termi"
)

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
	lines := ed.inp.Newline()
	ed.lines = append(append(before, lines...), after...)
	ed.row++
	ed.col = 0
	// row and col are already confined
}

func (ed *Editor) Backspace() {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	if !ed.inp.Backspace() {
		ed.Ring("nothing to delete, input is empty")
		return
	}
	ed.col--
	// col is already confined
}

func (ed *Editor) Main() {
	for !ed.quit {
		ed.Draw()

		key := termi.ReadKey()
		switch ed.mode {
		case ModeCommand:
			switch key.Kind {
			case termi.KeyRune:
				if key.Rune == termi.RuneEscape {
					if ed.parser.String() == "" {
						ed.Ring("already in vi command mode")
					}
					ed.parser.ClearAll()
					continue
				}
				ed.parser.InsertRune(key.Rune)

				c, ok := ed.parser.Parse()
				if ok {
					if ed.Run(c) {
						ed.parser.Clear()
					}
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
				ed.Ring("unknown key")
			}
		case ModeInsert:
			switch key.Kind {
			case termi.KeyRune:
				switch key.Rune {
				case termi.RuneEscape:
					ed.EnsureCommand()
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
				ed.MoveUp(1)
			case termi.KeyDown:
				ed.MoveDown(1)
			case termi.KeyRight:
				ed.MoveRight(1)
			case termi.KeyLeft:
				ed.MoveLeft(1)
			default:
				ed.Ring("unknown key")
			}
		}
	}
}
