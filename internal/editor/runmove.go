package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/cmd"
)

func (ed *Editor) RunMove(c cmd.Cmd, num int) (buf.Loc, bool) {
	ed.Commit()
	num *= c.Num
	switch c.Kind {

	case cmd.MoveLeft:
		return ed.MoveLeft(num)
	case cmd.MoveDown:
		return ed.MoveDown(num)
	case cmd.MoveHere:
		return ed.MoveHere(num)
	case cmd.MoveUp:
		return ed.MoveUp(num)
	case cmd.MoveRight:
		return ed.MoveRight(num)

	case cmd.MoveToStart:
		return ed.MoveToStart()
	case cmd.MoveToEnd:
		return ed.MoveToEnd(num)
	case cmd.MoveToAfterIndent:
		return ed.MoveToAfterIndent()
	case cmd.MoveToColumn:
		return ed.MoveToColumn(num)

	case cmd.MoveByWord:
		return ed.MoveByWord(num)
	case cmd.MoveByChangeWord:
		return ed.MoveByChangeWord(num)
	case cmd.MoveByDeleteWord:
		return ed.MoveByDeleteWord(num)
	case cmd.MoveBackwardByWord:
		return ed.MoveBackwardByWord(num)
	case cmd.MoveToEndOfWord:
		return ed.MoveToEndOfWord(num)
	case cmd.MoveByLooseWord:
		return ed.MoveByLooseWord(num)
	case cmd.MoveBackwardByLooseWord:
		return ed.MoveBackwardByLooseWord(num)
	case cmd.MoveToEndOfLooseWord:
		return ed.MoveToEndOfLooseWord(num)

	case cmd.MoveByLine:
		return ed.MoveByLine(num)
	case cmd.MoveBackwardByLine:
		return ed.MoveBackwardByLine(num)
	case cmd.MoveToLastLine:
		return ed.MoveToLastLine()
	case cmd.MoveToLine:
		return ed.MoveToLine(num)

	case cmd.MoveBySentence:
		return ed.MoveBySentence(num)
	case cmd.MoveBackwardBySentence:
		return ed.MoveBackwardBySentence(num)
	case cmd.MoveByParagraph:
		return ed.MoveByParagraph(num)
	case cmd.MoveBackwardByParagraph:
		return ed.MoveBackwardByParagraph(num)
	case cmd.MoveBySection:
		return ed.MoveBySection(num)
	case cmd.MoveBackwardBySection:
		return ed.MoveBackwardBySection(num)

	case cmd.MoveToTopOfView:
		return ed.MoveToTopOfView()
	case cmd.MoveToMiddleOfView:
		return ed.MoveToMiddleOfView()
	case cmd.MoveToBottomOfView:
		return ed.MoveToBottomOfView()
	case cmd.MoveToBelowTopOfView:
		return ed.MoveToBelowTopOfView(num)
	case cmd.MoveToAboveBottomOfView:
		return ed.MoveToAboveBottomOfView(num)

	case cmd.MoveToMark:
		return ed.MoveToMark(c.Rune)
	case cmd.MoveToMarkLine:
		return ed.MoveToMarkLine(c.Rune)

	case cmd.BackToMark:
		return ed.BackToMark()
	case cmd.BackToMarkLine:
		return ed.BackToMarkLine()

	case cmd.Search:
		return ed.Search()
	case cmd.SearchBackward:
		return ed.SearchBackward()
	case cmd.SearchNext:
		return ed.SearchNext()
	case cmd.SearchPrev:
		return ed.SearchPrev()
	case cmd.RepeatSearch:
		return ed.RepeatSearch()
	case cmd.RepeatBackwardSearch:
		return ed.RepeatBackwardSearch()

	case cmd.Find:
		return ed.Find(c.Rune, num)
	case cmd.FindBackward:
		return ed.FindBackward(c.Rune, num)
	case cmd.FindBefore:
		return ed.FindBefore(c.Rune, num)
	case cmd.FindBeforeBackward:
		return ed.FindBeforeBackward(c.Rune, num)
	case cmd.FindNext:
		return ed.FindNext(num)
	case cmd.FindPrev:
		return ed.FindPrev(num)

	}
	return buf.Loc{}, false
}
