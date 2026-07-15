package editor

import "tea.kareha.org/cup/levi/internal/cmd"

func (ed *Editor) RunOp(c cmd.Pair, replay bool) (bool, bool) {
	ed.Commit()

	switch c.Op.Kind {

	//
	// Insert Commands
	//

	case cmd.Insert:
		return ed.Insert(c.Op.Num, replay), true
	case cmd.InsertAfter:
		return ed.InsertAfter(c.Op.Num, replay), true
	case cmd.InsertAfterIndent:
		return ed.InsertAfterIndent(c.Op.Num, replay), true
	case cmd.InsertAfterEnd:
		return ed.InsertAfterEnd(c.Op.Num, replay), true

	case cmd.InsertLine:
		return ed.InsertLine(c.Op.Num, replay), true
	case cmd.InsertLineAbove:
		return ed.InsertLineAbove(c.Op.Num, replay), true

	case cmd.ChangeRegion:
		start := ed.Buf().Loc
		mv := c.Mv
		if mv.Kind == cmd.MoveByWord {
			mv.Kind = cmd.MoveByChangeWord
		}
		end, ok := ed.RunMove(mv, c.Op.Num)
		if !ok {
			ed.Error("Failed to move")
			return false, false
		}
		attr, ok := cmd.AttrOf[mv.Kind]
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
	case cmd.Subst:
		return ed.Subst(c.Reg, c.Op.Num, replay), true

	case cmd.Overwrite: // unsupported
		return ed.Overwrite(), true

	//
	// Edit Commands
	//

	case cmd.Paste:
		return ed.Paste(c.Reg, c.Op.Num), true
	case cmd.PasteBefore:
		return ed.PasteBefore(c.Reg, c.Op.Num), true

	case cmd.Delete:
		return ed.Delete(c.Reg, c.Op.Num), true
	case cmd.DeleteBefore:
		return ed.DeleteBefore(c.Reg, c.Op.Num), true
	case cmd.DeleteRegion:
		start := ed.Buf().Loc
		mv := c.Mv
		if mv.Kind == cmd.MoveByWord {
			mv.Kind = cmd.MoveByDeleteWord
		}
		end, ok := ed.RunMove(mv, c.Op.Num)
		if !ok {
			ed.Error("Failed to move")
			return false, false
		}
		attr, ok := cmd.AttrOf[mv.Kind]
		if !ok {
			ed.Error("Failed to retrieve move attr")
			return false, false
		}
		if attr.Linewise {
			return ed.DeleteLineRegion(c.Reg, start, end), true
		} else {
			return ed.DeleteRegion(c.Reg, start, end, attr.Inclusive), true
		}

	case cmd.Replace:
		return ed.Replace(c.Op.Rune, c.Op.Num), true
	case cmd.Join:
		return ed.Join(c.Op.Num), true
	case cmd.IndentRegion:
		start := ed.Buf().Loc
		mv := c.Mv
		if mv.Kind == cmd.MoveByWord {
			mv.Kind = cmd.MoveByChangeWord // XXX or cmd.MoveByDeleteWord?
		}
		end, ok := ed.RunMove(mv, c.Op.Num)
		if !ok {
			ed.Error("Failed to move")
			return false, true
		}
		return ed.IndentRegion(start, end), true
	case cmd.OutdentRegion:
		start := ed.Buf().Loc
		mv := c.Mv
		if mv.Kind == cmd.MoveByWord {
			mv.Kind = cmd.MoveByChangeWord // XXX or cmd.MoveByDeleteWord?
		}
		end, ok := ed.RunMove(mv, c.Op.Num)
		if !ok {
			ed.Error("Failed to move")
			return false, true
		}
		return ed.OutdentRegion(start, end), true

	case cmd.Restore:
		return ed.Restore(), true

	//
	// Mark Commands
	//

	case cmd.Mark:
		ed.Mark(c.Op.Rune)
		return false, true

	//
	// Copy Commands
	//

	case cmd.CopyRegion:
		start := ed.Buf().Loc
		end, ok := ed.RunMove(c.Mv, c.Op.Num)
		if !ok {
			ed.Error("Failed to move")
			return false, false
		}
		attr, ok := cmd.AttrOf[c.Mv.Kind]
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

	case cmd.ViewDown:
		ed.ViewDown(c.Op.Num)
		return false, true
	case cmd.ViewUp:
		ed.ViewUp(c.Op.Num)
		return false, true
	case cmd.ViewDownHalf:
		ed.ViewDownHalf(c.Op.Num)
		return false, true
	case cmd.ViewUpHalf:
		ed.ViewUpHalf(c.Op.Num)
		return false, true
	case cmd.ViewDownLine:
		ed.ViewDownLine(c.Op.Num)
		return false, true
	case cmd.ViewUpLine:
		ed.ViewUpLine(c.Op.Num)
		return false, true

	case cmd.ViewToTop:
		ed.ViewToTop()
		return false, true
	case cmd.ViewToMiddle:
		ed.ViewToMiddle()
		return false, true
	case cmd.ViewToBottom:
		ed.ViewToBottom()
		return false, true

	case cmd.Redraw:
		ed.Redraw()
		return false, true

	//
	// Miscellaneous Commands
	//

	case cmd.ShowInfo:
		ed.ShowInfo()
		return false, true
	case cmd.Repeat:
		ed.Repeat(c.Op.Num)
		return false, true
	case cmd.Undo:
		ed.Undo(c.Op.Num, replay)
		return false, true
	case cmd.SaveAndClose:
		ed.SaveAndClose()
		return false, true
	case cmd.Suspend:
		ed.Suspend()
		return false, true

	case cmd.LastBuf:
		ed.LastBuf()
		return false, true
	case cmd.GoToBuf:
		ed.GoToBuf(c.Op.Num)
		return false, true
	case cmd.NextBuf:
		ed.NextBuf()
		return false, true
	case cmd.PrevBuf:
		ed.PrevBuf()
		return false, true

	case cmd.Ring:
		ed.Ring("%s", c.Op.Pat)
		return false, true

	}

	return false, false
}
