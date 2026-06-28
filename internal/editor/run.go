package editor

import (
	"unicode/utf8"
)

func (ed *Editor) Run(cp CmdPair, replay bool) bool {
	ed.Commit()

	switch cp.Main.Kind {
	case CmdInvalid:
		ed.Notice("Not a levi command [%s]", ed.parser.String())
		return true
	}

	if meta, ok := MoveMetas[cp.Main.Kind]; ok {
		if loc, ok := ed.RunMove(cp.Main); ok {
			b := ed.Buf()
			if meta.Linewise {
				if meta.FreeCol {
					line := b.Line(loc.Row)
					rc := utf8.RuneCountInString(line)
					if b.VirtCol < rc {
						loc.Col = b.VirtCol
					} else {
						loc.Col = max(rc-1, 0)
					}
				}
			} else {
				loc = b.ConfineInclusive(loc)
				b.VirtCol = loc.Col
			}
			b.Loc = loc
			if b.Loc.Col < b.ViewLoc.Col {
				b.ViewLoc.Col = 0
			}
			if meta.Locate {
				ed.Locate()
			}
		}
		return true
	}

	switch cp.Main.Kind {

	case CmdMarkSet:
		ed.MarkSet(cp.Main.Letter)
		return true

	case CmdViewDown:
		ed.ViewDown(cp.Main.Num)
		return true
	case CmdViewUp:
		ed.ViewUp(cp.Main.Num)
		return true
	case CmdViewDownHalf:
		ed.ViewDownHalf(cp.Main.Num)
		return true
	case CmdViewUpHalf:
		ed.ViewUpHalf(cp.Main.Num)
		return true
	case CmdViewDownLine:
		ed.ViewDownLine(cp.Main.Num)
		return true
	case CmdViewUpLine:
		ed.ViewUpLine(cp.Main.Num)
		return true

	case CmdViewToTop:
		ed.ViewToTop()
		return true
	case CmdViewToMiddle:
		ed.ViewToMiddle()
		return true
	case CmdViewToBottom:
		ed.ViewToBottom()
		return true

	case CmdViewRedraw:
		ed.ViewRedraw()
		return true

	case CmdInsertBefore:
		ed.InsertBefore(cp.Main.Num, replay)
		return true
	case CmdInsertAfter:
		ed.InsertAfter(cp.Main.Num, replay)
		return true
	case CmdInsertBeforeNonBlank:
		ed.InsertBeforeNonBlank(cp.Main.Num, replay)
		return true
	case CmdInsertAfterEnd:
		ed.InsertAfterEnd(cp.Main.Num, replay)
		return true
	case CmdInsertOverwrite:
		ed.InsertOverwrite(cp.Main.Num, replay)
		return true

	case CmdInsertOpenBelow:
		ed.InsertOpenBelow(cp.Main.Num, replay)
		return true
	case CmdInsertOpenAbove:
		ed.InsertOpenAbove(cp.Main.Num, replay)
		return true

	case CmdOpCopyLine:
		ed.OpCopyLine(cp.Reg, cp.Main.Num)
		return true
	case CmdOpCopyRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(cp.Sub)
		if !ok {
			ed.Error("Failed to move")
			return false
		}
		meta, ok := MoveMetas[cp.Sub.Kind]
		if !ok {
			ed.Error("Failed to retrieve move meta")
			return false
		}
		if meta.Linewise {
			ed.OpCopyLineRegion(cp.Reg, start, end)
		} else {
			ed.OpCopyRegion(cp.Reg, start, end, meta.Inclusive)
		}
		return true
	case CmdOpCopyWord:
		ed.OpCopyWord(cp.Reg, cp.Main.Num)
		return true
	case CmdOpCopyToEnd:
		ed.OpCopyToEnd(cp.Reg, cp.Main.Num)
		return true

	case CmdOpPaste:
		ed.OpPaste(cp.Reg, cp.Main.Num)
		return true
	case CmdOpPasteBefore:
		ed.OpPasteBefore(cp.Reg, cp.Main.Num)
		return true

	case CmdOpDelete:
		ed.OpDelete(cp.Reg, cp.Main.Num)
		return true
	case CmdOpDeleteBefore:
		ed.OpDeleteBefore(cp.Reg, cp.Main.Num)
		return true
	case CmdOpDeleteLine:
		ed.OpDeleteLine(cp.Reg, cp.Main.Num)
		return true
	case CmdOpDeleteRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(cp.Sub)
		if !ok {
			ed.Error("Failed to move")
			return false
		}
		meta, ok := MoveMetas[cp.Sub.Kind]
		if !ok {
			ed.Error("Failed to retrieve move meta")
			return false
		}
		if meta.Linewise {
			ed.OpDeleteLineRegion(cp.Reg, start, end)
		} else {
			ed.OpDeleteRegion(cp.Reg, start, end, meta.Inclusive)
		}
		return true
	case CmdOpDeleteWord:
		ed.OpDeleteWord(cp.Reg, cp.Main.Num)
		return true
	case CmdOpDeleteToEnd:
		ed.OpDeleteToEnd(cp.Reg, cp.Main.Num)
		return true

	case CmdOpChangeLine:
		ed.OpChangeLine(cp.Reg, cp.Main.Num, replay)
		return true
	case CmdOpChangeRegion:
		start := ed.Buf().Loc
		cmd := cp.Sub
		if cmd.Kind == CmdMoveByWord {
			cmd.Kind = CmdMoveByWordEx
		}
		end, ok := ed.RunMove(cmd)
		if !ok {
			ed.Error("Failed to move")
			return false
		}
		meta, ok := MoveMetas[cmd.Kind]
		if !ok {
			ed.Error("Failed to retrieve move meta")
			return false
		}
		if meta.Linewise {
			ed.OpChangeLineRegion(cp.Reg, start, end, replay)
		} else {
			ed.OpChangeRegion(cp.Reg, start, end, meta.Inclusive, replay)
		}
		return true
	case CmdOpChangeWord:
		ed.OpChangeWord(cp.Reg, cp.Main.Num, replay)
		return true
	case CmdOpChangeToEnd:
		ed.OpChangeToEnd(cp.Reg, cp.Main.Num, replay)
		return true
	case CmdOpSubst:
		ed.OpSubst(cp.Reg, cp.Main.Num, replay)
		return true
	case CmdOpSubstLine:
		ed.OpSubstLine(cp.Reg, cp.Main.Num, replay)
		return true

	case CmdEditReplace:
		ed.EditReplace(cp.Main.Letter, cp.Main.Num, replay)
		return true
	case CmdEditJoin:
		ed.EditJoin(cp.Main.Num)
		return true
	case CmdEditIndent:
		ed.EditIndent(cp.Main.Num)
		return true
	case CmdEditOutdent:
		ed.EditOutdent(cp.Main.Num)
		return true
	case CmdEditIndentRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(cp.Sub)
		if !ok {
			ed.Error("Failed to move")
			return true
		}
		ed.EditIndentRegion(start, end)
		return true
	case CmdEditOutdentRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(cp.Sub)
		if !ok {
			ed.Error("Failed to move")
			return true
		}
		ed.EditOutdentRegion(start, end)
		return true

	case CmdMiscShowInfo:
		ed.MiscShowInfo()
		return true
	case CmdMiscRepeat:
		ed.MiscRepeat(cp.Main.Num)
		return true
	case CmdMiscUndo:
		ed.MiscUndo(cp.Main.Num, replay)
		return true
	case CmdMiscRestore:
		ed.MiscRestore()
		return true
	case CmdMiscSaveAndQuit:
		ed.MiscSaveAndQuit()
		return true
	case CmdMiscSuspend:
		ed.MiscSuspend()
		return true

	}

	return false
}
