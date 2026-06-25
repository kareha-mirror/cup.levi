package editor

import (
	"unicode/utf8"
)

func (ed *Editor) Run(c Cmd, replay bool) bool {
	switch c.Kind {
	case CmdInvalid:
		ed.Ring("not (yet) a vi command [" + ed.parser.String() + "]")
		ed.parser.Clear()
		return true
	}

	if _, ok := MoveCmds[c.Kind]; ok {
		if dest, ok := ed.RunMove(c); ok {
			b := ed.Buf()
			if dest.Linewise {
				if dest.FreeCol {
					line := b.Line(dest.Loc.Row)
					rc := utf8.RuneCountInString(line)
					if b.VirtCol < rc {
						dest.Loc.Col = b.VirtCol
					} else {
						dest.Loc.Col = max(rc-1, 0)
					}
				}
			} else {
				b.VirtCol = dest.Loc.Col
			}
			b.Loc = dest.Loc
			if b.Loc.Col < b.ViewLoc.Col {
				b.ViewLoc.Col = 0
			}
			if dest.Locate {
				ed.Locate()
			}
		}
		return true
	}

	switch c.Kind {
	case CmdMarkSet:
		ed.MarkSet(c.Letter)
		return true

	case CmdViewDown:
		ed.ViewDown(c.Num)
		return true
	case CmdViewUp:
		ed.ViewUp(c.Num)
		return true
	case CmdViewDownHalf:
		ed.ViewDownHalf(c.Num)
		return true
	case CmdViewUpHalf:
		ed.ViewUpHalf(c.Num)
		return true
	case CmdViewDownLine:
		ed.ViewDownLine(c.Num)
		return true
	case CmdViewUpLine:
		ed.ViewUpLine(c.Num)
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
		ed.InsertBefore(c.Num, replay)
		return true
	case CmdInsertAfter:
		ed.InsertAfter(c.Num, replay)
		return true
	case CmdInsertBeforeNonBlank:
		ed.InsertBeforeNonBlank(c.Num, replay)
		return true
	case CmdInsertAfterEnd:
		ed.InsertAfterEnd(c.Num, replay)
		return true
	case CmdInsertOverwrite:
		ed.InsertOverwrite(c.Num, replay)
		return true

	case CmdInsertOpenBelow:
		ed.InsertOpenBelow(c.Num, replay)
		return true
	case CmdInsertOpenAbove:
		ed.InsertOpenAbove(c.Num, replay)
		return true

	case CmdOpCopyLine:
		ed.OpCopyLine(c.Num)
		return true
	case CmdOpCopyRegion:
		ed.OpCopyRegion(c.Start, c.End, c.Inclusive)
		return true
	case CmdOpCopyLineRegion:
		ed.OpCopyLineRegion(c.StartRow, c.EndRow)
		return true
	case CmdOpCopyWord:
		ed.OpCopyWord(c.Num)
		return true
	case CmdOpCopyToEnd:
		ed.OpCopyToEnd(c.Num)
		return true
	case CmdOpCopyLineIntoReg:
		ed.OpCopyLineIntoReg(c.Reg, c.Num)
		return true

	case CmdOpPaste:
		ed.OpPaste(c.Num)
		return true
	case CmdOpPasteBefore:
		ed.OpPasteBefore(c.Num)
		return true
	case CmdOpPasteFromReg:
		ed.OpPasteFromReg(c.Reg, c.Num)
		return true

	case CmdOpDelete:
		ed.OpDelete(c.Num)
		return true
	case CmdOpDeleteBefore:
		ed.OpDeleteBefore(c.Num)
		return true
	case CmdOpDeleteLine:
		ed.OpDeleteLine(c.Num)
		return true
	case CmdOpDeleteRegion:
		ed.OpDeleteRegion(c.Start, c.End, c.Inclusive)
		return true
	case CmdOpDeleteLineRegion:
		ed.OpDeleteLineRegion(c.StartRow, c.EndRow)
		return true
	case CmdOpDeleteWord:
		ed.OpDeleteWord(c.Num)
		return true
	case CmdOpDeleteToEnd:
		ed.OpDeleteToEnd(c.Num)
		return true

	case CmdOpChangeLine:
		ed.OpChangeLine(c.Num, replay)
		return true
	case CmdOpChangeRegion:
		ed.OpChangeRegion(c.Start, c.End, c.Inclusive, replay)
		return true
	case CmdOpChangeLineRegion:
		ed.OpChangeLineRegion(c.StartRow, c.EndRow, replay)
		return true
	case CmdOpChangeWord:
		ed.OpChangeWord(c.Num, replay)
		return true
	case CmdOpChangeToEnd:
		ed.OpChangeToEnd(c.Num, replay)
		return true
	case CmdOpSubst:
		ed.OpSubst(c.Num, replay)
		return true
	case CmdOpSubstLine:
		ed.OpSubstLine(c.Num, replay)
		return true

	case CmdEditReplace:
		ed.EditReplace(c.Letter, c.Num, replay)
		return true
	case CmdEditJoin:
		ed.EditJoin(c.Num)
		return true
	case CmdEditIndent:
		ed.EditIndent(c.Num)
		return true
	case CmdEditOutdent:
		ed.EditOutdent(c.Num)
		return true
	case CmdEditIndentRegion:
		ed.EditIndentRegion(c.Start, c.End)
		return true
	case CmdEditOutdentRegion:
		ed.EditOutdentRegion(c.Start, c.End)
		return true

	case CmdMiscShowInfo:
		ed.MiscShowInfo()
		return true
	case CmdMiscRepeat:
		ed.MiscRepeat(c.Num)
		return true
	case CmdMiscUndo:
		ed.MiscUndo(c.Num, replay)
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
