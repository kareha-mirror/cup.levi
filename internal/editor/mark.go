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
	loc = b.Confine(loc)
	return buf.Dest{
		Loc:       loc,
		Linewise:  false,
		FreeCol:   false,
		Inclusive: false,
	}, true
}

// '<letter> : Move cursor to marked line labelled by <letter>.
func (ed *Editor) MoveToMarkLine(letter rune) (buf.Dest, bool) {
	ed.EnsureCommand()
	b := ed.Buf()
	loc, ok := b.Marks[letter]
	if !ok {
		return buf.Dest{}, false
	}
	loc = b.Confine(loc)
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
// Move by Context
//

// “ : Move cursor to previous position in context.
func (ed *Editor) MoveBackToMark() (buf.Dest, bool) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveBackToMark")
	return buf.Dest{}, false
}

// ” : Move cursor to previous line in context.
func (ed *Editor) MoveBackToMarkLine() (buf.Dest, bool) {
	ed.EnsureCommand()
	ed.Unimplemented("MarkBackToLine")
	return buf.Dest{}, false
}
