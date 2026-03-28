package editor

///////////////////////////////////////////////
// Operator Commands (Copy / Delte / Change) //
///////////////////////////////////////////////

//
// Copy (Yank)
//

// yy, Y : Copy current line.
func (ed *Editor) OpCopyLine(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpCopyLine")
}

// y<mv> : Copy region from current cursor to destination of motion <mv>.
func (ed *Editor) OpCopyRegion(start Loc, end Loc) {
	ed.EnsureCommand()
	ed.Unimplemented("OpCopyRegion")
}

// y<mv> : Copy region from current cursor to destination of motion <mv>.
func (ed *Editor) OpCopyLineRegion(start int, end int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpCopyLineRegion")
}

// yw : Copy word.
func (ed *Editor) OpCopyWord(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpCopyWord")
}

// y$ : Copy to end of current line.
func (ed *Editor) OpCopyToEnd(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpCopyToEnd")
}

// "<reg>yy : Copy current line into register <reg>.
func (ed *Editor) OpCopyLineIntoReg(reg rune, n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpCopyLineIntoReg")
}

//
// Paste (Put)
//

// p : Paste after cursor.
func (ed *Editor) OpPaste(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpPaste")
}

// P : Paste before cursor.
func (ed *Editor) OpPasteBefore(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpPasteBefore")
}

// "<reg>p : Paste from register <reg>.
func (ed *Editor) OpPasteFromReg(reg rune, n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpPasteFromReg")
}

//
// Delete
//

// x : Delete character under cursor.
func (ed *Editor) OpDelete(n int) {
	ed.EnsureCommand()
	if len(ed.CurrentLine()) < 1 {
		ed.Ring("nothing to delete, line is empty")
		return
	}
	rs := []rune(ed.CurrentLine())
	if ed.col < 1 {
		ed.lines[ed.row] = string(rs[1:])
	} else {
		head := string(rs[:ed.col])
		tail := string(rs[ed.col+1:])
		ed.lines[ed.row] = head + tail
	}
	ed.Confine()
}

// X : Delete character before cursor.
func (ed *Editor) OpDeleteBefore(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpDeleteBefore")
}

// dd : Delete current line.
func (ed *Editor) OpDeleteLine(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpDeleteLine")
}

// d<mv> : Delete region from current cursor to destination of motion <mv>.
func (ed *Editor) OpDeleteRegion(start Loc, end Loc) {
	ed.EnsureCommand()
	ed.Unimplemented("OpDeleteRegion")
}

// d<mv> : Delete region from current cursor to destination of motion <mv>.
func (ed *Editor) OpDeleteLineRegion(start int, end int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpDeleteLineRegion")
}

// dw : Delete word.
func (ed *Editor) OpDeleteWord(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpDeleteWord")
}

// d$, D : Delete to end of current line.
func (ed *Editor) OpDeleteToEnd(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpDeleteToEnd")
}

//
// Change / Substitute
//

// cc : Change current line.
func (ed *Editor) OpChangeLine(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpChangeLine")
}

// c<mv> : Change region from current cursor to destination of motion <mv>.
func (ed *Editor) OpChangeRegion(start Loc, end Loc) {
	ed.EnsureCommand()
	ed.Unimplemented("OpChangeRegion")
}

// c<mv> : Change region from current cursor to destination of motion <mv>.
func (ed *Editor) OpChangeLineRegion(start int, end int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpChangeLineRegion")
}

// cw : Change word.
func (ed *Editor) OpChangeWord(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpChangeWord")
}

// C : Change to end of current line.
func (ed *Editor) OpChangeToEnd(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpChangeToEnd")
}

// s : Substitute one character under cursor.
func (ed *Editor) OpSubst(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpSubst")
}

// S : Substtute current line (equals cc).
func (ed *Editor) OpSubstLine(n int) {
	ed.OpChangeLine(n)
}
