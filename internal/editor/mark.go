package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
)

//////////////////////
// Marking Commands //
//////////////////////

//
// Set Mark / Move to Mark
//

// m<char> : Mark current cursor position labelled by <char>.
func (ed *Editor) Mark(r rune) {
	ed.Buf().Mark(r)
}

// `<char> : Move cursor to marked position labelled by <char>.
func (ed *Editor) MoveToMark(r rune) (buf.Loc, bool) {
	b := ed.Buf()
	loc, ok := b.Marks[r]
	if !ok {
		ed.Notice("Mark not found")
		return buf.Loc{}, false
	}
	loc = b.Confine(loc)
	return loc, true
}

// '<char> : Move cursor to marked line labelled by <char>.
func (ed *Editor) MoveToMarkLine(r rune) (buf.Loc, bool) {
	b := ed.Buf()
	loc, ok := b.Marks[r]
	if !ok {
		ed.Notice("Mark not found")
		return buf.Loc{}, false
	}
	loc = b.Confine(loc)
	loc.Col = b.NonBlankColOfLine(loc.Row)
	return loc, true
}

//
// Move by Context
//

// “ : Move cursor to previous position in context.
func (ed *Editor) BackToMark() (buf.Loc, bool) {
	ed.Unimplemented("MoveBackToMark")
	return buf.Loc{}, false
}

// ” : Move cursor to previous line in context.
func (ed *Editor) BackToMarkLine() (buf.Loc, bool) {
	ed.Unimplemented("MarkBackToLine")
	return buf.Loc{}, false
}
