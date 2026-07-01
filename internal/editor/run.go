package editor

func (ed *Editor) Run(c CmdPair, replay bool) (bool, bool) {
	ed.Commit()

	// * motion commands are delegated to RunMove
	// * other commands are dispatched in this method
	// * compound commands also refer to RunMove

	//
	// Motion Commands
	//

	if c.Op.Kind == InvalidCmd {
		if attr, ok := MoveAttrs[c.Mv.Kind]; ok {
			if loc, ok := ed.RunMove(c.Mv, 1, false); ok {
				b := ed.Buf()
				if attr.Linewise {
					if attr.FreeCol {
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
				if attr.Locate {
					ed.Locate()
				}
			}
			return false, true
		}

		ed.Notice("Not a levi command [%s]", ed.parser.String())
		return false, true
	}

	switch c.Op.Kind {

	//
	// Insert Commands
	//

	case Insert:
		return ed.Insert(c.Op.Num, replay), true
	case InsertAfter:
		return ed.InsertAfter(c.Op.Num, replay), true
	case InsertAfterIndent:
		return ed.InsertAfterIndent(c.Op.Num, replay), true
	case InsertAfterEnd:
		return ed.InsertAfterEnd(c.Op.Num, replay), true
	case Overwrite:
		return ed.Overwrite(c.Op.Num, replay), true

	case OpenBelow:
		return ed.OpenBelow(c.Op.Num, replay), true
	case OpenAbove:
		return ed.OpenAbove(c.Op.Num, replay), true

	case ChangeRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(c.Mv, c.Op.Num, true)
		if !ok {
			ed.Error("Failed to move")
			return false, false
		}
		attr, ok := MoveAttrs[c.Mv.Kind]
		if !ok {
			ed.Error("Failed to retrieve move attr")
			return false, false
		}
		if attr.Linewise {
			return ed.ChangeLineRegion(c.Reg, start, end, replay), true
		} else {
			return ed.ChangeRegion(
				c.Reg, start, end, attr.Inclusive, replay,
			), true
		}
	case Subst:
		return ed.Subst(c.Reg, c.Op.Num, replay), true

	//
	// Edit Commands
	//

	case Paste:
		return ed.Paste(c.Reg, c.Op.Num), true
	case PasteBefore:
		return ed.PasteBefore(c.Reg, c.Op.Num), true

	case Delete:
		return ed.Delete(c.Reg, c.Op.Num), true
	case DeleteBefore:
		return ed.DeleteBefore(c.Reg, c.Op.Num), true
	case DeleteRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(c.Mv, c.Op.Num, true)
		if !ok {
			ed.Error("Failed to move")
			return false, false
		}
		attr, ok := MoveAttrs[c.Mv.Kind]
		if !ok {
			ed.Error("Failed to retrieve move attr")
			return false, false
		}
		if attr.Linewise {
			return ed.DeleteLineRegion(c.Reg, start, end), true
		} else {
			return ed.DeleteRegion(c.Reg, start, end, attr.Inclusive), true
		}

	case Replace:
		return ed.Replace(c.Op.Rune, c.Op.Num), true
	case Join:
		return ed.Join(c.Op.Num), true
	case IndentRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(c.Mv, c.Op.Num, true)
		if !ok {
			ed.Error("Failed to move")
			return false, true
		}
		return ed.IndentRegion(start, end), true
	case OutdentRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(c.Mv, c.Op.Num, true)
		if !ok {
			ed.Error("Failed to move")
			return false, true
		}
		return ed.OutdentRegion(start, end), true

	case Restore:
		return ed.Restore(), true

	//
	// Mark Commands
	//

	case Mark:
		ed.Mark(c.Op.Rune)
		return false, true

	//
	// Copy Commands
	//

	case CopyRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(c.Mv, c.Op.Num, true)
		if !ok {
			ed.Error("Failed to move")
			return false, false
		}
		attr, ok := MoveAttrs[c.Mv.Kind]
		if !ok {
			ed.Error("Failed to retrieve move attr")
			return false, false
		}
		if attr.Linewise {
			ed.CopyLineRegion(c.Reg, start, end)
		} else {
			ed.CopyRegion(c.Reg, start, end, attr.Inclusive)
		}
		return false, true

	//
	// View Commands
	//

	case ViewDown:
		ed.ViewDown(c.Op.Num)
		return false, true
	case ViewUp:
		ed.ViewUp(c.Op.Num)
		return false, true
	case ViewDownHalf:
		ed.ViewDownHalf(c.Op.Num)
		return false, true
	case ViewUpHalf:
		ed.ViewUpHalf(c.Op.Num)
		return false, true
	case ViewDownLine:
		ed.ViewDownLine(c.Op.Num)
		return false, true
	case ViewUpLine:
		ed.ViewUpLine(c.Op.Num)
		return false, true

	case ViewToTop:
		ed.ViewToTop()
		return false, true
	case ViewToMiddle:
		ed.ViewToMiddle()
		return false, true
	case ViewToBottom:
		ed.ViewToBottom()
		return false, true

	case Redraw:
		ed.Redraw()
		return false, true

	//
	// Miscellaneous Commands
	//

	case ShowInfo:
		ed.ShowInfo()
		return false, true
	case Repeat:
		ed.Repeat(c.Op.Num)
		return false, true
	case Undo:
		ed.Undo(c.Op.Num, replay)
		return false, true
	case SaveAndClose:
		ed.SaveAndClose()
		return false, true
	case Suspend:
		ed.Suspend()
		return false, true

	}

	return false, false
}
