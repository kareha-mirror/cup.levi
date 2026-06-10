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
	if n < 1 {
		return
	}
	if ed.row+n > len(ed.lines) {
		return
	}
	ed.killed.mode = KillLines
	ed.killed.lines = append([]string{}, ed.lines[ed.row:ed.row+n]...)
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
	if n < 1 {
		return
	}
	switch ed.killed.mode {
	case KillNone:
		ed.Ring("The default buffer is empty")
		return
	case KillRunes:
		runes := []rune(ed.Line(ed.row))
		runesLen := len(runes)
		for i := 0; i < n; i++ {
			rs := []rune{}
			if len(runes) > 0 && ed.col+1 <= len(runes) {
				rs = append(rs, runes[:ed.col+1]...)
			}
			rs = append(rs, ed.killed.runes...)
			if ed.col+1 < len(runes) {
				rs = append(rs, runes[ed.col+1:]...)
			}
			runes = rs
		}
		if runesLen > 0 {
			ed.col++
		}
		if len(ed.lines) < 1 {
			ed.lines = append(ed.lines, "")
		}
		ed.lines[ed.row] = string(runes)
	case KillLines:
		linesLen := len(ed.lines)
		for i := 0; i < n; i++ {
			lines := []string{}
			if ed.row+1 <= len(ed.lines) {
				lines = append(lines, ed.lines[:ed.row+1]...)
			}
			lines = append(lines, ed.killed.lines...)
			if ed.row+1 <= len(ed.lines)-1 {
				lines = append(lines, ed.lines[ed.row+1:]...)
			}
			ed.lines = lines
		}
		if linesLen > 0 {
			ed.MoveByLine(1)
		}
	}
	ed.modified = true
}

// P : Paste before cursor.
func (ed *Editor) OpPasteBefore(n int) {
	ed.EnsureCommand()
	if n < 1 {
		return
	}
	switch ed.killed.mode {
	case KillNone:
		ed.Ring("The default buffer is empty")
		return
	case KillRunes:
		runes := []rune(ed.Line(ed.row))
		for i := 0; i < n; i++ {
			rs := []rune{}
			rs = append(rs, runes[:ed.col]...)
			rs = append(rs, ed.killed.runes...)
			rs = append(rs, runes[ed.col:]...)
			runes = rs
		}
		if len(ed.lines) < 1 {
			ed.lines = append(ed.lines, "")
		}
		ed.lines[ed.row] = string(runes)
	case KillLines:
		for i := 0; i < n; i++ {
			lines := []string{}
			lines = append(lines, ed.lines[:ed.row]...)
			lines = append(lines, ed.killed.lines...)
			lines = append(lines, ed.lines[ed.row:]...)
			ed.lines = lines
			ed.MoveToNonBlank()
		}
	}
	ed.modified = true
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
	if n < 1 {
		return
	}
	if len(ed.CurrentLine()) < 1 {
		return
	}
	rs := []rune(ed.CurrentLine())
	n = min(n, len(rs)-ed.col)
	ed.killed.mode = KillRunes
	ed.killed.runes = append([]rune{}, rs[ed.col:ed.col+n]...)
	if ed.col < 1 {
		ed.lines[ed.row] = string(rs[n:])
	} else {
		head := string(rs[:ed.col])
		tail := string(rs[ed.col+n:])
		ed.lines[ed.row] = head + tail
	}
	ed.Confine()
	ed.modified = true
}

// X : Delete character before cursor.
func (ed *Editor) OpDeleteBefore(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("OpDeleteBefore")
}

// dd : Delete current line.
func (ed *Editor) OpDeleteLine(n int) {
	ed.EnsureCommand()
	if n < 1 {
		return
	}
	if ed.row+n > len(ed.lines) {
		return
	}
	lines := []string{}
	if ed.row > 0 {
		lines = append(lines, ed.lines[:ed.row]...)
	}
	ed.killed.mode = KillLines
	ed.killed.lines = append([]string{}, ed.lines[ed.row:ed.row+n]...)
	if ed.row+n <= len(ed.lines)-1 {
		lines = append(lines, ed.lines[ed.row+n:]...)
	}
	ed.lines = lines
	ed.Confine()
	ed.modified = true
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
