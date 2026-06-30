package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
)

func (ed *Editor) RunMove(c Cmd) (buf.Loc, bool) {
	ed.Commit()
	switch c.Kind {

	case MoveLeft:
		return ed.MoveLeft(c.Num)
	case MoveDown:
		return ed.MoveDown(c.Num)
	case MoveUp:
		return ed.MoveUp(c.Num)
	case MoveRight:
		return ed.MoveRight(c.Num)

	case MoveToStart:
		return ed.MoveToStart()
	case MoveToEnd:
		return ed.MoveToEnd()
	case MoveToAfterIndent:
		return ed.MoveToAfterIndent()
	case MoveToColumn:
		return ed.MoveToColumn(c.Num)

	case MoveByWord:
		return ed.MoveByWord(c.Num)
	case MoveByChangeWord:
		return ed.MoveByChangeWord(c.Num)
	case MoveBackwardByWord:
		return ed.MoveBackwardByWord(c.Num)
	case MoveToEndOfWord:
		return ed.MoveToEndOfWord(c.Num)
	case MoveByLooseWord:
		return ed.MoveByLooseWord(c.Num)
	case MoveBackwardByLooseWord:
		return ed.MoveBackwardByLooseWord(c.Num)
	case MoveToEndOfLooseWord:
		return ed.MoveToEndOfLooseWord(c.Num)

	case MoveByLine:
		return ed.MoveByLine(c.Num)
	case MoveBackwardByLine:
		return ed.MoveBackwardByLine(c.Num)
	case MoveToLastLine:
		return ed.MoveToLastLine()
	case MoveToLine:
		return ed.MoveToLine(c.Num)

	case MoveBySentence:
		return ed.MoveBySentence(c.Num)
	case MoveBackwardBySentence:
		return ed.MoveBackwardBySentence(c.Num)
	case MoveByParagraph:
		return ed.MoveByParagraph(c.Num)
	case MoveBackwardByParagraph:
		return ed.MoveBackwardByParagraph(c.Num)
	case MoveBySection:
		return ed.MoveBySection(c.Num)
	case MoveBackwardBySection:
		return ed.MoveBackwardBySection(c.Num)

	case MoveToTopOfView:
		return ed.MoveToTopOfView()
	case MoveToMiddleOfView:
		return ed.MoveToMiddleOfView()
	case MoveToBottomOfView:
		return ed.MoveToBottomOfView()
	case MoveToBelowTopOfView:
		return ed.MoveToBelowTopOfView(c.Num)
	case MoveToAboveBottomOfView:
		return ed.MoveToAboveBottomOfView(c.Num)

	case MoveToMark:
		return ed.MoveToMark(c.Rune)
	case MoveToMarkLine:
		return ed.MoveToMarkLine(c.Rune)

	case BackToMark:
		return ed.BackToMark()
	case BackToMarkLine:
		return ed.BackToMarkLine()

	case Search:
		return ed.Search()
	case SearchBackward:
		return ed.SearchBackward()
	case SearchNext:
		return ed.SearchNext()
	case SearchPrev:
		return ed.SearchPrev()
	case RepeatSearch:
		return ed.RepeatSearch()
	case RepeatBackwardSearch:
		return ed.RepeatBackwardSearch()

	case Find:
		return ed.Find(c.Rune, c.Num)
	case FindBackward:
		return ed.FindBackward(c.Rune, c.Num)
	case FindBefore:
		return ed.FindBefore(c.Rune, c.Num)
	case FindBeforeBackward:
		return ed.FindBeforeBackward(c.Rune, c.Num)
	case FindNext:
		return ed.FindNext(c.Num)
	case FindPrev:
		return ed.FindPrev(c.Num)

	}
	return buf.Loc{}, false
}
