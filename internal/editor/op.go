package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
)

func orderRegion(start, end buf.Loc) (buf.Loc, buf.Loc) {
	if start.Row < end.Row {
		return start, end
	}
	if end.Row < start.Row {
		return end, start
	}
	// start.Row == end.Row
	if start.Col < end.Col {
		return start, end
	}
	return end, start
}

func (ed *Editor) confineLoc(loc buf.Loc) buf.Loc {
	b := ed.Buf()
	if loc.Row < 0 {
		return buf.Loc{0, 0}
	}
	if loc.Row > b.NumLines() {
		return buf.Loc{0, b.NumLines()}
	}
	if loc.Col < 0 {
		return buf.Loc{0, loc.Row}
	}
	line := b.Line(loc.Row)
	rc := utf8.RuneCountInString(line)
	if loc.Col > rc {
		return buf.Loc{rc, loc.Row}
	}
	return loc
}

func (ed *Editor) confineRegion(
	start, end buf.Loc, inclusive bool,
) (buf.Loc, buf.Loc) {
	start, end = orderRegion(start, end)
	return ed.confineLoc(start), ed.confineLoc(end)
}

func (ed *Editor) getRegion(start, end buf.Loc) []string {
	b := ed.Buf()
	if start.Row == end.Row {
		line := b.Line(start.Row)
		if line == "" {
			return []string{""}
		}
		rs := []rune(line)
		return []string{string(rs[start.Col:end.Col])}
	}
	lines := []string{}
	line := b.Line(start.Row)
	if line == "" {
		lines = append(lines, "")
	} else {
		rs := []rune(line)
		lines = append(lines, string(rs[start.Col:]))
	}
	for row := start.Row + 1; row < end.Row; row++ {
		lines = append(lines, b.Line(row))
	}
	line = b.Line(end.Row)
	if line == "" {
		lines = append(lines, "")
	} else {
		rs := []rune(line)
		lines = append(lines, string(rs[:end.Col]))
	}
	return lines
}

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
	b := ed.Buf()
	if b.Loc.Row+n > b.NumLines() {
		return
	}
	ed.regs.SetLines("", b.Lines[b.Loc.Row:b.Loc.Row+n])
}

// y<mv> : Copy region from current cursor to destination of motion <mv>.
func (ed *Editor) OpCopyRegion(start buf.Loc, end buf.Loc, inclusive bool) {
	ed.EnsureCommand()
	start, end = ed.confineRegion(start, end, inclusive)
	lines := ed.getRegion(start, end)
	ed.regs.SetRunes("", lines)
}

