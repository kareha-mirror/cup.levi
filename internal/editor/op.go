package editor

import (
	"tea.kareha.org/cup/levi/internal/buffer"
)

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
	if b.Loc.Row+n > b.NumLines() {
		return
	}
	ed.killed.SetLines(b.Lines[b.Loc.Row : b.Loc.Row+n])
}

// y<mv> : Copy region from current cursor to destination of motion <mv>.
func (ed *Editor) OpCopyRegion(start buffer.Loc, end buffer.Loc) {
	ed.EnsureCommand()
	ed.Unimplemented("OpCopyRegion")
}

// y<mv> : Copy region from current cursor to destination of motion <mv>.
func (ed *Editor) OpCopyLineRegion(start int, end int) {
	ed.EnsureCommand()
	if end < start {
		start, end = end, start
	}
	b := ed.Buffer()
	if end+1 > b.NumLines() {
		return
	}
	ed.killed.SetLines(b.Lines[start : end+1])
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
		runes := []rune(b.CurrentLine())
		rs := []rune{}
		if b.Loc.Col+1 <= len(runes) {
			rs = append(rs, runes[:b.Loc.Col+1]...)
		}
		for i := 0; i < n; i++ {
			rs = append(rs, ed.killed.runes...)
		}
		if b.Loc.Col+1 < len(runes) {
			rs = append(rs, runes[b.Loc.Col+1:]...)
		}
		b.SetCurrentLine(string(rs))
		if len(runes) > 0 {
			b.Loc.Col++
		}
	case KillLines:
		lines := []string{}
		if b.Loc.Row+1 <= b.NumLines() {
			lines = append(lines, b.Lines[:b.Loc.Row+1]...)
		}
		for i := 0; i < n; i++ {
			lines = append(lines, ed.killed.lines...)
		}
		if b.Loc.Row+1 <= b.NumLines()-1 {
			lines = append(lines, b.Lines[b.Loc.Row+1:]...)
		}
		move := b.NumLines() > 0
		b.Lines = lines
		if move {
			ed.MoveByLine(1)
		}
	}
	b.Modified = true
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
		runes := []rune(b.CurrentLine())
		rs := append([]rune{}, runes[:b.Loc.Col]...)
		for i := 0; i < n; i++ {
			rs = append(rs, ed.killed.runes...)
		}
		rs = append(rs, runes[b.Loc.Col:]...)
		b.SetCurrentLine(string(rs))
	case KillLines:
		lines := append([]string{}, b.Lines[:b.Loc.Row]...)
		for i := 0; i < n; i++ {
			lines = append(lines, ed.killed.lines...)
		}
		lines = append(lines, b.Lines[b.Loc.Row:]...)
		b.Lines = lines
		ed.toNonBlankCol()
	}
	b.Modified = true
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
	n = min(n, len(rs)-b.Loc.Col)
	ed.killed.SetRunes(rs[b.Loc.Col : b.Loc.Col+n])
	if b.Loc.Col < 1 {
		b.SetCurrentLine(string(rs[n:]))
	} else {
		head := string(rs[:b.Loc.Col])
		tail := string(rs[b.Loc.Col+n:])
		b.SetCurrentLine(head + tail)
	}
	b.Confine()
	b.Modified = true
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
	if b.Loc.Row+n > b.NumLines() {
		return
	}
	lines := []string{}
	if b.Loc.Row > 0 {
		lines = append(lines, b.Lines[:b.Loc.Row]...)
	}
	ed.killed.SetLines(b.Lines[b.Loc.Row : b.Loc.Row+n])
	if b.Loc.Row+n <= b.NumLines()-1 {
		lines = append(lines, b.Lines[b.Loc.Row+n:]...)
	}
	b.Lines = lines
	b.Confine()
	b.Modified = true
}

// d<mv> : Delete region from current cursor to destination of motion <mv>.
func (ed *Editor) OpDeleteRegion(start buffer.Loc, end buffer.Loc) {
	ed.EnsureCommand()
	ed.Unimplemented("OpDeleteRegion")
}

// d<mv> : Delete region from current cursor to destination of motion <mv>.
func (ed *Editor) OpDeleteLineRegion(start int, end int) {
	ed.EnsureCommand()
	if end < start {
		start, end = end, start
	}
	b := ed.Buffer()
	if end+1 > b.NumLines() {
		return
	}
	lines := []string{}
	if start > 0 {
		lines = append(lines, b.Lines[:start]...)
	}
	ed.killed.SetLines(b.Lines[start : end+1])
	if end+1 <= b.NumLines()-1 {
		lines = append(lines, b.Lines[end+1:]...)
	}
	b.Lines = lines
	b.Confine()
	b.Modified = true
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
	b := ed.Buffer()
	rs := []rune(ed.CurrentLine())
	if b.Loc.Col < len(rs) {
		ed.killed.SetRunes(rs[b.Loc.Col:])
	}
	b.SetCurrentLine(string(rs[:b.Loc.Col]))
	b.Confine()
	b.Modified = true
	// TODO n
}

//
// Change / Substitute
//

// cc : Change current line.
func (ed *Editor) OpChangeLine(n int, replay bool) {
	if n < 1 {
		ed.Error("OpChangeLine: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("OpChangeLine")
}

// c<mv> : Change region from current cursor to destination of motion <mv>.
func (ed *Editor) OpChangeRegion(
	start buffer.Loc, end buffer.Loc, replay bool,
) {
	ed.EnsureCommand()
	ed.Unimplemented("OpChangeRegion")
}

// c<mv> : Change region from current cursor to destination of motion <mv>.
func (ed *Editor) OpChangeLineRegion(start int, end int, replay bool) {
	ed.EnsureCommand()
	ed.Unimplemented("OpChangeLineRegion")
}

// cw : Change word.
func (ed *Editor) OpChangeWord(n int, replay bool) {
	if n < 1 {
		ed.Error("OpChangeWord: n < 1")
		return
	}
	ed.EnsureCommand()
	ed.Unimplemented("OpChangeWord")
}

// C : Change to end of current line.
func (ed *Editor) OpChangeToEnd(n int, replay bool) {
	if n < 1 {
		ed.Error("OpChangeToEnd: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buffer()
	rs := []rune(ed.CurrentLine())
	if b.Loc.Col < len(rs) {
		ed.killed.SetRunes(rs[b.Loc.Col:])
	}
	line := string(rs[:b.Loc.Col])
	b.SetCurrentLine(line)
	ed.inp.Init(line, b.Loc.Col, ed.cfg.AutoIndent)
	ed.inpRow = b.Loc.Row
	ed.mode = ModeInsert
	// TODO n
}

// s : Substitute one character under cursor.
func (ed *Editor) OpSubst(n int, replay bool) {
	if n < 1 {
		ed.Error("OpSubst: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buffer()
	rs := []rune(ed.CurrentLine())
	nrs := append([]rune{}, rs[:b.Loc.Col]...)
	if b.Loc.Col+n <= len(rs)-1 {
		nrs = append(nrs, rs[b.Loc.Col+n:]...)
	}
	if replay {
		if len(ed.inserted) < 0 {
			return
		}
		ed.replayInsert(string(nrs))
	} else {
		ed.inp.Init(string(nrs), b.Loc.Col, ed.cfg.AutoIndent)
		ed.inpRow = b.Loc.Row
		ed.mode = ModeInsert
	}
}

// S : Substtute current line (equals cc).
func (ed *Editor) OpSubstLine(n int, replay bool) {
	if n < 1 {
		ed.Error("OpSubstLine: n < 1")
		return
	}
	ed.OpChangeLine(n, replay)
}
