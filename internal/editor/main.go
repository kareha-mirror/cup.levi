package editor

import (
	"tea.kareha.org/cup/termi"
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

		seq := termi.ReadSeq()
		switch ed.mode {
		case ModeCommand:
			switch seq.Kind {
			case termi.SeqRune:
				if seq.Rune == termi.RuneEscape {
					if ed.parser.String() == "" {
						ed.Ring("already in vi command mode")
					}
					ed.parser.ClearAll()
					continue
				}

				if ed.parser.String() == "" && seq.Rune == ':' {
					ed.mode = ModePrompt
					continue
				}
				ed.parser.InsertRune(seq.Rune)

				c, ok := ed.parser.Parse()
				if ok {
					if ed.Run(c) {
						ed.parser.Clear()
					}
				}
			case termi.SeqUp:
				ed.MoveUp(1)
			case termi.SeqDown:
				ed.MoveDown(1)
			case termi.SeqRight:
				ed.MoveRight(1)
			case termi.SeqLeft:
				ed.MoveLeft(1)
			default:
				ed.Ring("unknown sequence")
			}
		case ModeInsert:
			switch seq.Kind {
			case termi.SeqRune:
				switch seq.Rune {
				case termi.RuneEscape:
					ed.EnsureCommand()
				case termi.RuneEnter:
					ed.InsertNewline()
				case termi.RuneBackspace:
					ed.Backspace()
				case termi.RuneDelete:
					ed.Backspace()
				default:
					ed.InsertRune(seq.Rune)
				}
			case termi.SeqUp:
				ed.MoveUp(1)
			case termi.SeqDown:
				ed.MoveDown(1)
			case termi.SeqRight:
				ed.MoveRight(1)
			case termi.SeqLeft:
				ed.MoveLeft(1)
			default:
				ed.Ring("unknown sequence")
			}
		case ModePrompt:
			switch seq.Kind {
			case termi.SeqRune:
				switch seq.Rune {
				case termi.RuneEscape:
					ed.prompt.Reset()
					ed.mode = ModeCommand
				case termi.RuneEnter:
					c, ok := ed.ParsePrompt()
					if ok {
						ed.prompt.Reset()
						ok = ed.RunPrompt(c)
						if !ok {
							ed.Ring("prompt command failed")
						}
					} else {
						ed.Ring("unknown prompt command")
					}
				case termi.RuneBackspace:
					ed.prompt.RemoveTail()
				case termi.RuneDelete:
					ed.prompt.RemoveTail()
				default:
					ed.prompt.WriteRune(seq.Rune)
				}
			default:
				ed.Ring("unknown sequence")
			}
		}
	}
}
