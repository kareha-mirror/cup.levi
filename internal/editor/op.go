package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
)

///////////////////////////////////////////////
// Operator Commands (Copy / Delte / Change) //
///////////////////////////////////////////////

//
// Copy (Yank)
//

// "<reg>yy, "<reg>Y : Copy current line into register <reg>.
func (ed *Editor) OpCopyLine(reg string, n int) {
	if n < 1 {
		ed.Error("OpCopyLine: n < 1")
		return
	}
	b := ed.Buf()
	if b.Loc.Row+n > b.NumLines() {
		ed.Notice("Out of range")
		return
	}
	ed.ApplyRegLines(reg, b.Lines[b.Loc.Row:b.Loc.Row+n])
}

// y<mv> : Copy region from current cursor to destination of motion <mv>.
func (ed *Editor) OpCopyRegion(
	reg string, start buf.Loc, end buf.Loc, inclusive bool,
) {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, inclusive)
	lines := b.RegionRunewise(start, end)
	ed.ApplyRegRunes(reg, lines)
	b.Loc = start
}

// y<mv> : Copy region from current cursor to destination of motion <mv>.
func (ed *Editor) OpCopyLineRegion(
	reg string, start buf.Loc, end buf.Loc,
) {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, true)
	if end.Row+1 > b.NumLines() {
		ed.Notice("Out of range")
		return
	}
	ed.ApplyRegLines(reg, b.Lines[start.Row:end.Row+1])
	b.Loc = start
}

// yw : Copy word.
func (ed *Editor) OpCopyWord(reg string, n int) {
	if n < 1 {
		ed.Error("OpCopyWord: n < 1")
		return
	}
	ed.Unimplemented("OpCopyWord")
}

// y$ : Copy to end of current line.
func (ed *Editor) OpCopyToEnd(reg string, n int) {
	if n < 1 {
		ed.Error("OpCopyToEnd: n < 1")
		return
	}
	ed.Unimplemented("OpCopyToEnd")
}

//
// Paste (Put)
//

// "<reg>p : Paste after cursor from register <reg>.
func (ed *Editor) OpPaste(reg string, n int) {
	if n < 1 {
		ed.Error("OpPaste: n < 1")
		return
	}
	if ed.RegMode(reg) == KillNone {
		if reg == "" {
			ed.Ring("The default buffer is empty")
		} else {
			ed.Ring("Buffer %s is empty", reg)
		}
		return
	}
	killed := ed.RegKilled(reg)
	b := ed.Buf()
	switch ed.RegMode(reg) {
	case KillRunes:
		if len(killed) < 2 {
			runes := []rune(b.CurrentLine())
			rs := []rune{}
			if b.Loc.Col+1 <= len(runes) {
				rs = append(rs, runes[:b.Loc.Col+1]...)
			}
			for i := 0; i < n; i++ {
				rs = append(rs, []rune(killed[0])...)
			}
			if b.Loc.Col+1 < len(runes) {
				rs = append(rs, runes[b.Loc.Col+1:]...)
			}
			b.SetCurrentLine(string(rs))
			if len(runes) > 0 {
				b.Loc.Col++
				b.VirtCol = b.Loc.Col
			}
		} else {
			lines := []string{}
			lines = append(lines, b.Lines[:b.Loc.Row]...)

			runes := []rune(b.CurrentLine())
			rs := []rune{}
			if b.Loc.Col+1 <= len(runes) {
				rs = append(rs, runes[:b.Loc.Col+1]...)
			}
			rs = append(rs, []rune(killed[0])...)
			lines = append(lines, string(rs))

			if len(killed) > 2 {
				lines = append(lines, killed[1:len(killed)-1]...)
			}
			rs = []rune(killed[len(killed)-1])
			if b.Loc.Col+1 < len(runes) {
				rs = append(rs, runes[b.Loc.Col+1:]...)
			}
			lines = append(lines, string(rs))

			if b.Loc.Row+1 <= b.NumLines()-1 {
				lines = append(lines, b.Lines[b.Loc.Row+1:]...)
			}

			b.Lines = lines
		}
	case KillLines:
		lines := []string{}
		if b.Loc.Row+1 <= b.NumLines() {
			lines = append(lines, b.Lines[:b.Loc.Row+1]...)
		}
		for i := 0; i < n; i++ {
			lines = append(lines, killed...)
		}
		if b.Loc.Row+1 <= b.NumLines()-1 {
			lines = append(lines, b.Lines[b.Loc.Row+1:]...)
		}
		move := b.NumLines() > 0
		b.Lines = lines
		if move {
			b.Loc.Row++
			b.Loc = b.Confine(b.Loc)
			b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
			b.VirtCol = b.Loc.Col
		}
	}
	b.Modified = true
}

