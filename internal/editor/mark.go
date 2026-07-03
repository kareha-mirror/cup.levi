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
	b := ed.Buf()
	loc := b.Context
	bLoc := b.Loc
	b.Context = bLoc
	return loc, true
}

// ” : Move cursor to previous line in context.
func (ed *Editor) BackToMarkLine() (buf.Loc, bool) {
	b := ed.Buf()
	loc := b.Context
	loc.Col = b.NonBlankColOfLine(loc.Row)
	bLoc := b.Loc
	bLoc.Col = b.NonBlankColOfLine(bLoc.Row)
	b.Context = bLoc
	return loc, true
}
