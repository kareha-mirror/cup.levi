package editor

///////////////////////////////////////////////
// Operator Commands (Copy / Delte / Change) //
///////////////////////////////////////////////

//
// Copy (Yank)
//

// yy, Y : Copy current line.
func (ed *Editor) OpCopyLine(n int) {
	if n < 1 {
		ed.Error("OpCopyLine: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buffer()
	if b.row+n > len(b.lines) {
		return
	}
	ed.killed.SetLines(b.lines[b.row : b.row+n])
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
	if n < 1 {
		ed.Error("OpCopyWord: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("OpCopyWord")
}

// y$ : Copy to end of current line.
func (ed *Editor) OpCopyToEnd(n int) {
	if n < 1 {
		ed.Error("OpCopyToEnd: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("OpCopyToEnd")
}

// "<reg>yy : Copy current line into register <reg>.
func (ed *Editor) OpCopyLineIntoReg(reg rune, n int) {
	if n < 1 {
		ed.Error("OpCopyLineIntoReg: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("OpCopyLineIntoReg")
}

//
// Paste (Put)
//

// p : Paste after cursor.
func (ed *Editor) OpPaste(n int) {
	if n < 1 {
		ed.Error("OpPaste: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buffer()
	if ed.killed.mode == KillNone {
		ed.Ring("The default buffer is empty")
		return
	}
	switch ed.killed.mode {
	case KillRunes:
		if len(b.lines) < 1 {
			b.lines = append(b.lines, "")
		}
		runes := []rune(b.lines[b.row])
		rs := []rune{}
		if b.col+1 <= len(runes) {
			rs = append(rs, runes[:b.col+1]...)
		}
		for i := 0; i < n; i++ {
			rs = append(rs, ed.killed.runes...)
		}
		if b.col+1 < len(runes) {
			rs = append(rs, runes[b.col+1:]...)
		}
		b.lines[b.row] = string(rs)
		if len(runes) > 0 {
			b.col++
		}
	case KillLines:
		lines := []string{}
		if b.row+1 <= len(b.lines) {
			lines = append(lines, b.lines[:b.row+1]...)
		}
		for i := 0; i < n; i++ {
			lines = append(lines, ed.killed.lines...)
		}
		if b.row+1 <= len(b.lines)-1 {
			lines = append(lines, b.lines[b.row+1:]...)
		}
		move := len(b.lines) > 0
		b.lines = lines
		if move {
			ed.MoveByLine(1)
		}
	}
	b.modified = true
}

// P : Paste before cursor.
func (ed *Editor) OpPasteBefore(n int) {
	if n < 1 {
		ed.Error("OpPasteBefore: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buffer()
	if ed.killed.mode == KillNone {
		ed.Ring("The default buffer is empty")
		return
	}
	switch ed.killed.mode {
	case KillRunes:
		if len(b.lines) < 1 {
			b.lines = append(b.lines, "")
		}
		runes := []rune(b.lines[b.row])
		rs := append([]rune{}, runes[:b.col]...)
		for i := 0; i < n; i++ {
			rs = append(rs, ed.killed.runes...)
		}
		rs = append(rs, runes[b.col:]...)
		b.lines[b.row] = string(rs)
	case KillLines:
		lines := append([]string{}, b.lines[:b.row]...)
		for i := 0; i < n; i++ {
			lines = append(lines, ed.killed.lines...)
		}
		lines = append(lines, b.lines[b.row:]...)
		b.lines = lines
		ed.toNonBlankCol()
	}
	b.modified = true
}

// "<reg>p : Paste from register <reg>.
func (ed *Editor) OpPasteFromReg(reg rune, n int) {
	if n < 1 {
		ed.Error("OpPasteFromReg: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("OpPasteFromReg")
}

//
// Delete
//

// x : Delete character under cursor.
func (ed *Editor) OpDelete(n int) {
	if n < 1 {
		ed.Error("OpDelete: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buffer()
	if len(ed.CurrentLine()) < 1 {
		return
	}
	rs := []rune(ed.CurrentLine())
	n = min(n, len(rs)-b.col)
	ed.killed.SetRunes(rs[b.col : b.col+n])
	if b.col < 1 {
		b.lines[b.row] = string(rs[n:])
	} else {
		head := string(rs[:b.col])
		tail := string(rs[b.col+n:])
		b.lines[b.row] = head + tail
	}
	ed.confine()
	b.modified = true
}

// X : Delete character before cursor.
func (ed *Editor) OpDeleteBefore(n int) {
	if n < 1 {
		ed.Error("OpDeleteBefore: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("OpDeleteBefore")
}

// dd : Delete current line.
func (ed *Editor) OpDeleteLine(n int) {
	if n < 1 {
		ed.Error("OpDeleteLine: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buffer()
	if b.row+n > len(b.lines) {
		return
	}
	lines := []string{}
	if b.row > 0 {
		lines = append(lines, b.lines[:b.row]...)
	}
	ed.killed.SetLines(b.lines[b.row : b.row+n])
	if b.row+n <= len(b.lines)-1 {
		lines = append(lines, b.lines[b.row+n:]...)
	}
	b.lines = lines
	ed.confine()
	b.modified = true
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
	if n < 1 {
		ed.Error("OpDeleteWord: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("OpDeleteWord")
}

// d$, D : Delete to end of current line.
func (ed *Editor) OpDeleteToEnd(n int) {
	if n < 1 {
		ed.Error("OpDeleteToEnd: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("OpDeleteToEnd")
}

//
// Change / Substitute
//

// cc : Change current line.
func (ed *Editor) OpChangeLine(n int) {
	if n < 1 {
		ed.Error("OpChangeLine: n < 1")
		return
	}
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
	if n < 1 {
		ed.Error("OpChangeWord: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("OpChangeWord")
}

// C : Change to end of current line.
func (ed *Editor) OpChangeToEnd(n int) {
	if n < 1 {
		ed.Error("OpChangeToEnd: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("OpChangeToEnd")
}

// s : Substitute one character under cursor.
func (ed *Editor) OpSubst(n int) {
	if n < 1 {
		ed.Error("OpSubst: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buffer()
	rs := []rune(ed.CurrentLine())
	nrs := append([]rune{}, rs[:b.col]...)
	if b.col+n <= len(rs)-1 {
		nrs = append(nrs, rs[b.col+n:]...)
	}
	ed.inp.Init(string(nrs), b.col)
	ed.inpRow = b.row
	ed.mode = ModeInsert
}

// S : Substtute current line (equals cc).
func (ed *Editor) OpSubstLine(n int) {
	if n < 1 {
		ed.Error("OpSubstLine: n < 1")
		return
	}
	ed.OpChangeLine(n)
}
