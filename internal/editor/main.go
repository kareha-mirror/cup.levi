package editor

import (
	"fmt"
	"regexp"

	"tea.kareha.org/cup/termi"

	"tea.kareha.org/cup/levi/internal/cmd"
	"tea.kareha.org/cup/levi/internal/prompt"
)

func (ed *Editor) MainCommand(key termi.Key) {
	switch key.Kind {
	case termi.KeyRune:
		if key.Rune == termi.RuneEscape {
			if ed.cmdInp.String() == "" {
				ed.Notice("Already in command mode")
			}
			ed.Reset()
			return
		}

		if ed.cmdInp.String() == "" {
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
			if !ed.cmdInp.Backspace() {
				ed.Notice("No tokens left")
				return
			}
		} else {
			ed.cmdInp.WriteRune(key.Rune)
		}

		a := cmd.Parse(ed.cmdInp.String())
		ed.args = a
		c, ok := a.Parse()
		ed.cmdOk = ok
		if ok {
			b := ed.Buf()
			prevRow := b.Loc.Row
			if _, ok := cmd.IsModifying[c.Op.Kind]; ok {
				ed.BeginUndoRecord()
			}
			if modified, ok := ed.Run(c, false); ok {
				if modified {
					b.Modified = true
				}
				if ed.alive && ed.Buf() == b && b.Loc.Row != prevRow {
					b.StoreLine()
				}
				if _, ok := cmd.IsInsert[c.Op.Kind]; ok {
					ed.lastCmd = c
				} else if _, ok := cmd.IsEdit[c.Op.Kind]; ok {
					if modified {
						ed.EndUndoRecord()
						ed.lastCmd = c
					} else {
						ed.CancelUndoRecord()
					}
				} else if c.Op.Kind == cmd.Undo {
					// undo is not included in modifying commands
					// it is not usual edit or insert but repeatable
					ed.lastCmd = c
				}
				// reset undo/redo toggle if command is not undo/repeat
				if c.Op.Kind != cmd.Undo && c.Op.Kind != cmd.Repeat {
					ed.undo = false
				}
				// reset buffer select mode
				if _, ok := cmd.IsBufMove[c.Op.Kind]; !ok {
					ed.bufMove = false
				}
				ed.cmdInp.Reset()
			} else {
				ed.Error("Failed to run")
				if _, ok := cmd.IsModifying[c.Op.Kind]; ok {
					ed.CancelUndoRecord()
				}
			}
		}
	case termi.KeyUp:
		ed.Run(cmd.Pair{Mv: cmd.Cmd{Kind: cmd.MoveUp, Num: 1}}, false)
	case termi.KeyDown:
		ed.Run(cmd.Pair{Mv: cmd.Cmd{Kind: cmd.MoveDown, Num: 1}}, false)
	case termi.KeyRight:
		ed.Run(cmd.Pair{Mv: cmd.Cmd{Kind: cmd.MoveRight, Num: 1}}, false)
	case termi.KeyLeft:
		ed.Run(cmd.Pair{Mv: cmd.Cmd{Kind: cmd.MoveLeft, Num: 1}}, false)
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
		ed.Run(cmd.Pair{Mv: cmd.Cmd{Kind: cmd.MoveUp, Num: 1}}, false)
	case termi.KeyDown:
		ed.Run(cmd.Pair{Mv: cmd.Cmd{Kind: cmd.MoveDown, Num: 1}}, false)
	case termi.KeyRight:
		ed.Run(cmd.Pair{Mv: cmd.Cmd{Kind: cmd.MoveRight, Num: 1}}, false)
	case termi.KeyLeft:
		ed.Run(cmd.Pair{Mv: cmd.Cmd{Kind: cmd.MoveLeft, Num: 1}}, false)
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
			c, ok := prompt.Parse(ed.prompt.String())
			if ok {
				ed.prompt.Reset()
			}
			ok = ed.RunPrompt(c)
			if !ok {
				//ed.Error("Prompt command failed")
			}
			// reset buffer select mode
			if _, ok := prompt.IsBufMove[c.Kind]; !ok {
				ed.bufMove = false
			}
			ed.mode = ModeCommand
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
					ed.Run(cmd.Pair{Mv: cmd.Cmd{Kind: cmd.RepeatBackwardSearch}}, false)
				} else {
					ed.Run(cmd.Pair{Mv: cmd.Cmd{Kind: cmd.RepeatSearch}}, false)
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
				ed.Run(cmd.Pair{Mv: cmd.Cmd{Kind: cmd.SearchBackward}}, false)
			} else {
				ed.Run(cmd.Pair{Mv: cmd.Cmd{Kind: cmd.Search}}, false)
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
