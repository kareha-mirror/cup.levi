package editor

import (
	"fmt"
	"regexp"

	"tea.kareha.org/cup/termi"
)

func (ed *Editor) MainCommand(key termi.Key) {
	switch key.Kind {
	case termi.KeyRune:
		if key.Rune == termi.RuneEscape {
			if ed.parser.String() == "" {
				ed.Notice("Already in command mode")
			}
			ed.Reset()
			return
		}

		if ed.parser.String() == "" {
			switch key.Rune {
			case ':':
				ed.Reset()
				ed.mode = ModePrompt
				return
			case '/':
				ed.Reset()
				ed.mode = ModeSearch
				ed.search.backward = false
				return
			case '?':
				ed.Reset()
				ed.mode = ModeSearch
				ed.search.backward = true
				return
			}
		}

		if key.Rune == termi.RuneBackspace || key.Rune == termi.RuneDelete {
			if !ed.parser.Backspace() {
				ed.Notice("No tokens left")
				return
			}
		} else {
			ed.parser.WriteRune(key.Rune)
		}

		c, ok := ed.Parse()
		if ok {
			if _, ok := InsertCmds[c.Kind]; ok {
				ed.BeginMemory()
			} else if _, ok := EditCmds[c.Kind]; ok {
				ed.BeginMemory()
			}
			if ed.Run(c, false) {
				if _, ok := InsertCmds[c.Kind]; ok {
					ed.lastCmd = c
				} else if _, ok := EditCmds[c.Kind]; ok {
					ed.EndMemory()
					ed.lastCmd = c
				}
				ed.parser.Reset()
			} else {
				if _, ok := InsertCmds[c.Kind]; ok {
					ed.CancelMemory()
				} else if _, ok := EditCmds[c.Kind]; ok {
					ed.CancelMemory()
				}
			}
		}
	case termi.KeyUp:
		ed.Run(Cmd{Kind: CmdMoveUp, Num: 1}, false)
	case termi.KeyDown:
		ed.Run(Cmd{Kind: CmdMoveDown, Num: 1}, false)
	case termi.KeyRight:
		ed.Run(Cmd{Kind: CmdMoveRight, Num: 1}, false)
	case termi.KeyLeft:
		ed.Run(Cmd{Kind: CmdMoveLeft, Num: 1}, false)
	default:
		ed.Ring("unknown key")
	}
}

func (ed *Editor) MainInsert(key termi.Key) {
	switch key.Kind {
	case termi.KeyRune:
		switch key.Rune {
		case termi.RuneEscape:
			ed.Commit()
		case termi.RuneEnter, termi.RuneNewline:
			ed.InputNewline()
		case termi.RuneBackspace, termi.RuneDelete:
			ed.InputBackspace()
		default:
			ed.InputWriteRune(key.Rune)
		}
	case termi.KeyUp:
		ed.Run(Cmd{Kind: CmdMoveUp, Num: 1}, false)
	case termi.KeyDown:
		ed.Run(Cmd{Kind: CmdMoveDown, Num: 1}, false)
	case termi.KeyRight:
		ed.Run(Cmd{Kind: CmdMoveRight, Num: 1}, false)
	case termi.KeyLeft:
		ed.Run(Cmd{Kind: CmdMoveLeft, Num: 1}, false)
	default:
		ed.Error("Unknown key")
	}
}

func (ed *Editor) MainPrompt(key termi.Key) {
	switch key.Kind {
	case termi.KeyRune:
		switch key.Rune {
		case termi.RuneEscape:
			ed.prompt.Reset()
			ed.mode = ModeCommand
		case termi.RuneEnter, termi.RuneNewline:
			c, ok := ed.ParsePrompt()
			ed.prompt.Reset()
			if ok {
				ok = ed.RunPrompt(c)
				if !ok {
					ed.Error("Prompt command failed")
				}
			} else {
				ed.mode = ModeCommand
			}
		case termi.RuneBackspace, termi.RuneDelete:
			if !ed.prompt.RemoveTail() {
				ed.mode = ModeCommand
			}
		default:
			ed.prompt.WriteRune(key.Rune)
		}
	default:
		ed.Error("Unknown key")
	}
}

func (ed *Editor) MainSearch(key termi.Key) {
	switch key.Kind {
	case termi.KeyRune:
		switch key.Rune {
		case termi.RuneEscape:
			ed.search.pattern.Reset()
			ed.mode = ModeCommand
		case termi.RuneEnter, termi.RuneNewline:
			if ed.search.pattern.Len() < 1 {
				if ed.search.backward {
					ed.Run(Cmd{
						Kind: CmdMoveSearchRepeatBackward,
					}, false)
				} else {
					ed.Run(Cmd{
						Kind: CmdMoveSearchRepeatForward,
					}, false)
				}
				return
			}
			re, err := regexp.Compile(ed.search.pattern.String())
			if err != nil {
				ed.Ring("%v", err)
				return
			}
			ed.search.regexp = re
			ed.search.pattern.Reset()
			if ed.search.backward {
				ed.Run(Cmd{Kind: CmdMoveSearchBackward}, false)
			} else {
				ed.Run(Cmd{Kind: CmdMoveSearchForward}, false)
			}
		case termi.RuneBackspace, termi.RuneDelete:
			if !ed.search.pattern.RemoveTail() {
				ed.mode = ModeCommand
			}
		default:
			ed.search.pattern.WriteRune(key.Rune)
		}
	default:
		ed.Error("Unknown key")
	}
}

func (ed *Editor) Main() {
	for ed.alive {
		ed.Draw()

		select {
		case key := <-termi.Keys():
			switch ed.mode {
			case ModeCommand:
				ed.MainCommand(key)
			case ModeInsert:
				ed.MainInsert(key)
			case ModePrompt:
				ed.MainPrompt(key)
			case ModeSearch:
				ed.MainSearch(key)
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
