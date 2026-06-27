package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
)

/////////////////////
// Motion Commands //
/////////////////////

// Note: Marking, Search, Character Finding Commands also have Motion Commands.

//
// Move by Character / Move by Line
//

// h : Move cursor left by character.
func (ed *Editor) MoveLeft(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveLeft: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Col -= n
	loc.Col = b.ConfineCol(loc)
	return loc, true
}

// j : Move cursor down by line.
func (ed *Editor) MoveDown(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveDown: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Row += n
	if !b.CheckRowInclusive(loc.Row) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	return loc, true
}

// k : Move cursor up by line.
func (ed *Editor) MoveUp(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveUp: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Row -= n
	if !b.CheckRowInclusive(loc.Row) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	return loc, true
}

// l : Move cursor right by character.
func (ed *Editor) MoveRight(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveRight: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Col += n
	loc.Col = b.ConfineCol(loc)
	return loc, true
}

//
// Move in Line
//

// 0 : Move cursor to start of current line.
func (ed *Editor) MoveToStart() (buf.Loc, bool) {
	loc := ed.Buf().Loc
	loc.Col = 0
	return loc, true
}

// $ : Move cursor to end of current line.
func (ed *Editor) MoveToEnd() (buf.Loc, bool) {
	b := ed.Buf()
	loc := b.Loc
	loc.Col = utf8.RuneCountInString(b.Line(loc.Row))
	return loc, true
}

// ^ : Move cursor to first non-blank character of current line.
func (ed *Editor) MoveToNonBlank() (buf.Loc, bool) {
	b := ed.Buf()
	loc := b.Loc
	loc.Col = b.NonBlankColOfLine(loc.Row)
	return loc, true
}

// <num>| : Move cursor to column <num> of current line.
// (Note: Proper vi's column number is visual-based, but levi' is rune-based.)
func (ed *Editor) MoveToColumn(n int) (buf.Loc, bool) { // n: 1-based
	if n < 1 {
		ed.Error("MoveToColumn: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Col = n - 1
	loc.Col = b.ConfineCol(loc)
	return loc, true
}

//
// Move by Word / Move by Loose Word
//

// w : Move cursor forward by word.
func (ed *Editor) MoveByWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByWord: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 0; i < n; i++ {
		if loc, found = b.MoveByWord(loc); found {
			continue
		}
		loc.Row++
		loc.Col = 0
		if loc, found = b.SkipBlankLines(loc); !found {
			return loc, true
		}
	}
	return loc, true
}

// internal use : Move cursor forward by word used by cw and dw.
func (ed *Editor) MoveByWordEx(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByWordEx: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 0; i < n; i++ {
		if loc, found = b.MoveByWordEx(loc); found {
			continue
		}
		return loc, true
	}
	return loc, true
}

// b : Move cursor backward by word.
func (ed *Editor) MoveBackwardByWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBackwardByWord: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 0; i < n; i++ {
		if loc.Col > 0 {
			loc.Col--
		} else {
			if loc.Row < 1 {
				return loc, true
			}
			loc.Row--
			line := b.Line(loc.Row)
			rc := utf8.RuneCountInString(line)
			loc.Col = max(rc-1, 0)
		}
		if loc, found = b.SkipBackwardBlankLines(loc); !found {
			return loc, true
		}
		if loc, found = b.MoveBackwardByWord(loc); !found {
			return loc, true
		}
	}
	return loc, true
}

// e : Move cursor to end of word.
func (ed *Editor) MoveToEndOfWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveToEndOfWord: n < 1")
		return buf.Loc{}, false
	}
	ed.Unimplemented("MoveToEndOfWord")
	return buf.Loc{}, false
}

// W : Move cursor forward by loose word.
func (ed *Editor) MoveByLooseWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByLooseWord: n < 1")
		return buf.Loc{}, false
	}
	ed.Unimplemented("MoveByLooseWord")
	return buf.Loc{}, false
}

// B : Move cursor backward by loose word.
func (ed *Editor) MoveBackwardByLooseWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBackwardByLooseWord: n < 1")
		return buf.Loc{}, false
	}
	ed.Unimplemented("MoveBackwardByLooseWord")
	return buf.Loc{}, false
}

// E : Move cursor to end of loose word.
func (ed *Editor) MoveToEndOfLooseWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveToEndOfLooseWord: n < 1")
		return buf.Loc{}, false
	}
	ed.Unimplemented("MoveToEndOfLooseWord")
	return buf.Loc{}, false
}

//
// Move by Line
//

// Enter, + : Move cursor to first non-blank character of next line.
func (ed *Editor) MoveByLine(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByLine: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Row += n
	if !b.CheckRowInclusive(loc.Row) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	loc.Col = b.NonBlankColOfLine(loc.Row)
	return loc, true
}

