package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/rkind"
	"tea.kareha.org/cup/levi/internal/rutil"
)

////////////////////////////////
// Operator Commands (Change) //
////////////////////////////////

//
// Change / Substitute
//

// cc : Change current line.
func (ed *Editor) ChangeLine(reg string, n int, replay bool) bool {
	if n < 1 {
		ed.Error("ChangeLine: n < 1")
		return false
	}
	b := ed.Buf()
	// empty case
	if b.NumLines() == 0 && n == 1 {
		return ed.Insert(1, replay)
	}
	if b.Loc.Row+n > b.NumLines() {
		ed.Notice("Out of range")
		return false
	}
	ed.ApplyRegLines(reg, b.Lines[b.Loc.Row:b.Loc.Row+n])

	// above
	lines := append([]string{}, b.Lines[:b.Loc.Row]...)
	// current
	if ed.cfg.AutoIndent {
		lines = append(lines, rkind.IndentOf(b.Lines[b.Loc.Row]))
	} else {
		lines = append(lines, "")
	}
	// below
	if b.Loc.Row+n < b.NumLines() {
		lines = append(lines, b.Lines[b.Loc.Row+n:]...)
	}
	b.Lines = lines

	return ed.InsertAfterEnd(1, replay)
}

// c<mv> : Change region from current cursor to destination of motion <mv>.
func (ed *Editor) ChangeRegion(
	reg string, start buf.Loc, end buf.Loc, inclusive bool, replay bool,
) bool {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, inclusive)
	after := false
	if end.Row < b.NumLines() {
		rc := utf8.RuneCountInString(b.Line(end.Row))
		if end.Col >= rc {
			after = true
		}
	}
	ed.DeleteRegion(reg, start, end, inclusive)
	if after {
		b.Loc.Col++
	}
	if replay {
		if len(ed.inserted) < 0 {
			ed.Notice("Not inserted yet")
			return false
		}
		ed.replayInsert()
		b.Loc.Col--
		b.Loc = b.ConfineInclusive(b.Loc)
		b.VirtCol = b.Loc.Col
		return true
	}
	ed.inp.Init(b.CurrentLine(), b.Loc.Col, ed.cfg.AutoIndent)
	ed.inpRow = b.Loc.Row
	ed.mode = ModeInsert
	return false
}

// c<mv> : Change region from current cursor to destination of motion <mv>.
func (ed *Editor) ChangeLineRegion(
	reg string, start buf.Loc, end buf.Loc, replay bool,
) bool {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, true)
	n := end.Row - start.Row + 1
	b.Loc = start
	return ed.ChangeLine(reg, n, replay)
}

// cw : Change word.
func (ed *Editor) ChangeWord(reg string, n int, replay bool) bool {
	if n < 1 {
		ed.Error("ChangeWord: n < 1")
		return false
	}
	start := ed.Buf().Loc
	end, ok := ed.MoveByWordAlt(n)
	if !ok {
		ed.Error("Failed to move")
		return false
	}
	return ed.ChangeRegion(reg, start, end, false, replay)
}

// C : Change to end of current line.
func (ed *Editor) ChangeToEnd(reg string, n int, replay bool) bool {
	if n < 1 {
		ed.Error("ChangeToEnd: n < 1")
		return false
	}
	b := ed.Buf()
	// empty case
	if b.NumLines() == 0 && n == 1 {
		return ed.Insert(1, replay)
	}
	if b.Loc.Row+n > b.NumLines() {
		ed.Notice("Out of range")
		return false
	}
	head, tail := rutil.Split(b.CurrentLine(), b.Loc.Col)
	killed := []string{tail}
	killed = append(killed, b.Lines[b.Loc.Row+1:b.Loc.Row+n]...)
	ed.ApplyRegRunes(reg, killed)

	// above
	lines := append([]string{}, b.Lines[:b.Loc.Row]...)
	// current
	if rkind.IsBlankLine(head) {
		if ed.cfg.AutoIndent {
			lines = append(lines, rkind.IndentOf(b.Lines[b.Loc.Row]))
		} else {
			lines = append(lines, "")
		}
	} else {
		lines = append(lines, head)
	}
	// below
	if b.Loc.Row+n < b.NumLines() {
		lines = append(lines, b.Lines[b.Loc.Row+n:]...)
	}
	b.Lines = lines

	return ed.InsertAfterEnd(1, replay)
}

// s : Substitute one character under cursor.
func (ed *Editor) Subst(reg string, n int, replay bool) bool {
	if n < 1 {
		ed.Error("Subst: n < 1")
		return false
	}
	ed.internalDelete(reg, n)
	b := ed.Buf()
	if replay {
		if len(ed.inserted) < 0 {
			ed.Notice("Not inserted yet")
			return false
		}
		ed.replayInsert()
		b.Loc.Col--
		b.Loc = b.ConfineInclusive(b.Loc)
		b.VirtCol = b.Loc.Col
		return true
	}
	ed.inp.Init(b.CurrentLine(), b.Loc.Col, ed.cfg.AutoIndent)
	ed.inpRow = b.Loc.Row
	ed.mode = ModeInsert
	return false
}

// S : Substtute current line (equals cc).
func (ed *Editor) SubstLine(reg string, n int, replay bool) bool {
	if n < 1 {
		ed.Error("SubstLine: n < 1")
		return false
	}
	return ed.ChangeLine(reg, n, replay)
}
