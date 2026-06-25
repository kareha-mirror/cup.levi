package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
)

func isBlankRune(r rune) bool {
	return r == ' ' || r == '\t'
}

func nonBlankCol(s string) int {
	i := 0
	for _, r := range s {
		if !isBlankRune(r) {
			break
		}
		i++
	}
	return i
}

/////////////////////
// Motion Commands //
/////////////////////

// Note: Marking, Search, Character Finding Commands also have Motion Commands.

//
// Move by Character / Move by Line
//

// h : Move cursor left by character.
func (ed *Editor) MoveLeft(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveLeft: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	b := ed.Buf()
	loc := b.Loc
	loc.Col -= n
	loc.Col = b.ConfineCol(loc)
	return buf.Dest{
		Loc:       loc,
		Linewise:  false,
		FreeCol:   false,
		Inclusive: false,
	}, true
}

// j : Move cursor down by line.
func (ed *Editor) MoveDown(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveDown: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	b := ed.Buf()
	loc := b.Loc
	if !b.CheckRow(loc.Row + n) {
		return buf.Dest{}, false
	}
	loc.Row += n
	return buf.Dest{
		Loc:       loc,
		Linewise:  true,
		FreeCol:   true,
		Inclusive: true,
	}, true
}

// k : Move cursor up by line.
func (ed *Editor) MoveUp(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveUp: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	b := ed.Buf()
	loc := b.Loc
	if !b.CheckRow(loc.Row - n) {
		return buf.Dest{}, false
	}
	loc.Row -= n
	return buf.Dest{
		Loc:       loc,
		Linewise:  true,
		FreeCol:   true,
		Inclusive: true,
	}, true
}

// l : Move cursor right by character.
func (ed *Editor) MoveRight(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveRight: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	b := ed.Buf()
	loc := b.Loc
	loc.Col += n
	loc.Col = b.ConfineCol(loc)
	return buf.Dest{
		Loc:       loc,
		Linewise:  false,
		FreeCol:   false,
		Inclusive: false,
	}, true
}

//
// Move in Line
//

// 0 : Move cursor to start of current line.
func (ed *Editor) MoveToStart() (buf.Dest, bool) {
	ed.EnsureCommand()
	b := ed.Buf()
	loc := b.Loc
	loc.Col = 0
	return buf.Dest{
		Loc:       loc,
		Linewise:  false,
		FreeCol:   false,
		Inclusive: true,
	}, true
}

// $ : Move cursor to end of current line.
func (ed *Editor) MoveToEnd() (buf.Dest, bool) {
	ed.EnsureCommand()
	b := ed.Buf()
	loc := b.Loc
	line := b.Line(loc.Row)
	rc := utf8.RuneCountInString(line)
	loc.Col = max(rc-1, 0)
	return buf.Dest{
		Loc:       loc,
		Linewise:  false,
		FreeCol:   false,
		Inclusive: true,
	}, true
}

// ^ : Move cursor to first non-blank character of current line.
func (ed *Editor) MoveToNonBlank() (buf.Dest, bool) {
	ed.EnsureCommand()
	b := ed.Buf()
	loc := b.Loc
	line := b.Line(loc.Row)
	loc.Col = nonBlankCol(line)
	return buf.Dest{
		Loc:       loc,
		Linewise:  false,
		FreeCol:   false,
		Inclusive: true,
	}, true
}

