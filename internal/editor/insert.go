package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/rkind"
	"tea.kareha.org/cup/levi/internal/rutil"
)

////////////////////////////////
// Commands Using Insert Mode //
////////////////////////////////

// Internal Use //

// Repeats insertion.
func (ed *Editor) replayInsert() {
	inserted := append([]string{}, ed.inserted...)
	b := ed.Buf()
	line := b.CurrentLine()
	head, tail := rutil.Split(line, b.Loc.Col)
	inserted[0] = head + inserted[0]
	inserted[len(inserted)-1] = inserted[len(inserted)-1] + tail
	if ed.cfg.AutoIndent {
		for i := 0; i < len(inserted); i++ {
			if rkind.IsBlankLine(inserted[i]) {
				inserted[i] = ""
			}
		}
	}
	lines := append([]string{}, b.Lines[:b.Loc.Row]...)
	lines = append(lines, inserted...)
	if b.Loc.Row+1 < b.NumLines() {
		lines = append(lines, b.Lines[b.Loc.Row+1:]...)
	}
	b.Lines = lines
	if len(inserted) < 2 {
		b.Loc.Col += utf8.RuneCountInString(ed.inserted[0])
	} else {
		b.Loc.Row += len(ed.inserted) - 1
		b.Loc.Col = utf8.RuneCountInString(ed.inserted[len(ed.inserted)-1])
	}
}

//
// Insert
//

// Switches to insert mode.
// Key: i
func (ed *Editor) Insert(n int, replay bool) bool {
	if n < 1 {
		ed.Error("Insert: n < 1")
		return false
	}
	if replay {
		if len(ed.inserted) < 0 {
			ed.Notice("Not inserted yet")
			return false
		}
		for i := 0; i < n; i++ {
			ed.replayInsert()
		}
		b := ed.Buf()
		b.Loc.Col--
		b.Loc = b.ConfineInclusive(b.Loc)
		b.VirtCol = b.Loc.Col
		return true
	}
	b := ed.Buf()
	ed.inp.Init(b.CurrentLine(), b.Loc.Col, ed.cfg.AutoIndent)
	ed.inpRow = b.Loc.Row
	ed.mode = ModeInsert
	return false
}

// Switches to insert mode after cursor.
// Key: a
func (ed *Editor) InsertAfter(n int, replay bool) bool {
	if n < 1 {
		ed.Error("InsertAfter: n < 1")
		return false
	}
	if replay {
		if len(ed.inserted) < 0 {
			ed.Notice("Not inserted yet")
			return false
		}
		b := ed.Buf()
		rc := utf8.RuneCountInString(b.CurrentLine())
		if b.Loc.Col >= rc-1 {
			b.Loc.Col = rc
		} else {
			b.Loc.Col++
		}
		for i := 0; i < n; i++ {
			ed.replayInsert()
		}
		b.Loc.Col--
		b.Loc = b.ConfineInclusive(b.Loc)
		b.VirtCol = b.Loc.Col
		return true
	}
	b := ed.Buf()
	rc := utf8.RuneCountInString(b.CurrentLine())
	if b.Loc.Col >= rc-1 {
		b.Loc.Col = rc
	} else {
		b.Loc.Col++
	}
	ed.inp.Init(b.CurrentLine(), b.Loc.Col, ed.cfg.AutoIndent)
	ed.inpRow = b.Loc.Row
	ed.mode = ModeInsert
	return false
}

// Switches to insert mode after indent of current line.
// Key: I
func (ed *Editor) InsertAfterIndent(n int, replay bool) bool {
	if n < 1 {
		ed.Error("InsertAfterIndent: n < 1")
		return false
	}
	b := ed.Buf()
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
	return ed.Insert(n, replay)
}

// Switches to insert mode after end of current line.
// Key: A
func (ed *Editor) InsertAfterEnd(n int, replay bool) bool {
	if n < 1 {
		ed.Error("InsertAfterEnd: n < 1")
		return false
	}
	b := ed.Buf()
	rc := utf8.RuneCountInString(b.CurrentLine())
	b.Loc.Col = max(rc-1, 0)
	return ed.InsertAfter(n, replay)
}

//
// Insert Line
//

// Inserts blank line below cursor and switches to insert mode.
// Key: o
func (ed *Editor) InsertLine(n int, replay bool) bool {
	if n < 1 {
		ed.Error("InsertLine: n < 1")
		return false
	}
	b := ed.Buf()
	if b.NumLines() < 1 {
		return ed.InsertAfter(n, replay)
	}
	indent := ""
	if ed.cfg.AutoIndent {
		indent = rkind.IndentOf(b.CurrentLine())
	}
	lines := []string{}
	if b.NumLines() > 0 {
		lines = append(lines, b.Lines[:b.Loc.Row+1]...)
	}
	lines = append(lines, indent)
	if b.Loc.Row+1 <= b.NumLines()-1 {
		lines = append(lines, b.Lines[b.Loc.Row+1:]...)
	}
	b.Lines = lines
	b.Loc.Row++
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
	return ed.InsertAfter(n, replay)
}

// Inserts blank line above cursor and switches to insert mode.
// Key: O
func (ed *Editor) InsertLineAbove(n int, replay bool) bool {
	if n < 1 {
		ed.Error("InsertLineAbove: n < 1")
		return false
	}
	indent := ""
	b := ed.Buf()
	if ed.cfg.AutoIndent {
		indent = rkind.IndentOf(b.CurrentLine())
	}
	lines := []string{}
	if b.Loc.Row > 0 {
		lines = append(lines, b.Lines[:b.Loc.Row]...)
	}
	lines = append(lines, indent)
	if b.Loc.Row <= b.NumLines()-1 {
		lines = append(lines, b.Lines[b.Loc.Row:]...)
	}
	b.Lines = lines
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
	return ed.InsertAfter(n, replay)
}

//
// Change / Substitute
//

// Changes region from cursor to destination of motion.
// Key: c<mv> (when <mv> is characterwise motion)
func (ed *Editor) ChangeRegion(
	reg rune, start buf.Loc, end buf.Loc, inclusive bool, replay bool,
) bool {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, inclusive, false)
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

// Changes line region from cursor to destination of motion.
// Key: c<mv> (when <mv> is linewise motion)
func (ed *Editor) ChangeLineRegion(
	reg rune, start buf.Loc, end buf.Loc, replay bool,
) bool {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, true, true)
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

// Substitutes character under cursor.
// Key: s<char>
func (ed *Editor) Subst(reg rune, n int, replay bool) bool {
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

//
// Unsupported Commands
//

// Switches to insert mode and overwrites current line.
// Key: R
func (ed *Editor) Overwrite() bool {
	ed.Unsupported("Overwrite")
	return false
}
