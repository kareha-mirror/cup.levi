package editor

/////////////////////
// Motion Commands //
/////////////////////

//
// Move by Character / Move by Line
//

// h : Move cursor left by character.
func (ed *Editor) MoveLeft(n int) {
	ed.EnsureCommand()
	ed.col -= n
	ed.Confine()
}

// j : Move cursor down by line.
func (ed *Editor) MoveDown(n int) {
	ed.EnsureCommand()
	ed.row += n
	ed.Confine()
}

// k : Move cursor up by line.
func (ed *Editor) MoveUp(n int) {
	ed.EnsureCommand()
	ed.row -= n
	ed.Confine()
}

// l : Move cursor right by character.
func (ed *Editor) MoveRight(n int) {
	ed.EnsureCommand()
	ed.col += n
	ed.Confine()
}

//
// Move in Line
//

// 0 : Move cursor to start of current line.
func (ed *Editor) MoveToStart() {
	ed.EnsureCommand()
	ed.col = 0
	ed.Confine() // redundant
}

// $ : Move cursor to end of current line.
func (ed *Editor) MoveToEnd() {
	ed.EnsureCommand()
	ed.col = ed.RuneCount() - 1
	ed.Confine()
}

// ^ : Move cursor to first non-blank character of current line.
func (ed *Editor) MoveToNonBlank() {
	ed.EnsureCommand()
	line := ed.CurrentLine()
	i := 0
	for _, r := range line {
		if r != ' ' && r != '\t' {
			break
		}
		i++
	}
	ed.col = i
	ed.Confine()
}

// <num>| : Move cursor to column <num> of current line.
// (Note: Proper vi's column number is visual-based, but levi' is rune-based.)
func (ed *Editor) MoveToColumn(n int) {
	ed.EnsureCommand()
	ed.col = n - 1
	ed.Confine()
}

//
// Move by Word / Move by Loose Word
//

// w : Move cursor forward by word.
func (ed *Editor) MoveByWord(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveByWord")
}

// b : Move cursor backward by word.
func (ed *Editor) MoveBackwardByWord(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardByWord")
}

// e : Move cursor to end of word.
func (ed *Editor) MoveToEndOfWord(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveToEndOfWord")
}

// W : Move cursor forward by loose word.
func (ed *Editor) MoveByLooseWord(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveByLooseWord")
}

// B : Move cursor backward by loose word.
func (ed *Editor) MoveBackwardByLooseWord(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardByLooseWord")
}

// E : Move cursor to end of loose word.
func (ed *Editor) MoveToEndOfLooseWord(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveToEndOfLooseWord")
}

//
// Move by Line
//

// Enter, + : Move cursor to first non-blank character of next line.
func (ed *Editor) MoveToNonBlankOfNextLine(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveToNonBlankOfNextLine")
}

// - : Move cursor to first non-blank character of previous line.
func (ed *Editor) MoveToNonBlankOfPrevLine(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveToNonBlankOfPrevLine")
}

// G : Move cursor to last line.
func (ed *Editor) MoveToLastLine() {
	ed.EnsureCommand()
	ed.Unimplemented("MoveToLastLine")
}

// <num>G : Move cursor to line <num>.
func (ed *Editor) MoveToLine(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveToLine")
}

//
// Move by Block
//

// ) : Move cursor forward by sentence.
func (ed *Editor) MoveBySentence(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveBySentence")
}

// ( : Move cursor backward by sentence.
func (ed *Editor) MoveBackwardBySentence(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardBySentence")
}

// } : Move cursor forward by paragraph.
func (ed *Editor) MoveByParagraph(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveByParagraph")
}

// { : Move cursor backward by paragraph.
func (ed *Editor) MoveBackwardByParagraph(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardByParagraph")
}

// ]] : Move cursor forward by section.
func (ed *Editor) MoveBySection(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveBySection")
}

// [[ : Move cursor backward by section.
func (ed *Editor) MoveBackwardBySection(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackwardBySection")
}

//
// Move in View
//

// H : Move cursor to top of view.
func (ed *Editor) MoveToTopOfView() {
	ed.EnsureCommand()
	ed.Unimplemented("MoveToTopOfView")
}

// M : Move cursor to middle of view.
func (ed *Editor) MoveToMiddleOfView() {
	ed.EnsureCommand()
	ed.Unimplemented("MoveToMiddleOfView")
}

// L : Move cursor to bottom of view.
func (ed *Editor) MoveToBottomOfView() {
	ed.EnsureCommand()
	ed.Unimplemented("MoveToBottomOfView")
}

// <num>H : Move cursor below <num> lines from top of view.
func (ed *Editor) MoveToBelowTopOfView(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveToBelowTopOfView")
}

// <num>L : Move cursor above <num> lines from bottom of view.
func (ed *Editor) MoveToAboveBottomOfView(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveToAboveBottomOfView")
}
