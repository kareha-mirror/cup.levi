package editor

func (ed *Editor) setRow(row int) bool {
	if row < 0 {
		return false
	}
	linesLen := len(ed.lines)
	if linesLen == 0 && row == 0 {
		ed.row = row
		return true
	}
	if row >= linesLen {
		return false
	}
	ed.row = row
	return true
}

func (ed *Editor) adjustRow(n int) bool {
	return ed.setRow(ed.row + n)
}

func (ed *Editor) confineRow() {
	n := len(ed.lines)
	if ed.row < 0 {
		ed.row = 0
	} else if ed.row >= n {
		ed.row = max(n-1, 0)
	}
}

func (ed *Editor) confineCol() {
	if ed.col < 0 {
		ed.col = 0
		return
	}
	rc := ed.RuneCount()
	if ed.col >= rc {
		ed.col = max(rc-1, 0)
	}
}

func (ed *Editor) confine() {
	ed.confineRow()
	ed.confineCol()
}

func (ed *Editor) saveVirtCol() {
	ed.virtCol = ed.col
}

func (ed *Editor) loadVirtCol() {
	ed.col = ed.virtCol
}

func (ed *Editor) setCol(col int) {
	ed.col = col
	ed.confineCol()
	ed.saveVirtCol()
}

func (ed *Editor) adjustCol(n int) {
	ed.setCol(ed.col + n)
}

func isBlank(r rune) bool {
	return r == ' ' || r == '\t'
}

func nonBlankCol(s string) int {
	i := 0
	for _, r := range s {
		if !isBlank(r) {
			break
		}
		i++
	}
	return i
}

func (ed *Editor) toNonBlankCol() {
	line := ed.CurrentLine()
	ed.setCol(nonBlankCol(line))
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
	ed.adjustCol(-n)
}

// j : Move cursor down by line.
func (ed *Editor) MoveDown(n int) {
	if n < 1 {
		ed.Error("MoveDown: n < 1")
		return
	}
	ed.EnsureCommand()
	if !ed.adjustRow(n) {
		return
	}
	ed.loadVirtCol()
	ed.confineCol()
}

// k : Move cursor up by line.
func (ed *Editor) MoveUp(n int) {
	if n < 1 {
		ed.Error("MoveUp: n < 1")
		return
	}
	ed.EnsureCommand()
	if !ed.adjustRow(-n) {
		return
	}
	ed.loadVirtCol()
	ed.confineCol()
}

// l : Move cursor right by character.
func (ed *Editor) MoveRight(n int) {
	if n < 1 {
		ed.Error("MoveRight: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.adjustCol(n)
}

//
// Move in Line
//

// 0 : Move cursor to start of current line.
func (ed *Editor) MoveToStart() {
	ed.EnsureCommand()
	ed.setCol(0)
}

// $ : Move cursor to end of current line.
func (ed *Editor) MoveToEnd() {
	ed.EnsureCommand()
	ed.setCol(ed.RuneCount() - 1)
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
	ed.setCol(n - 1)
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
	ed.Unimplemented("MoveByWord")
}

// b : Move cursor backward by word.
func (ed *Editor) MoveBackwardByWord(n int) {
	if n < 1 {
		ed.Error("MoveBackwardByWord: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardByWord")
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
	if !ed.adjustRow(n) {
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
	if !ed.adjustRow(-n) {
		return
	}
	ed.toNonBlankCol()
}

// G : Move cursor to first non-blank character of last line.
func (ed *Editor) MoveToLastLine() {
	ed.EnsureCommand()
	ed.row = len(ed.lines) - 1
	ed.confineRow()
	ed.toNonBlankCol()
}

// <num>G : Move cursor to first non-blank character of line specified by <num>.
func (ed *Editor) MoveToLine(n int) { // n: 1-based
	if n < 1 {
		ed.Error("MoveToLine: n < 1")
		return
	}
	ed.EnsureCommand()
	if !ed.setRow(n - 1) {
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
	ed.row = ed.vrow
	ed.confineRow()
	ed.toNonBlankCol()
}

// M : Move cursor to middle of view.
func (ed *Editor) MoveToMiddleOfView() {
	ed.EnsureCommand()
	ed.row = ed.vrow + ed.h/2 - 1
	ed.confineRow()
	ed.toNonBlankCol()
}

// L : Move cursor to bottom of view.
func (ed *Editor) MoveToBottomOfView() {
	ed.EnsureCommand()
	ed.row = ed.vrow + ed.h - 2
	ed.confineRow()
	ed.toNonBlankCol()
}

// <num>H : Move cursor below <num> lines from top of view.
func (ed *Editor) MoveToBelowTopOfView(n int) {
	if n < 1 {
		ed.Error("MoveToBelowTopOfView: n < 1")
		return
	}
	ed.EnsureCommand()
	if n-1 > ed.h-2 {
		ed.Ring("Out of range")
		return
	}
	ed.row = ed.vrow + n - 1
	ed.confineRow()
	ed.toNonBlankCol()
}

// <num>L : Move cursor above <num> lines from bottom of view.
func (ed *Editor) MoveToAboveBottomOfView(n int) {
	if n < 1 {
		ed.Error("MoveToAboveBottomOfView: n < 1")
		return
	}
	ed.EnsureCommand()
	if n-1 > ed.h-2 {
		ed.Ring("Out of range")
		return
	}
	ed.row = ed.vrow + ed.h - 2 - (n - 1)
	ed.confineRow()
	ed.toNonBlankCol()
}
