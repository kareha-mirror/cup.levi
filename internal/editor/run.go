package editor

func (ed *Editor) Run(c Cmd) bool {
	switch c.Kind {
	case CmdInvalid:
		ed.Ring("not (yet) a vi command [" + ed.parser.String() + "]")
		ed.parser.Clear()
		return true
	}

	switch c.Kind {
	case CmdMoveLeft:
		ed.MoveLeft(c.Num)
		return true
	case CmdMoveDown:
		ed.MoveDown(c.Num)
		return true
	case CmdMoveUp:
		ed.MoveUp(c.Num)
		return true
	case CmdMoveRight:
		ed.MoveRight(c.Num)
		return true

	case CmdMoveToStart:
		ed.MoveToStart()
		return true
	case CmdMoveToEnd:
		ed.MoveToEnd()
		return true
	case CmdMoveToNonBlank:
		ed.MoveToNonBlank()
		return true
	case CmdMoveToColumn:
		ed.MoveToColumn(c.Num)
		return true

	case CmdMoveByWord:
		ed.MoveByWord(c.Num)
		return true
	case CmdMoveBackwardByWord:
		ed.MoveBackwardByWord(c.Num)
		return true
	case CmdMoveToEndOfWord:
		ed.MoveToEndOfWord(c.Num)
		return true
	case CmdMoveByLooseWord:
		ed.MoveByLooseWord(c.Num)
		return true
	case CmdMoveBackwardByLooseWord:
		ed.MoveBackwardByLooseWord(c.Num)
		return true
	case CmdMoveToEndOfLooseWord:
		ed.MoveToEndOfLooseWord(c.Num)
		return true

	case CmdMoveToNonBlankOfNextLine:
		ed.MoveToNonBlankOfNextLine(c.Num)
		return true
	case CmdMoveToNonBlankOfPrevLine:
		ed.MoveToNonBlankOfPrevLine(c.Num)
		return true
	case CmdMoveToLastLine:
		ed.MoveToLastLine()
		return true
	case CmdMoveToLine:
		ed.MoveToLine(c.Num)
		return true

	case CmdMoveBySentence:
		ed.MoveBySentence(c.Num)
		return true
	case CmdMoveBackwardBySentence:
		ed.MoveBackwardBySentence(c.Num)
		return true
	case CmdMoveByParagraph:
		ed.MoveByParagraph(c.Num)
		return true
	case CmdMoveBackwardByParagraph:
		ed.MoveBackwardByParagraph(c.Num)
		return true
	case CmdMoveBySection:
		ed.MoveBySection(c.Num)
		return true
	case CmdMoveBackwardBySection:
		ed.MoveBackwardBySection(c.Num)
		return true

	case CmdMoveToTopOfView:
		ed.MoveToTopOfView()
		return true
	case CmdMoveToMiddleOfView:
		ed.MoveToMiddleOfView()
		return true
	case CmdMoveToBottomOfView:
		ed.MoveToBottomOfView()
		return true
	case CmdMoveToBelowTopOfView:
		ed.MoveToBelowTopOfView(c.Num)
		return true
	case CmdMoveToAboveBottomOfView:
		ed.MoveToAboveBottomOfView(c.Num)
		return true
	}

	switch c.Kind {
	case CmdMarkSet:
		ed.MarkSet(c.Letter)
		return true
	case CmdMarkMoveTo:
		ed.MarkMoveTo(c.Letter)
		return true
	case CmdMarkMoveToLine:
		ed.MarkMoveToLine(c.Letter)
		return true

	case CmdMarkBack:
		ed.MarkBack()
		return true
	case CmdMarkBackToLine:
		ed.MarkBackToLine()
		return true
	}

	switch c.Kind {
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
	}

	switch c.Kind {
	case CmdSearchForward:
		ed.SearchForward(c.Pat)
		return true
	case CmdSearchBackward:
		ed.SearchBackward(c.Pat)
		return true
	case CmdSearchNextMatch:
		ed.SearchNextMatch()
		return true
	case CmdSearchPrevMatch:
		ed.SearchPrevMatch()
		return true
	case CmdSearchRepeatForward:
		ed.SearchRepeatForward()
		return true
	case CmdSearchRepeatBackward:
		ed.SearchRepeatBackward()
		return true
	}

	switch c.Kind {
	case CmdFindForward:
		ed.FindForward(c.Letter, c.Num)
		return true
	case CmdFindBackward:
		ed.FindBackward(c.Letter, c.Num)
		return true
	case CmdFindBeforeForward:
		ed.FindBeforeForward(c.Letter, c.Num)
		return true
	case CmdFindBeforeBackward:
		ed.FindBeforeBackward(c.Letter, c.Num)
		return true
	case CmdFindNextMatch:
		ed.FindNextMatch(c.Num)
		return true
	case CmdFindPrevMatch:
		ed.FindPrevMatch(c.Num)
		return true
	}

	switch c.Kind {
	case CmdInsertBefore:
		ed.InsertBefore(c.Num)
		return true
	case CmdInsertAfter:
		ed.InsertAfter(c.Num)
		return true
	case CmdInsertBeforeNonBlank:
		ed.InsertBeforeNonBlank(c.Num)
		return true
	case CmdInsertAfterEnd:
		ed.InsertAfterEnd(c.Num)
		return true
	case CmdInsertOverwrite:
		ed.InsertOverwrite(c.Num)
		return true

	case CmdInsertOpenBelow:
		ed.InsertOpenBelow(c.Num)
		return true
	case CmdInsertOpenAbove:
		ed.InsertOpenAbove(c.Num)
		return true

	case CmdOpCopyLine:
		ed.OpCopyLine(c.Num)
		return true
	case CmdOpCopyRegion:
		ed.OpCopyRegion(c.Start, c.End)
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
		ed.OpDeleteRegion(c.Start, c.End)
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
		ed.OpChangeLine(c.Num)
		return true
	case CmdOpChangeRegion:
		ed.OpChangeRegion(c.Start, c.End)
		return true
	case CmdOpChangeLineRegion:
		ed.OpChangeLineRegion(c.StartRow, c.EndRow)
		return true
	case CmdOpChangeWord:
		ed.OpChangeWord(c.Num)
		return true
	case CmdOpChangeToEnd:
		ed.OpChangeToEnd(c.Num)
		return true
	case CmdOpSubst:
		ed.OpSubst(c.Num)
		return true
	case CmdOpSubstLine:
		ed.OpSubstLine(c.Num)
		return true
	}

	switch c.Kind {
	case CmdEditReplace:
		ed.EditReplace(c.Letter, c.Num)
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
	}

	switch c.Kind {
	case CmdMiscShowInfo:
		ed.MiscShowInfo()
		return true
	case CmdMiscRepeat:
		ed.MiscRepeat(c.Num)
		return true
	case CmdMiscUndo:
		ed.MiscUndo(c.Num)
		return true
	case CmdMiscRestore:
		ed.MiscRestore()
		return true
	case CmdMiscSaveAndQuit:
		ed.MiscSaveAndQuit()
		return true
	}

	return false
}
