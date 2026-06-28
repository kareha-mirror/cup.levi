package editor

import (
	"unicode/utf8"
)

func (ed *Editor) Run(c Cmd, replay bool) bool {
	ed.Commit()

	switch c.Kind {
	case CmdInvalid:
		ed.Notice("Not a levi command [%s]", ed.parser.String())
		return true
	}

	if meta, ok := MoveMetas[c.Kind]; ok {
		if loc, ok := ed.RunMove(c); ok {
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
		ed.OpCopyLine(c.Reg, c.Num)
		return true
	case CmdOpCopyRegion:
		ed.OpCopyRegion(c.Reg, c.Start, c.End, c.Inclusive)
		return true
	case CmdOpCopyLineRegion:
		ed.OpCopyLineRegion(c.Reg, c.Start, c.End)
		return true
	case CmdOpCopyWord:
		ed.OpCopyWord(c.Reg, c.Num)
		return true
	case CmdOpCopyToEnd:
		ed.OpCopyToEnd(c.Reg, c.Num)
		return true

	case CmdOpPaste:
		ed.OpPaste(c.Reg, c.Num)
		return true
	case CmdOpPasteBefore:
		ed.OpPasteBefore(c.Reg, c.Num)
		return true

	case CmdOpDelete:
		ed.OpDelete(c.Reg, c.Num)
		return true
	case CmdOpDeleteBefore:
		ed.OpDeleteBefore(c.Reg, c.Num)
		return true
	case CmdOpDeleteLine:
		ed.OpDeleteLine(c.Reg, c.Num)
		return true
	case CmdOpDeleteRegion:
		ed.OpDeleteRegion(c.Reg, c.Start, c.End, c.Inclusive)
		return true
	case CmdOpDeleteLineRegion:
		ed.OpDeleteLineRegion(c.Reg, c.Start, c.End)
		return true
	case CmdOpDeleteWord:
		ed.OpDeleteWord(c.Reg, c.Num)
		return true
	case CmdOpDeleteToEnd:
		ed.OpDeleteToEnd(c.Reg, c.Num)
		return true

	case CmdOpChangeLine:
		ed.OpChangeLine(c.Reg, c.Num, replay)
		return true
	case CmdOpChangeRegion:
		ed.OpChangeRegion(c.Reg, c.Start, c.End, c.Inclusive, replay)
		return true
	case CmdOpChangeLineRegion:
		ed.OpChangeLineRegion(c.Reg, c.Start, c.End, replay)
		return true
	case CmdOpChangeWord:
		ed.OpChangeWord(c.Reg, c.Num, replay)
		return true
	case CmdOpChangeToEnd:
		ed.OpChangeToEnd(c.Reg, c.Num, replay)
		return true
	case CmdOpSubst:
		ed.OpSubst(c.Reg, c.Num, replay)
		return true
	case CmdOpSubstLine:
		ed.OpSubstLine(c.Reg, c.Num, replay)
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
