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
	case cmd.MoveToFirstNonBlank:
		return ed.MoveToFirstNonBlank()
	case cmd.MoveToColumn:
		return ed.MoveToColumn(num)

	case cmd.MoveByWord:
		return ed.MoveByWord(num)
	case cmd.MoveByChangeWord:
		return ed.MoveByChangeWord(num)
	case cmd.MoveByDeleteWord:
		return ed.MoveByDeleteWord(num)
	case cmd.MoveBackByWord:
		return ed.MoveBackByWord(num)
	case cmd.MoveToEndOfWord:
		return ed.MoveToEndOfWord(num)
	case cmd.MoveByBigword:
		return ed.MoveByBigword(num)
	case cmd.MoveByChangeBigword:
		return ed.MoveByChangeBigword(num)
	case cmd.MoveByDeleteBigword:
		return ed.MoveByDeleteBigword(num)
	case cmd.MoveBackByBigword:
		return ed.MoveBackByBigword(num)
	case cmd.MoveToEndOfBigword:
		return ed.MoveToEndOfBigword(num)

	case cmd.MoveByLine:
		return ed.MoveByLine(num)
	case cmd.MoveBackByLine:
		return ed.MoveBackByLine(num)
	case cmd.MoveToLastLine:
		return ed.MoveToLastLine()
	case cmd.MoveToLine:
		return ed.MoveToLine(num)

	case cmd.MoveBySentence:
		return ed.MoveBySentence(num)
	case cmd.MoveBackBySentence:
		return ed.MoveBackBySentence(num)
	case cmd.MoveByParagraph:
		return ed.MoveByParagraph(num)
	case cmd.MoveBackByParagraph:
		return ed.MoveBackByParagraph(num)
	case cmd.MoveBySection:
		return ed.MoveBySection(num)
	case cmd.MoveBackBySection:
		return ed.MoveBackBySection(num)

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
	case cmd.SearchBack:
		return ed.SearchBack()
	case cmd.SearchNext:
		return ed.SearchNext()
	case cmd.SearchPrev:
		return ed.SearchPrev()
	case cmd.RepeatSearch:
		return ed.RepeatSearch()
	case cmd.RepeatBackSearch:
		return ed.RepeatBackSearch()

	case cmd.Find:
		return ed.Find(c.Rune, num)
	case cmd.FindBack:
		return ed.FindBack(c.Rune, num)
	case cmd.FindBefore:
		return ed.FindBefore(c.Rune, num)
	case cmd.FindBeforeBack:
		return ed.FindBeforeBack(c.Rune, num)
	case cmd.FindNext:
		return ed.FindNext(num)
	case cmd.FindPrev:
		return ed.FindPrev(num)

	}
	return buf.Loc{}, false
}
