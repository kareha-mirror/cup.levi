package editor

func (ed *Editor) Run(cp CmdPair, replay bool) bool {
	ed.Commit()

	if cp.Op.Kind == InvalidCmd {
		if meta, ok := MoveAttrs[cp.Mv.Kind]; ok {
			if loc, ok := ed.RunMove(cp.Mv); ok {
				b := ed.Buf()
				if meta.Linewise {
					if meta.FreeCol {
						loc.Col = b.ConfineFreeColInclusive(loc.Row)
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

		ed.Notice("Not a levi command [%s]", ed.parser.String())
		return true
	}

	switch cp.Op.Kind {

	case MarkSet:
		ed.MarkSet(cp.Op.Ltr)
		return true

	case ViewDown:
		ed.ViewDown(cp.Op.Num)
		return true
	case ViewUp:
		ed.ViewUp(cp.Op.Num)
		return true
	case ViewDownHalf:
		ed.ViewDownHalf(cp.Op.Num)
		return true
	case ViewUpHalf:
		ed.ViewUpHalf(cp.Op.Num)
		return true
	case ViewDownLine:
		ed.ViewDownLine(cp.Op.Num)
		return true
	case ViewUpLine:
		ed.ViewUpLine(cp.Op.Num)
		return true

	case ViewToTop:
		ed.ViewToTop()
		return true
	case ViewToMiddle:
		ed.ViewToMiddle()
		return true
	case ViewToBottom:
		ed.ViewToBottom()
		return true

	case Redraw:
		ed.Redraw()
		return true

	case InsertBefore:
		ed.InsertBefore(cp.Op.Num, replay)
		return true
	case InsertAfter:
		ed.InsertAfter(cp.Op.Num, replay)
		return true
	case InsertBeforeNonBlank:
		ed.InsertBeforeNonBlank(cp.Op.Num, replay)
		return true
	case InsertAfterEnd:
		ed.InsertAfterEnd(cp.Op.Num, replay)
		return true
	case Overwrite:
		ed.Overwrite(cp.Op.Num, replay)
		return true

	case OpenBelow:
		ed.OpenBelow(cp.Op.Num, replay)
		return true
	case OpenAbove:
		ed.OpenAbove(cp.Op.Num, replay)
		return true

	case CopyLine:
		ed.CopyLine(cp.Reg, cp.Op.Num)
		return true
	case CopyRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(cp.Mv)
		if !ok {
			ed.Error("Failed to move")
			return false
		}
		meta, ok := MoveAttrs[cp.Mv.Kind]
		if !ok {
			ed.Error("Failed to retrieve move meta")
			return false
		}
		if meta.Linewise {
			ed.CopyLineRegion(cp.Reg, start, end)
		} else {
			ed.CopyRegion(cp.Reg, start, end, meta.Inclusive)
		}
		return true
	case CopyWord:
		ed.CopyWord(cp.Reg, cp.Op.Num)
		return true
	case CopyToEnd:
		ed.CopyToEnd(cp.Reg, cp.Op.Num)
		return true

	case Paste:
		ed.Paste(cp.Reg, cp.Op.Num)
		return true
	case PasteBefore:
		ed.PasteBefore(cp.Reg, cp.Op.Num)
		return true

	case Delete:
		ed.Delete(cp.Reg, cp.Op.Num)
		return true
	case DeleteBefore:
		ed.DeleteBefore(cp.Reg, cp.Op.Num)
		return true
	case DeleteLine:
		ed.DeleteLine(cp.Reg, cp.Op.Num)
		return true
	case DeleteRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(cp.Mv)
		if !ok {
			ed.Error("Failed to move")
			return false
		}
		meta, ok := MoveAttrs[cp.Mv.Kind]
		if !ok {
			ed.Error("Failed to retrieve move meta")
			return false
		}
		if meta.Linewise {
			ed.DeleteLineRegion(cp.Reg, start, end)
		} else {
			ed.DeleteRegion(cp.Reg, start, end, meta.Inclusive)
		}
		return true
	case DeleteWord:
		ed.DeleteWord(cp.Reg, cp.Op.Num)
		return true
	case DeleteToEnd:
		ed.DeleteToEnd(cp.Reg, cp.Op.Num)
		return true

	case ChangeLine:
		ed.ChangeLine(cp.Reg, cp.Op.Num, replay)
		return true
	case ChangeRegion:
		start := ed.Buf().Loc
		cmd := cp.Mv
		if cmd.Kind == MoveByWord {
			cmd.Kind = MoveByChangeWord
		}
		end, ok := ed.RunMove(cmd)
		if !ok {
			ed.Error("Failed to move")
			return false
		}
		meta, ok := MoveAttrs[cmd.Kind]
		if !ok {
			ed.Error("Failed to retrieve move meta")
			return false
		}
		if meta.Linewise {
			ed.ChangeLineRegion(cp.Reg, start, end, replay)
		} else {
			ed.ChangeRegion(cp.Reg, start, end, meta.Inclusive, replay)
		}
		return true
	case ChangeWord:
		ed.ChangeWord(cp.Reg, cp.Op.Num, replay)
		return true
	case ChangeToEnd:
		ed.ChangeToEnd(cp.Reg, cp.Op.Num, replay)
		return true
	case Subst:
		ed.Subst(cp.Reg, cp.Op.Num, replay)
		return true
	case SubstLine:
		ed.SubstLine(cp.Reg, cp.Op.Num, replay)
		return true

	case Replace:
		ed.Replace(cp.Op.Ltr, cp.Op.Num, replay)
		return true
	case Join:
		ed.Join(cp.Op.Num)
		return true
	case Indent:
		ed.Indent(cp.Op.Num)
		return true
	case Outdent:
		ed.Outdent(cp.Op.Num)
		return true
	case IndentRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(cp.Mv)
		if !ok {
			ed.Error("Failed to move")
			return true
		}
		ed.IndentRegion(start, end)
		return true
	case OutdentRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(cp.Mv)
		if !ok {
			ed.Error("Failed to move")
			return true
		}
		ed.OutdentRegion(start, end)
		return true

	case ShowInfo:
		ed.ShowInfo()
		return true
	case Repeat:
		ed.Repeat(cp.Op.Num)
		return true
	case Undo:
		ed.Undo(cp.Op.Num, replay)
		return true
	case Restore:
		ed.Restore()
		return true
	case SaveAndQuit:
		ed.SaveAndQuit()
		return true
	case Suspend:
		ed.Suspend()
		return true

	}

	return false
}
