package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/rutil"
)

//
// Change / Substitute
//

// cc : Change current line.
func (ed *Editor) ChangeLine(reg string, n int, replay bool) bool {
	if n < 1 {
		ed.Error("ChangeLine: n < 1")
		return false
	}
	ed.Unimplemented("ChangeLine")
	return false
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
	ed.Unimplemented("ChangeLineRegion")
	return false
}

// cw : Change word.
func (ed *Editor) ChangeWord(reg string, n int, replay bool) bool {
	if n < 1 {
		ed.Error("ChangeWord: n < 1")
		return false
	}
	start := ed.Buf().Loc
	end, ok := ed.MoveByChangeWord(n)
	if !ok {
		ed.Error("Failed to move")
		return false
	}
	return ed.ChangeRegion(reg, start, end, false, replay)
	// TODO n
}

// C : Change to end of current line.
func (ed *Editor) ChangeToEnd(reg string, n int, replay bool) bool {
	if n < 1 {
		ed.Error("ChangeToEnd: n < 1")
		return false
	}
	b := ed.Buf()
	head, tail := rutil.Split(b.CurrentLine(), b.Loc.Col)
	ed.ApplyRegRunes(reg, []string{tail})
	b.SetCurrentLine(head)
	ed.inp.Init(head, b.Loc.Col, ed.cfg.AutoIndent)
	ed.inpRow = b.Loc.Row
	ed.mode = ModeInsert
	return false
	// TODO n
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