// <num>| : Move cursor to column <num> of current line.
// (Note: Proper vi's column number is visual-based, but levi' is rune-based.)
func (ed *Editor) MoveToColumn(n int) (buf.Dest, bool) { // n: 1-based
	if n < 1 {
		ed.Error("MoveToColumn: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	b := ed.Buf()
	loc := b.Loc
	loc.Col = n - 1
	loc.Col = b.ConfineCol(loc)
	return buf.Dest{
		Loc:       loc,
		Linewise:  false,
		FreeCol:   false,
		Inclusive: true,
	}, true
}

//
// Move by Word / Move by Loose Word
//

// w : Move cursor forward by word.
func (ed *Editor) MoveByWord(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveByWord: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	b := ed.Buf()
	for i := 0; i < n; i++ {
		if b.MoveByWord() {
			continue
		}
		if b.Loc.Row >= b.NumLines()-1 {
			ed.MoveToEnd()
			return buf.Dest{}, false // TODO
		}
		b.Loc.Row++
		b.Loc.Col = 0
		if !b.SkipBlankLines() {
			return buf.Dest{}, false // TODO
		}
	}
	return buf.Dest{}, false // TODO
}

// internal use : Move cursor forward by word used by cw and dw.
func (ed *Editor) MoveByWordEx(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveByWordEx: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	b := ed.Buf()
	for i := 0; i < n; i++ {
		if b.MoveByWordEx() {
			continue
		}
		if b.Loc.Row >= b.NumLines()-1 {
			ed.MoveToEnd()
			return buf.Dest{}, false // TODO
		}
		b.Loc.Row++
		b.Loc.Col = 0
		if !b.SkipBlankLines() {
			return buf.Dest{}, false // TODO
		}
	}
	return buf.Dest{}, false // TODO
}

// b : Move cursor backward by word.
func (ed *Editor) MoveBackwardByWord(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveBackwardByWord: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	b := ed.Buf()
	for i := 0; i < n; i++ {
		if b.Loc.Col > 0 {
			b.Loc.Col--
		} else {
			if b.Loc.Row < 1 {
				return buf.Dest{}, false // TODO
			}
			b.Loc.Row--
			b.Loc.Col = max(utf8.RuneCountInString(b.CurrentLine())-1, 0)
		}
		if !b.SkipBackwardBlankLines() {
			return buf.Dest{}, false // TODO
		}
		if !b.MoveBackwardByWord() {
			return buf.Dest{}, false // TODO
		}
	}
	return buf.Dest{}, false // TODO
}

// e : Move cursor to end of word.
func (ed *Editor) MoveToEndOfWord(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveToEndOfWord: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveToEndOfWord")
	return buf.Dest{}, false
}

// W : Move cursor forward by loose word.
func (ed *Editor) MoveByLooseWord(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveByLooseWord: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveByLooseWord")
	return buf.Dest{}, false
}

// B : Move cursor backward by loose word.
func (ed *Editor) MoveBackwardByLooseWord(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveBackwardByLooseWord: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardByLooseWord")
	return buf.Dest{}, false
}

// E : Move cursor to end of loose word.
func (ed *Editor) MoveToEndOfLooseWord(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveToEndOfLooseWord: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveToEndOfLooseWord")
	return buf.Dest{}, false
}

//
// Move by Line
//

// Enter, + : Move cursor to first non-blank character of next line.
func (ed *Editor) MoveByLine(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveByLine: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	b := ed.Buf()
	loc := b.Loc
	if !b.CheckRow(loc.Row + n) {
		return buf.Dest{}, false
	}
	loc.Row += n
	line := b.Line(loc.Row)
	loc.Col = nonBlankCol(line)
	return buf.Dest{
		Loc:       loc,
		Linewise:  true,
		FreeCol:   false,
		Inclusive: true,
	}, true
}

// - : Move cursor to first non-blank character of previous line.
func (ed *Editor) MoveBackwardByLine(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveBackwardByLine: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	b := ed.Buf()
	loc := b.Loc
	if !b.CheckRow(loc.Row - n) {
		return buf.Dest{}, false
	}
	loc.Row -= n
	line := b.Line(loc.Row)
	loc.Col = nonBlankCol(line)
	return buf.Dest{
		Loc:       loc,
		Linewise:  true,
		FreeCol:   false,
		Inclusive: true,
	}, true
}

// G : Move cursor to first non-blank character of last line.
func (ed *Editor) MoveToLastLine() (buf.Dest, bool) {
	ed.EnsureCommand()
	b := ed.Buf()
	loc := b.Loc
	loc.Row = b.ConfineRow(b.NumLines() - 1)
	line := b.Line(loc.Row)
	loc.Col = nonBlankCol(line)
	return buf.Dest{
		Loc:       loc,
		Linewise:  true,
		FreeCol:   false,
		Inclusive: true,
	}, true
}

// <num>G : Move cursor to first non-blank character of line specified by <num>.
func (ed *Editor) MoveToLine(n int) (buf.Dest, bool) { // n: 1-based
	if n < 1 {
		ed.Error("MoveToLine: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	b := ed.Buf()
	loc := b.Loc
	if !b.CheckRow(n - 1) {
		return buf.Dest{}, false
	}
	loc.Row = n - 1
	line := b.Line(loc.Row)
	loc.Col = nonBlankCol(line)
	return buf.Dest{
		Loc:       loc,
		Linewise:  true,
		FreeCol:   false,
		Inclusive: true,
	}, true
}

//
// Move by Block
//

// ) : Move cursor forward by sentence.
func (ed *Editor) MoveBySentence(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveBySentence: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBySentence")
	return buf.Dest{}, false
}

// ( : Move cursor backward by sentence.
func (ed *Editor) MoveBackwardBySentence(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveBackwardBySentence: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardBySentence")
	return buf.Dest{}, false
}

// } : Move cursor forward by paragraph.
func (ed *Editor) MoveByParagraph(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveByParagraph: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveByParagraph")
	return buf.Dest{}, false
}

// { : Move cursor backward by paragraph.
func (ed *Editor) MoveBackwardByParagraph(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveBackwardByParagraph: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardByParagraph")
	return buf.Dest{}, false
}

// ]] : Move cursor forward by section.
func (ed *Editor) MoveBySection(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveBySection: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBySection")
	return buf.Dest{}, false
}

// [[ : Move cursor backward by section.
func (ed *Editor) MoveBackwardBySection(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveBackwardBySection: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardBySection")
	return buf.Dest{}, false
}

//
// Move in View
//

// H : Move cursor to top of view.
func (ed *Editor) MoveToTopOfView() (buf.Dest, bool) {
	ed.EnsureCommand()
	if len(ed.vMeta) < 1 {
		return buf.Dest{}, false
	}
	i := 0
	loc := ed.vMeta[i].Loc
	if loc.Col < 1 {
		b := ed.Buf()
		line := b.Line(loc.Row)
		loc.Col = nonBlankCol(line)
	}
	return buf.Dest{
		Loc:       loc,
		Linewise:  true,
		FreeCol:   false,
		Inclusive: true,
	}, true
}

// M : Move cursor to middle of view.
func (ed *Editor) MoveToMiddleOfView() (buf.Dest, bool) {
	ed.EnsureCommand()
	if len(ed.vMeta) < 1 {
		return buf.Dest{}, false
	}
	i := len(ed.vMeta)/2 - 1
	loc := ed.vMeta[i].Loc
	if loc.Col < 1 {
		b := ed.Buf()
		line := b.Line(loc.Row)
		loc.Col = nonBlankCol(line)
	}
	return buf.Dest{
		Loc:       loc,
		Linewise:  true,
		FreeCol:   false,
		Inclusive: true,
	}, true
}

// L : Move cursor to bottom of view.
func (ed *Editor) MoveToBottomOfView() (buf.Dest, bool) {
	ed.EnsureCommand()
	if len(ed.vMeta) < 1 {
		return buf.Dest{}, false
	}
	i := len(ed.vMeta) - 1
	loc := ed.vMeta[i].Loc
	if loc.Col < 1 {
		b := ed.Buf()
		line := b.Line(loc.Row)
		loc.Col = nonBlankCol(line)
	}
	return buf.Dest{
		Loc:       loc,
		Linewise:  true,
		FreeCol:   false,
		Inclusive: true,
	}, true
}

// <num>H : Move cursor below <num> lines from top of view.
func (ed *Editor) MoveToBelowTopOfView(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveToBelowTopOfView: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	if len(ed.vMeta) < 1 {
		return buf.Dest{}, false
	}
	i := n - 1
	loc := ed.vMeta[i].Loc
	if loc.Col < 1 {
		b := ed.Buf()
		line := b.Line(loc.Row)
		loc.Col = nonBlankCol(line)
	}
	return buf.Dest{
		Loc:       loc,
		Linewise:  true,
		FreeCol:   false,
		Inclusive: true,
	}, true
}

// <num>L : Move cursor above <num> lines from bottom of view.
func (ed *Editor) MoveToAboveBottomOfView(n int) (buf.Dest, bool) {
	if n < 1 {
		ed.Error("MoveToAboveBottomOfView: n < 1")
		return buf.Dest{}, false
	}
	ed.EnsureCommand()
	i := len(ed.vMeta) - n
	if i < 0 {
		return buf.Dest{}, false
	}
	loc := ed.vMeta[i].Loc
	if loc.Col < 1 {
		b := ed.Buf()
		line := b.Line(loc.Row)
		loc.Col = nonBlankCol(line)
	}
	return buf.Dest{
		Loc:       loc,
		Linewise:  true,
		FreeCol:   false,
		Inclusive: true,
	}, true
}
