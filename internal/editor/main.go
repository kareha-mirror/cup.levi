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
				ed.searchs.backward = false
				return
			case '?':
				ed.Reset()
				ed.mode = ModeSearch
				ed.searchs.backward = true
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
			b := ed.Buf()
			prevRow := b.Loc.Row
			if _, ok := IsModifyingCmd[c.Op.Kind]; ok {
				ed.BeginUndoRecord()
			}
			if modified, ok := ed.Run(c, false); ok {
				if modified {
					b.Modified = true
				}
				if ed.alive && ed.Buf() == b && b.Loc.Row != prevRow {
					b.StoreLine()
				}
				if _, ok := IsModifyingCmd[c.Op.Kind]; ok {
					if modified {
						ed.EndUndoRecord()
					} else {
						ed.CancelUndoRecord()
					}
					ed.lastCmd = c
				} else if c.Op.Kind == Undo {
					// undo is not included in modifying commands
					// it is not usual edit or insert but repeatable
					ed.lastCmd = c
				}
				// reset undo/redo toggle if command is not undo/repeat
				if c.Op.Kind != Undo && c.Op.Kind != Repeat {
					ed.undo = false
				}
				ed.parser.Reset()
			} else {
				ed.Error("Failed to run")
				if _, ok := IsModifyingCmd[c.Op.Kind]; ok {
					ed.CancelUndoRecord()
				}
			}
		}
	case termi.KeyUp:
		ed.Run(CmdPair{Mv: Cmd{Kind: MoveUp, Num: 1}}, false)
	case termi.KeyDown:
		ed.Run(CmdPair{Mv: Cmd{Kind: MoveDown, Num: 1}}, false)
	case termi.KeyRight:
		ed.Run(CmdPair{Mv: Cmd{Kind: MoveRight, Num: 1}}, false)
	case termi.KeyLeft:
		ed.Run(CmdPair{Mv: Cmd{Kind: MoveLeft, Num: 1}}, false)
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
		ed.Run(CmdPair{Mv: Cmd{Kind: MoveUp, Num: 1}}, false)
	case termi.KeyDown:
		ed.Run(CmdPair{Mv: Cmd{Kind: MoveDown, Num: 1}}, false)
	case termi.KeyRight:
		ed.Run(CmdPair{Mv: Cmd{Kind: MoveRight, Num: 1}}, false)
	case termi.KeyLeft:
		ed.Run(CmdPair{Mv: Cmd{Kind: MoveLeft, Num: 1}}, false)
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
			ed.searchs.pattern.Reset()
			ed.mode = ModeCommand
		case termi.RuneEnter, termi.RuneNewline:
			if ed.searchs.pattern.RuneCount() < 1 {
				if ed.searchs.backward {
					ed.Run(CmdPair{Mv: Cmd{Kind: RepeatBackwardSearch}}, false)
				} else {
					ed.Run(CmdPair{Mv: Cmd{Kind: RepeatSearch}}, false)
				}
				return
			}
			re, err := regexp.Compile(ed.searchs.pattern.String())
			if err != nil {
				ed.Ring("%v", err)
				return
			}
			ed.searchs.regexp = re
			ed.searchs.pattern.Reset()
			if ed.searchs.backward {
				ed.Run(CmdPair{Mv: Cmd{Kind: SearchBackward}}, false)
			} else {
				ed.Run(CmdPair{Mv: Cmd{Kind: Search}}, false)
			}
		case termi.RuneBackspace, termi.RuneDelete:
			if !ed.searchs.pattern.RemoveTail() {
				ed.mode = ModeCommand
			}
		default:
			ed.searchs.pattern.WriteRune(key.Rune)
		}
	default:
		ed.Error("Unknown key")
	}
}

func (ed *Editor) Main() error {
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
	return nil
}
