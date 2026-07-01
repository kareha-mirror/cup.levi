package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/rkind"
)

////////////////////////////////
// Operator Commands (Change) //
////////////////////////////////

//
// Change / Substitute
//

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
	// empty case
	if b.NumLines() == 0 && n == 1 {
		return ed.Insert(1, replay)
	}
	b.Loc = start
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