// - : Move cursor to first non-blank character of previous line.
func (ed *Editor) MoveBackwardByLine(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBackwardByLine: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Row -= n
	if !b.CheckRowInclusive(loc.Row) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	loc.Col = b.NonBlankColOfLine(loc.Row)
	return loc, true
}

// G : Move cursor to first non-blank character of last line.
func (ed *Editor) MoveToLastLine() (buf.Loc, bool) {
	var loc buf.Loc
	b := ed.Buf()
	loc.Row = b.ConfineRow(b.NumLines() - 1)
	loc.Col = b.NonBlankColOfLine(loc.Row)
	return loc, true
}

// <num>G : Move cursor to first non-blank character of line specified by <num>.
func (ed *Editor) MoveToLine(n int) (buf.Loc, bool) { // n: 1-based
	if n < 1 {
		ed.Error("MoveToLine: n < 1")
		return buf.Loc{}, false
	}
	var loc buf.Loc
	loc.Row = n - 1
	b := ed.Buf()
	if !b.CheckRowInclusive(loc.Row) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	loc.Col = b.NonBlankColOfLine(loc.Row)
	return loc, true
}

//
// Move by Block
//

// ) : Move cursor forward by sentence.
func (ed *Editor) MoveBySentence(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBySentence: n < 1")
		return buf.Loc{}, false
	}
	ed.Unimplemented("MoveBySentence")
	return buf.Loc{}, false
}

// ( : Move cursor backward by sentence.
func (ed *Editor) MoveBackwardBySentence(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBackwardBySentence: n < 1")
		return buf.Loc{}, false
	}
	ed.Unimplemented("MoveBackwardBySentence")
	return buf.Loc{}, false
}

// } : Move cursor forward by paragraph.
func (ed *Editor) MoveByParagraph(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByParagraph: n < 1")
		return buf.Loc{}, false
	}
	ed.Unimplemented("MoveByParagraph")
	return buf.Loc{}, false
}

// { : Move cursor backward by paragraph.
func (ed *Editor) MoveBackwardByParagraph(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBackwardByParagraph: n < 1")
		return buf.Loc{}, false
	}
	ed.Unimplemented("MoveBackwardByParagraph")
	return buf.Loc{}, false
}

// ]] : Move cursor forward by section.
func (ed *Editor) MoveBySection(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBySection: n < 1")
		return buf.Loc{}, false
	}
	ed.Unimplemented("MoveBySection")
	return buf.Loc{}, false
}

// [[ : Move cursor backward by section.
func (ed *Editor) MoveBackwardBySection(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBackwardBySection: n < 1")
		return buf.Loc{}, false
	}
	ed.Unimplemented("MoveBackwardBySection")
	return buf.Loc{}, false
}

//
// Move in View
//

// H : Move cursor to top of view.
func (ed *Editor) MoveToTopOfView() (buf.Loc, bool) {
	if len(ed.viewMeta) < 1 {
		return buf.Loc{}, false
	}
	loc := ed.viewMeta[0].Loc
	if loc.Col < 1 {
		loc.Col = ed.Buf().NonBlankColOfLine(loc.Row)
	}
	return loc, true
}

// M : Move cursor to middle of view.
func (ed *Editor) MoveToMiddleOfView() (buf.Loc, bool) {
	if len(ed.viewMeta) < 1 {
		return buf.Loc{}, false
	}
	i := len(ed.viewMeta)/2 - 1
	loc := ed.viewMeta[i].Loc
	if loc.Col < 1 {
		loc.Col = ed.Buf().NonBlankColOfLine(loc.Row)
	}
	return loc, true
}

// L : Move cursor to bottom of view.
func (ed *Editor) MoveToBottomOfView() (buf.Loc, bool) {
	if len(ed.viewMeta) < 1 {
		return buf.Loc{}, false
	}
	i := len(ed.viewMeta) - 1
	loc := ed.viewMeta[i].Loc
	if loc.Col < 1 {
		loc.Col = ed.Buf().NonBlankColOfLine(loc.Row)
	}
	return loc, true
}

// <num>H : Move cursor below <num> lines from top of view.
func (ed *Editor) MoveToBelowTopOfView(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveToBelowTopOfView: n < 1")
		return buf.Loc{}, false
	}
	i := n - 1
	if i >= len(ed.viewMeta) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	loc := ed.viewMeta[i].Loc
	if loc.Col < 1 {
		loc.Col = ed.Buf().NonBlankColOfLine(loc.Row)
	}
	return loc, true
}

// <num>L : Move cursor above <num> lines from bottom of view.
func (ed *Editor) MoveToAboveBottomOfView(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveToAboveBottomOfView: n < 1")
		return buf.Loc{}, false
	}
	i := len(ed.viewMeta) - n
	if i < 0 {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	loc := ed.viewMeta[i].Loc
	if loc.Col < 1 {
		loc.Col = ed.Buf().NonBlankColOfLine(loc.Row)
	}
	return loc, true
}
