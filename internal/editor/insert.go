package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/rkind"
)

////////////////////////
// Insertion Commands //
////////////////////////

//
// Enter Insert Mode
//

func (ed *Editor) replayInsert() {
	inserted := append([]string{}, ed.inserted...)
	b := ed.Buf()
	rs := []rune(b.CurrentLine())
	head := string(rs[:b.Loc.Col])
	tail := ""
	if b.Loc.Col < len(rs) {
		tail = string(rs[b.Loc.Col:])
	}
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

// i : Switch to insert mode before cursor.
func (ed *Editor) InsertBefore(n int, replay bool) {
	if n < 1 {
		ed.Error("InsertBefore: n < 1")
		return
	}
	if replay {
		if len(ed.inserted) < 0 {
			ed.Notice("Not inserted yet")
			return
		}
		for i := 0; i < n; i++ {
			ed.replayInsert()
		}
		b := ed.Buf()
		b.Loc.Col--
		b.Loc = b.ConfineInclusive(b.Loc)
		b.VirtCol = b.Loc.Col
		b.Modified = true
	} else {
		b := ed.Buf()
		ed.inp.Init(b.CurrentLine(), b.Loc.Col, ed.cfg.AutoIndent)
		ed.inpRow = b.Loc.Row
		ed.mode = ModeInsert
	}
}

// a : Switch to insert mode after cursor.
func (ed *Editor) InsertAfter(n int, replay bool) {
	if n < 1 {
		ed.Error("InsertAfter: n < 1")
		return
	}
	if replay {
		if len(ed.inserted) < 0 {
			ed.Notice("Not inserted yet")
			return
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
		b.Modified = true
	} else {
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
	}
}

// I : Switch to insert mode before first non-blank character of current line.
func (ed *Editor) InsertBeforeNonBlank(n int, replay bool) {
	if n < 1 {
		ed.Error("InsertBeforeNonBlank: n < 1")
		return
	}
	b := ed.Buf()
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
	ed.InsertBefore(n, replay)
}

// A : Switch to insert mode after end of current line.
func (ed *Editor) InsertAfterEnd(n int, replay bool) {
	if n < 1 {
		ed.Error("InsertAfterEnd: n < 1")
		return
	}
	b := ed.Buf()
	rc := utf8.RuneCountInString(b.CurrentLine())
	b.Loc.Col = max(rc-1, 0)
	ed.InsertAfter(n, replay)
}

// R : Switch to replace (overwrite) mode.
func (ed *Editor) InsertOverwrite(n int, replay bool) {
	if n < 1 {
		ed.Error("InsertOverwrite: n < 1")
		return
	}
	ed.Unimplemented("InsertOverwrite")
}

//
// Open Line
//

// o : Open a new line below and switch to insert mode.
func (ed *Editor) InsertOpenBelow(n int, replay bool) {
	if n < 1 {
		ed.Error("InsertOpenBelow: n < 1")
		return
	}
	b := ed.Buf()
	if b.NumLines() < 1 {
		ed.InsertAfter(n, replay)
		return
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
	ed.InsertAfter(n, replay)
}

// O : Open a new line above and switch to insert mode.
func (ed *Editor) InsertOpenAbove(n int, replay bool) {
	if n < 1 {
		ed.Error("InsertOpenAbove: n < 1")
		return
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
	ed.InsertAfter(n, replay)
}
