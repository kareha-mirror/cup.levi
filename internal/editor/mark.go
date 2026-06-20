package editor

//////////////////////
// Marking Commands //
//////////////////////

//
// Set Mark / Move to Mark
//

// m<letter> : Mark current cursor position labelled by <letter>.
func (ed *Editor) MarkSet(letter rune) {
	ed.EnsureCommand()
	ed.Buffer().Mark(letter)
}

// `<letter> : Move cursor to marked position labelled by <letter>.
func (ed *Editor) MarkMoveTo(letter rune) {
	ed.EnsureCommand()
	b := ed.Buffer()
	loc, ok := b.Marks[letter]
	if !ok {
		return
	}
	b.Loc = loc
	b.Confine()
}

// '<letter> : Move cursor to marked line labelled by <letter>.
func (ed *Editor) MarkMoveToLine(letter rune) {
	ed.EnsureCommand()
	b := ed.Buffer()
	loc, ok := b.Marks[letter]
	if !ok {
		return
	}
	b.Loc = loc
	b.Confine()
	ed.toNonBlankCol()
}

//
// Move by Context
//

// “ : Move cursor to previous position in context.
func (ed *Editor) MarkBack() {
	ed.EnsureCommand()
	ed.Unimplemented("MarkBack")
}

// ” : Move cursor to previous line in context.
func (ed *Editor) MarkBackToLine() {
	ed.EnsureCommand()
	ed.Unimplemented("MarkBackToLine")
}
