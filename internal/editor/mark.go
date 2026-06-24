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
	ed.EnsureCommand()
	ed.Buf().Mark(letter)
}

// `<letter> : Move cursor to marked position labelled by <letter>.
func (ed *Editor) MoveToMark(letter rune) (buf.Dest, bool) {
	ed.EnsureCommand()
	b := ed.Buf()
	loc, ok := b.Marks[letter]
	if !ok {
		return buf.Dest{}, false
	}
	b.Loc = loc
	b.Loc = b.Confine(b.Loc)
	return buf.Dest{}, false // TODO
}

// '<letter> : Move cursor to marked line labelled by <letter>.
func (ed *Editor) MoveToMarkLine(letter rune) (buf.Dest, bool) {
	ed.EnsureCommand()
	b := ed.Buf()
	loc, ok := b.Marks[letter]
	if !ok {
		return buf.Dest{}, false
	}
	b.Loc = loc
	b.Loc = b.Confine(b.Loc)
	ed.toNonBlankCol()
	return buf.Dest{}, false // TODO
}

//
// Move by Context
//

// “ : Move cursor to previous position in context.
func (ed *Editor) MoveBackToMark() (buf.Dest, bool) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackToMark")
	return buf.Dest{}, false // TODO
}

// ” : Move cursor to previous line in context.
func (ed *Editor) MoveBackToMarkLine() (buf.Dest, bool) {
	ed.EnsureCommand()
	ed.Unimplemented("MarkBackToLine")
	return buf.Dest{}, false // TODO
}