// "<reg>p : Paste before cursor from register <reg>.
func (ed *Editor) OpPasteBefore(reg string, n int) {
	if n < 1 {
		ed.Error("OpPasteBefore: n < 1")
		return
	}
	b := ed.Buf()
	if ed.RegMode(reg) == KillNone {
		if reg == "" {
			ed.Ring("The default buffer is empty")
		} else {
			ed.Ring("Buffer %s is empty", reg)
		}
		return
	}
	killed := ed.RegKilled(reg)
	switch ed.RegMode(reg) {
	case KillRunes:
		if len(killed) < 2 {
			runes := []rune(b.CurrentLine())
			rs := append([]rune{}, runes[:b.Loc.Col]...)
			for i := 0; i < n; i++ {
				rs = append(rs, []rune(killed[0])...)
			}
			rs = append(rs, runes[b.Loc.Col:]...)
			b.SetCurrentLine(string(rs))
		} else {
			lines := append([]string{}, b.Lines[:b.Loc.Row]...)

			runes := []rune(b.CurrentLine())
			rs := append([]rune{}, runes[:b.Loc.Col]...)
			rs = append(rs, []rune(killed[0])...)
			lines = append(lines, string(rs))

			if len(killed) > 2 {
				lines = append(
					lines, killed[1:len(killed)-1]...,
				)
			}

			rs = []rune(killed[len(killed)-1])
			rs = append(rs, runes[b.Loc.Col:]...)
			lines = append(lines, string(rs))

			if b.Loc.Row+1 < b.NumLines() {
				lines = append(lines, b.Lines[b.Loc.Row+1:]...)
			}

			b.Lines = lines
		}
	case KillLines:
		lines := append([]string{}, b.Lines[:b.Loc.Row]...)
		for i := 0; i < n; i++ {
			lines = append(lines, killed...)
		}
		lines = append(lines, b.Lines[b.Loc.Row:]...)
		b.Lines = lines
		b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
		b.VirtCol = b.Loc.Col
	}
	b.Modified = true
}

//
// Delete
//

func (ed *Editor) internalOpDelete(reg string, n int) bool {
	b := ed.Buf()
	if len(b.CurrentLine()) < 1 {
		return false
	}
	rs := []rune(b.CurrentLine())
	n = min(n, len(rs)-b.Loc.Col)
	ed.ApplyRegRunes(reg, []string{string(rs[b.Loc.Col : b.Loc.Col+n])})
	if b.Loc.Col < 1 {
		b.SetCurrentLine(string(rs[n:]))
	} else {
		head := string(rs[:b.Loc.Col])
		tail := string(rs[b.Loc.Col+n:])
		b.SetCurrentLine(head + tail)
	}
	return true
}

// x : Delete character under cursor.
func (ed *Editor) OpDelete(reg string, n int) {
	if n < 1 {
		ed.Error("OpDelete: n < 1")
		return
	}
	if !ed.internalOpDelete(reg, n) {
		ed.Notice("Nothing to delete")
		return
	}
	b := ed.Buf()
	b.Loc = b.ConfineInclusive(b.Loc)
	b.Modified = true
}

// X : Delete character before cursor.
func (ed *Editor) OpDeleteBefore(reg string, n int) {
	if n < 1 {
		ed.Error("OpDeleteBefore: n < 1")
		return
	}
	ed.Unimplemented("OpDeleteBefore")
}

// dd : Delete current line.
func (ed *Editor) OpDeleteLine(reg string, n int) {
	if n < 1 {
		ed.Error("OpDeleteLine: n < 1")
		return
	}
	b := ed.Buf()
	if b.Loc.Row+n > b.NumLines() {
		ed.Notice("Out of range")
		return
	}
	lines := append([]string{}, b.Lines[:b.Loc.Row]...)
	ed.ApplyRegLines(reg, b.Lines[b.Loc.Row:b.Loc.Row+n])
	if b.Loc.Row+n < b.NumLines() {
		lines = append(lines, b.Lines[b.Loc.Row+n:]...)
	}
	b.Lines = lines
	b.Loc = b.ConfineInclusive(b.Loc)
	b.Modified = true
}

// d<mv> : Delete region from current cursor to destination of motion <mv>.
func (ed *Editor) OpDeleteRegion(
	reg string, start buf.Loc, end buf.Loc, inclusive bool,
) {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, inclusive)
	lines := b.RegionRunewise(start, end)
	ed.ApplyRegRunes(reg, lines)

	lines = append([]string{}, b.Lines[:start.Row]...)

	runes := []rune(b.Line(start.Row))
	rs := append([]rune{}, runes[:start.Col]...)
	runes = []rune(b.Line(end.Row))
	rs = append(rs, runes[end.Col:]...)
	lines = append(lines, string(rs))

	if end.Row+1 < b.NumLines() {
		lines = append(lines, b.Lines[end.Row+1:]...)
	}
	b.Lines = lines
	b.Loc = start
	b.Loc = b.ConfineInclusive(b.Loc)
	b.Modified = true
}

// d<mv> : Delete region from current cursor to destination of motion <mv>.
func (ed *Editor) OpDeleteLineRegion(
	reg string, start buf.Loc, end buf.Loc,
) {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, true)
	if end.Row+1 > b.NumLines() {
		ed.Notice("Out of range")
		return
	}
	lines := append([]string{}, b.Lines[:start.Row]...)
	ed.ApplyRegLines(reg, b.Lines[start.Row:end.Row+1])
	if end.Row+1 < b.NumLines() {
		lines = append(lines, b.Lines[end.Row+1:]...)
	}
	b.Lines = lines
	b.Loc = start
	b.Loc = b.ConfineInclusive(b.Loc)
	b.Modified = true
}

