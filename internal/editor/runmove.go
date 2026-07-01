package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
)

func (ed *Editor) RunMove(c Cmd, num int, alt bool) (buf.Loc, bool) {
	ed.Commit()
	num *= c.Num
	switch c.Kind {

	case MoveLeft:
		return ed.MoveLeft(num)
	case MoveDown:
		return ed.MoveDown(num)
	case MoveHere:
		return ed.MoveHere(num)
	case MoveUp:
		return ed.MoveUp(num)
	case MoveRight:
		return ed.MoveRight(num)

	case MoveToStart:
		return ed.MoveToStart()
	case MoveToEnd:
		return ed.MoveToEnd(num)
	case MoveToAfterIndent:
		return ed.MoveToAfterIndent()
	case MoveToColumn:
		return ed.MoveToColumn(num)

	case MoveByWord:
		if alt {
			return ed.MoveByWordAlt(num)
		} else {
			return ed.MoveByWord(num)
		}
	case MoveByWordAlt: // XXX debug
		return ed.MoveByWordAlt(num)
	case MoveBackwardByWord:
		return ed.MoveBackwardByWord(num)
	case MoveToEndOfWord:
		return ed.MoveToEndOfWord(num)
	case MoveByLooseWord:
		return ed.MoveByLooseWord(num)
	case MoveBackwardByLooseWord:
		return ed.MoveBackwardByLooseWord(num)
	case MoveToEndOfLooseWord:
		return ed.MoveToEndOfLooseWord(num)

	case MoveByLine:
		return ed.MoveByLine(num)
	case MoveBackwardByLine:
		return ed.MoveBackwardByLine(num)
	case MoveToLastLine:
		return ed.MoveToLastLine()
	case MoveToLine:
		return ed.MoveToLine(num)

	case MoveBySentence:
		return ed.MoveBySentence(num)
	case MoveBackwardBySentence:
		return ed.MoveBackwardBySentence(num)
	case MoveByParagraph:
		return ed.MoveByParagraph(num)
	case MoveBackwardByParagraph:
		return ed.MoveBackwardByParagraph(num)
	case MoveBySection:
		return ed.MoveBySection(num)
	case MoveBackwardBySection:
		return ed.MoveBackwardBySection(num)

	case MoveToTopOfView:
		return ed.MoveToTopOfView()
	case MoveToMiddleOfView:
		return ed.MoveToMiddleOfView()
	case MoveToBottomOfView:
		return ed.MoveToBottomOfView()
	case MoveToBelowTopOfView:
		return ed.MoveToBelowTopOfView(num)
	case MoveToAboveBottomOfView:
		return ed.MoveToAboveBottomOfView(num)

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
		return ed.Find(c.Rune, num)
	case FindBackward:
		return ed.FindBackward(c.Rune, num)
	case FindBefore:
		return ed.FindBefore(c.Rune, num)
	case FindBeforeBackward:
		return ed.FindBeforeBackward(c.Rune, num)
	case FindNext:
		return ed.FindNext(num)
	case FindPrev:
		return ed.FindPrev(num)

	}
	return buf.Loc{}, false
}
