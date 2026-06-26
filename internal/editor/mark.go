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

// m<letter> : Mark current cursor position labelled by <letter>.
func (ed *Editor) MarkSet(letter rune) {
	ed.Commit()
	ed.Buf().Mark(letter)
}

// `<letter> : Move cursor to marked position labelled by <letter>.
func (ed *Editor) MoveToMark(letter rune) (buf.Loc, bool) {
	ed.Commit()
	b := ed.Buf()
	loc, ok := b.Marks[letter]
	if !ok {
		return buf.Loc{}, false
	}
	loc = b.Confine(loc)
	return loc, true
}

// '<letter> : Move cursor to marked line labelled by <letter>.
func (ed *Editor) MoveToMarkLine(letter rune) (buf.Loc, bool) {
	ed.Commit()
	b := ed.Buf()
	loc, ok := b.Marks[letter]
	if !ok {
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
func (ed *Editor) MoveBackToMark() (buf.Loc, bool) {
	ed.Commit()
	ed.Unimplemented("MoveBackToMark")
	return buf.Loc{}, false
}

// ” : Move cursor to previous line in context.
func (ed *Editor) MoveBackToMarkLine() (buf.Loc, bool) {
	ed.Commit()
	ed.Unimplemented("MarkBackToLine")
	return buf.Loc{}, false
}
