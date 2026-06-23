package editor

import (
	"unicode/utf8"
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

func (ed *Editor) toNonBlankCol() {
	line := ed.CurrentLine()
	ed.Buffer().MoveCol(nonBlankCol(line))
}

/////////////////////
// Motion Commands //
/////////////////////

//
// Move by Character / Move by Line
//

// h : Move cursor left by character.
func (ed *Editor) MoveLeft(n int) {
	if n < 1 {
		ed.Error("MoveLeft: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Buffer().AdjustCol(-n)
}

// j : Move cursor down by line.
func (ed *Editor) MoveDown(n int) {
	if n < 1 {
		ed.Error("MoveDown: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buffer()
	if !b.AdjustRow(n) {
		return
	}
	b.LoadVirtCol()
	b.ConfineCol()
}

// k : Move cursor up by line.
func (ed *Editor) MoveUp(n int) {
	if n < 1 {
		ed.Error("MoveUp: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buffer()
	if !b.AdjustRow(-n) {
		return
	}
	b.LoadVirtCol()
	b.ConfineCol()
}

// l : Move cursor right by character.
func (ed *Editor) MoveRight(n int) {
	if n < 1 {
		ed.Error("MoveRight: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Buffer().AdjustCol(n)
}

//
// Move in Line
//

// 0 : Move cursor to start of current line.
func (ed *Editor) MoveToStart() {
	ed.EnsureCommand()
	ed.Buffer().MoveCol(0)
}

// $ : Move cursor to end of current line.
func (ed *Editor) MoveToEnd() {
	ed.EnsureCommand()
	ed.Buffer().MoveCol(ed.RuneCount() - 1)
}

// ^ : Move cursor to first non-blank character of current line.
func (ed *Editor) MoveToNonBlank() {
	ed.EnsureCommand()
	ed.toNonBlankCol()
}

// <num>| : Move cursor to column <num> of current line.
// (Note: Proper vi's column number is visual-based, but levi' is rune-based.)
func (ed *Editor) MoveToColumn(n int) { // n: 1-based
	if n < 1 {
		ed.Error("MoveToColumn: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Buffer().MoveCol(n - 1)
}

//
// Move by Word / Move by Loose Word
//

// w : Move cursor forward by word.
func (ed *Editor) MoveByWord(n int) {
	if n < 1 {
		ed.Error("MoveByWord: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buffer()
	for i := 0; i < n; i++ {
		if b.MoveByWord() {
			continue
		}
		if b.Loc.Row >= b.NumLines()-1 {
			ed.MoveToEnd()
			return
		}
		b.Loc.Row++
		b.Loc.Col = 0
		if !b.SkipBlankLines() {
			return
		}
	}
}

// b : Move cursor backward by word.
func (ed *Editor) MoveBackwardByWord(n int) {
	if n < 1 {
		ed.Error("MoveBackwardByWord: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buffer()
	for i := 0; i < n; i++ {
		if b.Loc.Col > 0 {
			b.Loc.Col--
		} else {
			if b.Loc.Row < 1 {
				return
			}
			b.Loc.Row--
			b.Loc.Col = max(utf8.RuneCountInString(b.CurrentLine())-1, 0)
		}
		if !b.SkipBackwardBlankLines() {
			return
		}
		if !b.MoveBackwardByWord() {
			return
		}
	}
}

// e : Move cursor to end of word.
func (ed *Editor) MoveToEndOfWord(n int) {
	if n < 1 {
		ed.Error("MoveToEndOfWord: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveToEndOfWord")
}

// W : Move cursor forward by loose word.
func (ed *Editor) MoveByLooseWord(n int) {
	if n < 1 {
		ed.Error("MoveByLooseWord: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveByLooseWord")
}

// B : Move cursor backward by loose word.
func (ed *Editor) MoveBackwardByLooseWord(n int) {
	if n < 1 {
		ed.Error("MoveBackwardByLooseWord: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardByLooseWord")
}

// E : Move cursor to end of loose word.
func (ed *Editor) MoveToEndOfLooseWord(n int) {
	if n < 1 {
		ed.Error("MoveToEndOfLooseWord: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveToEndOfLooseWord")
}

//
// Move by Line
//

// Enter, + : Move cursor to first non-blank character of next line.
func (ed *Editor) MoveByLine(n int) {
	if n < 1 {
		ed.Error("MoveByLine: n < 1")
		return
	}
	ed.EnsureCommand()
	if !ed.Buffer().AdjustRow(n) {
		return
	}
	ed.toNonBlankCol()
}

// - : Move cursor to first non-blank character of previous line.
func (ed *Editor) MoveBackwardByLine(n int) {
	if n < 1 {
		ed.Error("MoveBackwardByLine: n < 1")
		return
	}
	ed.EnsureCommand()
	if !ed.Buffer().AdjustRow(-n) {
		return
	}
	ed.toNonBlankCol()
}

// G : Move cursor to first non-blank character of last line.
func (ed *Editor) MoveToLastLine() {
	ed.EnsureCommand()
	b := ed.Buffer()
	b.Loc.Row = b.NumLines() - 1
	b.ConfineRow()
	ed.toNonBlankCol()
}

// <num>G : Move cursor to first non-blank character of line specified by <num>.
func (ed *Editor) MoveToLine(n int) { // n: 1-based
	if n < 1 {
		ed.Error("MoveToLine: n < 1")
		return
	}
	ed.EnsureCommand()
	if !ed.Buffer().MoveRow(n - 1) {
		return
	}
	ed.toNonBlankCol()
}

//
// Move by Block
//

// ) : Move cursor forward by sentence.
func (ed *Editor) MoveBySentence(n int) {
	if n < 1 {
		ed.Error("MoveBySentence: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBySentence")
}

// ( : Move cursor backward by sentence.
func (ed *Editor) MoveBackwardBySentence(n int) {
	if n < 1 {
		ed.Error("MoveBackwardBySentence: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardBySentence")
}

// } : Move cursor forward by paragraph.
func (ed *Editor) MoveByParagraph(n int) {
	if n < 1 {
		ed.Error("MoveByParagraph: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveByParagraph")
}

// { : Move cursor backward by paragraph.
func (ed *Editor) MoveBackwardByParagraph(n int) {
	if n < 1 {
		ed.Error("MoveBackwardByParagraph: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardByParagraph")
}

// ]] : Move cursor forward by section.
func (ed *Editor) MoveBySection(n int) {
	if n < 1 {
		ed.Error("MoveBySection: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBySection")
}

// [[ : Move cursor backward by section.
func (ed *Editor) MoveBackwardBySection(n int) {
	if n < 1 {
		ed.Error("MoveBackwardBySection: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardBySection")
}

//
// Move in View
//

// H : Move cursor to top of view.
func (ed *Editor) MoveToTopOfView() {
	ed.EnsureCommand()
	if len(ed.vMeta) < 1 {
		return
	}
	i := 0
	b := ed.Buffer()
	b.Loc = ed.vMeta[i].Loc
	if b.Loc.Col < 1 {
		ed.toNonBlankCol()
	}
}

// M : Move cursor to middle of view.
func (ed *Editor) MoveToMiddleOfView() {
	ed.EnsureCommand()
	if len(ed.vMeta) < 1 {
		return
	}
	i := len(ed.vMeta)/2 - 1
	b := ed.Buffer()
	b.Loc = ed.vMeta[i].Loc
	if b.Loc.Col < 1 {
		ed.toNonBlankCol()
	}
}

// L : Move cursor to bottom of view.
func (ed *Editor) MoveToBottomOfView() {
	ed.EnsureCommand()
	if len(ed.vMeta) < 1 {
		return
	}
	i := len(ed.vMeta) - 1
	b := ed.Buffer()
	b.Loc = ed.vMeta[i].Loc
	if b.Loc.Col < 1 {
		ed.toNonBlankCol()
	}
}

// <num>H : Move cursor below <num> lines from top of view.
func (ed *Editor) MoveToBelowTopOfView(n int) {
	if n < 1 {
		ed.Error("MoveToBelowTopOfView: n < 1")
		return
	}
	ed.EnsureCommand()
	if len(ed.vMeta) < 1 {
		return
	}
	i := n - 1
	b := ed.Buffer()
	b.Loc = ed.vMeta[i].Loc
	if b.Loc.Col < 1 {
		ed.toNonBlankCol()
	}
}

// <num>L : Move cursor above <num> lines from bottom of view.
func (ed *Editor) MoveToAboveBottomOfView(n int) {
	if n < 1 {
		ed.Error("MoveToAboveBottomOfView: n < 1")
		return
	}
	ed.EnsureCommand()
	i := len(ed.vMeta) - n
	if i < 0 {
		return
	}
	b := ed.Buffer()
	b.Loc = ed.vMeta[i].Loc
	if b.Loc.Col < 1 {
		ed.toNonBlankCol()
	}
}
