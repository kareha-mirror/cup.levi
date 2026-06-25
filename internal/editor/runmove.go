package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
)

func (ed *Editor) RunMove(c Cmd) (buf.Loc, bool) {
	switch c.Kind {

	case CmdMoveLeft:
		return ed.MoveLeft(c.Num)
	case CmdMoveDown:
		return ed.MoveDown(c.Num)
	case CmdMoveUp:
		return ed.MoveUp(c.Num)
	case CmdMoveRight:
		return ed.MoveRight(c.Num)

	case CmdMoveToStart:
		return ed.MoveToStart()
	case CmdMoveToEnd:
		return ed.MoveToEnd()
	case CmdMoveToNonBlank:
		return ed.MoveToNonBlank()
	case CmdMoveToColumn:
		return ed.MoveToColumn(c.Num)

	case CmdMoveByWord:
		return ed.MoveByWord(c.Num)
	case CmdMoveBackwardByWord:
		return ed.MoveBackwardByWord(c.Num)
	case CmdMoveToEndOfWord:
		return ed.MoveToEndOfWord(c.Num)
	case CmdMoveByLooseWord:
		return ed.MoveByLooseWord(c.Num)
	case CmdMoveBackwardByLooseWord:
		return ed.MoveBackwardByLooseWord(c.Num)
	case CmdMoveToEndOfLooseWord:
		return ed.MoveToEndOfLooseWord(c.Num)

	case CmdMoveByLine:
		return ed.MoveByLine(c.Num)
	case CmdMoveBackwardByLine:
		return ed.MoveBackwardByLine(c.Num)
	case CmdMoveToLastLine:
		return ed.MoveToLastLine()
	case CmdMoveToLine:
		return ed.MoveToLine(c.Num)

	case CmdMoveBySentence:
		return ed.MoveBySentence(c.Num)
	case CmdMoveBackwardBySentence:
		return ed.MoveBackwardBySentence(c.Num)
	case CmdMoveByParagraph:
		return ed.MoveByParagraph(c.Num)
	case CmdMoveBackwardByParagraph:
		return ed.MoveBackwardByParagraph(c.Num)
	case CmdMoveBySection:
		return ed.MoveBySection(c.Num)
	case CmdMoveBackwardBySection:
		return ed.MoveBackwardBySection(c.Num)

	case CmdMoveToTopOfView:
		return ed.MoveToTopOfView()
	case CmdMoveToMiddleOfView:
		return ed.MoveToMiddleOfView()
	case CmdMoveToBottomOfView:
		return ed.MoveToBottomOfView()
	case CmdMoveToBelowTopOfView:
		return ed.MoveToBelowTopOfView(c.Num)
	case CmdMoveToAboveBottomOfView:
		return ed.MoveToAboveBottomOfView(c.Num)

	case CmdMoveToMark:
		return ed.MoveToMark(c.Letter)
	case CmdMoveToMarkLine:
		return ed.MoveToMarkLine(c.Letter)

	case CmdMoveBackToMark:
		return ed.MoveBackToMark()
	case CmdMoveBackToMarkLine:
		return ed.MoveBackToMarkLine()

	case CmdMoveSearchForward:
		return ed.MoveSearchForward()
	case CmdMoveSearchBackward:
		return ed.MoveSearchBackward()
	case CmdMoveSearchNextMatch:
		return ed.MoveSearchNextMatch()
	case CmdMoveSearchPrevMatch:
		return ed.MoveSearchPrevMatch()
	case CmdMoveSearchRepeatForward:
		return ed.MoveSearchRepeatForward()
	case CmdMoveSearchRepeatBackward:
		return ed.MoveSearchRepeatBackward()

	case CmdMoveFindForward:
		return ed.MoveFindForward(c.Letter, c.Num)
	case CmdMoveFindBackward:
		return ed.MoveFindBackward(c.Letter, c.Num)
	case CmdMoveFindBeforeForward:
		return ed.MoveFindBeforeForward(c.Letter, c.Num)
	case CmdMoveFindBeforeBackward:
		return ed.MoveFindBeforeBackward(c.Letter, c.Num)
	case CmdMoveFindNextMatch:
		return ed.MoveFindNextMatch(c.Num)
	case CmdMoveFindPrevMatch:
		return ed.MoveFindPrevMatch(c.Num)

	}

	return buf.Loc{}, false
}
