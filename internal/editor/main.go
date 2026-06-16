package editor

import (
	"tea.kareha.org/cup/termi"
)

func (ed *Editor) Main() {
	for ed.alive {
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
					if ed.Run(c, false) {
						if RepeatableCmds[c.Kind] {
							ed.lastCmd = c
						}
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
				case termi.RuneEnter, '\n':
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
				case termi.RuneEnter, '\n':
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
