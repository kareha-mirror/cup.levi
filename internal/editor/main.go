package editor

import (
	"fmt"

	"tea.kareha.org/cup/termi"
)

func (ed *Editor) Main() {
	for ed.alive {
		ed.Draw()

		select {
		case key := <-termi.Keys():
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

					if ed.parser.String() == "" && key.Rune == ':' {
						ed.parser.ClearAll()
						ed.mode = ModePrompt
						continue
					}
					ed.parser.InsertRune(key.Rune)

					c, ok := ed.Parse()
					if ok {
						if ed.Run(c, false) {
							if RepeatableCmds[c.Kind] {
								ed.lastCmd = c
							}
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
					case termi.RuneEnter, termi.RuneNewline:
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
			case ModePrompt:
				switch key.Kind {
				case termi.KeyRune:
					switch key.Rune {
					case termi.RuneEscape:
						ed.prompt.Reset()
						ed.mode = ModeCommand
					case termi.RuneEnter, termi.RuneNewline:
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
					case termi.RuneBackspace, termi.RuneDelete:
						if !ed.prompt.RemoveTail() {
							ed.mode = ModeCommand
						}
					default:
						ed.prompt.WriteRune(key.Rune)
					}
				default:
					ed.Ring("unknown key")
				}
			}
		case sig := <-termi.Sigs():
			if sig == termi.SigStop {
				fmt.Print(termi.Clear)
				fmt.Print(termi.HomeCursor)
				termi.StopKey()
				fmt.Print(termi.ResetAlternate)
				termi.Cooked()
				fmt.Print(termi.ShowCursor)
				ed.redraw = true

				termi.ForceSuspend()
				for {
					sig := <-termi.Sigs()
					if sig == termi.SigCont {
						termi.Raw()
						fmt.Print(termi.SetAlternate)
						termi.StartKey()
						break
					}
				}
			}
		}
	}
}