// dw : Delete word.
func (ed *Editor) OpDeleteWord(reg string, n int) {
	if n < 1 {
		ed.Error("OpDeleteWord: n < 1")
		return
	}
	b := ed.Buf()
	start := b.Loc
	end, ok := ed.MoveByWord(n)
	if !ok {
		ed.Error("Failed to move")
		return
	}
	ed.OpDeleteRegion(reg, start, end, false)
	// TODO n
}

// d$, D : Delete to end of current line.
func (ed *Editor) OpDeleteToEnd(reg string, n int) {
	if n < 1 {
		ed.Error("OpDeleteToEnd: n < 1")
		return
	}
	b := ed.Buf()
	rs := []rune(b.CurrentLine())
	if b.Loc.Col < len(rs) {
		ed.ApplyRegRunes(reg, []string{string(rs[b.Loc.Col:])})
	}
	b.SetCurrentLine(string(rs[:b.Loc.Col]))
	b.Loc = b.ConfineInclusive(b.Loc)
	b.Modified = true
	// TODO n
}

//
// Change / Substitute
//

// cc : Change current line.
func (ed *Editor) OpChangeLine(reg string, n int, replay bool) {
	if n < 1 {
		ed.Error("OpChangeLine: n < 1")
		return
	}
	ed.Unimplemented("OpChangeLine")
}

// c<mv> : Change region from current cursor to destination of motion <mv>.
func (ed *Editor) OpChangeRegion(
	reg string, start buf.Loc, end buf.Loc, inclusive bool, replay bool,
) {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, inclusive)
	after := false
	if end.Row < b.NumLines() {
		rc := utf8.RuneCountInString(b.Line(end.Row))
		if end.Col >= rc {
			after = true
		}
	}
	ed.OpDeleteRegion(reg, start, end, inclusive)
	if after {
		b.Loc.Col++
	}
	if replay {
		if len(ed.inserted) < 0 {
			ed.Notice("Not inserted yet")
			return
		}
		ed.replayInsert()
		b.Loc.Col--
		b.Loc = b.ConfineInclusive(b.Loc)
		b.VirtCol = b.Loc.Col
		b.Modified = true
	} else {
		ed.inp.Init(b.CurrentLine(), b.Loc.Col, ed.cfg.AutoIndent)
		ed.inpRow = b.Loc.Row
		ed.mode = ModeInsert
	}
}

// c<mv> : Change region from current cursor to destination of motion <mv>.
func (ed *Editor) OpChangeLineRegion(
	reg string, start buf.Loc, end buf.Loc, replay bool,
) {
	ed.Unimplemented("OpChangeLineRegion")
}

// cw : Change word.
func (ed *Editor) OpChangeWord(reg string, n int, replay bool) {
	if n < 1 {
		ed.Error("OpChangeWord: n < 1")
		return
	}
	start := ed.Buf().Loc
	end, ok := ed.MoveByWordEx(n)
	if !ok {
		ed.Error("Failed to move")
		return
	}
	ed.OpChangeRegion(reg, start, end, false, replay)
	// TODO n
}

// C : Change to end of current line.
func (ed *Editor) OpChangeToEnd(reg string, n int, replay bool) {
	if n < 1 {
		ed.Error("OpChangeToEnd: n < 1")
		return
	}
	b := ed.Buf()
	rs := []rune(b.CurrentLine())
	if b.Loc.Col < len(rs) {
		ed.ApplyRegRunes(reg, []string{string(rs[b.Loc.Col:])})
	}
	line := string(rs[:b.Loc.Col])
	b.SetCurrentLine(line)
	ed.inp.Init(line, b.Loc.Col, ed.cfg.AutoIndent)
	ed.inpRow = b.Loc.Row
	ed.mode = ModeInsert
	// TODO n
}

// s : Substitute one character under cursor.
func (ed *Editor) OpSubst(reg string, n int, replay bool) {
	if n < 1 {
		ed.Error("OpSubst: n < 1")
		return
	}
	ed.internalOpDelete(reg, n)
	b := ed.Buf()
	if replay {
		if len(ed.inserted) < 0 {
			ed.Notice("Not inserted yet")
			return
		}
		ed.replayInsert()
		b.Loc.Col--
		b.Loc = b.ConfineInclusive(b.Loc)
		b.VirtCol = b.Loc.Col
		b.Modified = true
	} else {
		ed.inp.Init(b.CurrentLine(), b.Loc.Col, ed.cfg.AutoIndent)
		ed.inpRow = b.Loc.Row
		ed.mode = ModeInsert
	}
}

// S : Substtute current line (equals cc).
func (ed *Editor) OpSubstLine(reg string, n int, replay bool) {
	if n < 1 {
		ed.Error("OpSubstLine: n < 1")
		return
	}
	ed.OpChangeLine(reg, n, replay)
}