// y<mv> : Copy region from current cursor to destination of motion <mv>.
func (ed *Editor) OpCopyLineRegion(start int, end int) {
	ed.EnsureCommand()
	if end < start {
		start, end = end, start
	}
	b := ed.Buf()
	if end+1 > b.NumLines() {
		return
	}
	ed.regs.SetLines("", b.Lines[start:end+1])
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
	b := ed.Buf()
	if ed.regs.Mode("") == KillNone {
		ed.Ring("The default buffer is empty")
		return
	}
	killed := ed.regs.Killed("")
	switch ed.regs.Mode("") {
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
				lines = append(
					lines, killed[1:len(killed)-2]...,
				)
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
	b := ed.Buf()
	if ed.regs.Mode("") == KillNone {
		ed.Ring("The default buffer is empty")
		return
	}
	killed := ed.regs.Killed("")
	switch ed.regs.Mode("") {
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
			lines := []string{}
			lines = append(lines, b.Lines[:b.Loc.Row]...)

			runes := []rune(b.CurrentLine())
			rs := []rune{}
			rs = append(rs, runes[:b.Loc.Col]...)
			rs = append(rs, []rune(killed[0])...)
			lines = append(lines, string(rs))

			if len(killed) > 2 {
				lines = append(
					lines, killed[1:len(killed)-2]...,
				)
			}

			rs = []rune(killed[len(killed)-1])
			rs = append(rs, runes[b.Loc.Col:]...)
			lines = append(lines, string(rs))

			if b.Loc.Row+1 <= b.NumLines()-1 {
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
	b := ed.Buf()
	if len(ed.CurrentLine()) < 1 {
		return
	}
	rs := []rune(ed.CurrentLine())
	n = min(n, len(rs)-b.Loc.Col)
	ed.regs.SetRunes("", []string{string(rs[b.Loc.Col : b.Loc.Col+n])})
	if b.Loc.Col < 1 {
		b.SetCurrentLine(string(rs[n:]))
	} else {
		head := string(rs[:b.Loc.Col])
		tail := string(rs[b.Loc.Col+n:])
		b.SetCurrentLine(head + tail)
	}
	b.Loc = b.ConfineInclusive(b.Loc)
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
	b := ed.Buf()
	if b.Loc.Row+n > b.NumLines() {
		return
	}
	lines := []string{}
	if b.Loc.Row > 0 {
		lines = append(lines, b.Lines[:b.Loc.Row]...)
	}
	ed.regs.SetLines("", b.Lines[b.Loc.Row:b.Loc.Row+n])
	if b.Loc.Row+n <= b.NumLines()-1 {
		lines = append(lines, b.Lines[b.Loc.Row+n:]...)
	}
	b.Lines = lines
	b.Loc = b.ConfineInclusive(b.Loc)
	b.Modified = true
}

// d<mv> : Delete region from current cursor to destination of motion <mv>.
func (ed *Editor) OpDeleteRegion(start buf.Loc, end buf.Loc, inclusive bool) {
	ed.EnsureCommand()
	start, end = ed.confineRegion(start, end, inclusive)
	lines := ed.getRegion(start, end)
	ed.regs.SetRunes("", lines)

	b := ed.Buf()
	lines = []string{}
	if start.Row > 0 {
		lines = append(lines, b.Lines[:start.Row]...)
	}

	runes := []rune(b.Line(start.Row))
	rs := append([]rune{}, runes[:start.Col]...)
	runes = []rune(b.Line(end.Row))
	rs = append(rs, runes[end.Col:]...)
	lines = append(lines, string(rs))

	if end.Row+1 <= b.NumLines()-1 {
		lines = append(lines, b.Lines[end.Row+1:]...)
	}
	b.Lines = lines
	b.Loc = start
	b.Loc = b.ConfineInclusive(b.Loc)
	b.Modified = true
}

// d<mv> : Delete region from current cursor to destination of motion <mv>.
func (ed *Editor) OpDeleteLineRegion(start int, end int) {
	ed.EnsureCommand()
	if end < start {
		start, end = end, start
	}
	b := ed.Buf()
	if end+1 > b.NumLines() {
		return
	}
	lines := []string{}
	if start > 0 {
		lines = append(lines, b.Lines[:start]...)
	}
	ed.regs.SetLines("", b.Lines[start:end+1])
	if end+1 <= b.NumLines()-1 {
		lines = append(lines, b.Lines[end+1:]...)
	}
	b.Lines = lines
	b.Loc.Row = start
	b.Loc.Col = 0
	b.Loc = b.ConfineInclusive(b.Loc)
	b.Modified = true
}

// dw : Delete word.
func (ed *Editor) OpDeleteWord(n int) {
	if n < 1 {
		ed.Error("OpDeleteWord: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buf()
	start := b.Loc
	end, ok := ed.MoveByWordEx(n)
	if !ok {
		return
	}
	ed.OpDeleteRegion(start, end, false)
	// TODO n
}

// d$, D : Delete to end of current line.
func (ed *Editor) OpDeleteToEnd(n int) {
	if n < 1 {
		ed.Error("OpDeleteToEnd: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buf()
	rs := []rune(ed.CurrentLine())
	if b.Loc.Col < len(rs) {
		ed.regs.SetRunes("", []string{string(rs[b.Loc.Col:])})
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
	start buf.Loc, end buf.Loc, inclusive bool, replay bool,
) {
	ed.EnsureCommand()
	start, end = ed.confineRegion(start, end, inclusive)
	after := false
	b := ed.Buf()
	if end.Row < b.NumLines() {
		line := b.Line(end.Row)
		rc := utf8.RuneCountInString(line)
		if end.Col >= rc {
			after = true
		}
	}
	ed.OpDeleteRegion(start, end, inclusive)
	if after {
		b.Loc.Col++
	}
	// XXX replay?
	ed.inp.Init(b.CurrentLine(), b.Loc.Col, ed.cfg.AutoIndent)
	ed.inpRow = b.Loc.Row
	ed.mode = ModeInsert
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
	b := ed.Buf()
	start := b.Loc
	end, ok := ed.MoveByWordEx(n)
	if !ok {
		return
	}
	ed.OpChangeRegion(start, end, false, replay)
	// TODO n
}

// C : Change to end of current line.
func (ed *Editor) OpChangeToEnd(n int, replay bool) {
	if n < 1 {
		ed.Error("OpChangeToEnd: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buf()
	rs := []rune(ed.CurrentLine())
	if b.Loc.Col < len(rs) {
		ed.regs.SetRunes("", []string{string(rs[b.Loc.Col:])})
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
	b := ed.Buf()
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
