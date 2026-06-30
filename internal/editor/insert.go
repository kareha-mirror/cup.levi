package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/rkind"
	"tea.kareha.org/cup/levi/internal/rutil"
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

// i : Switch to insert mode before cursor.
func (ed *Editor) InsertBefore(n int, replay bool) bool {
	if n < 1 {
		ed.Error("InsertBefore: n < 1")
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

// a : Switch to insert mode after cursor.
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

// I : Switch to insert mode before first non-blank character of current line.
func (ed *Editor) InsertBeforeNonBlank(n int, replay bool) bool {
	if n < 1 {
		ed.Error("InsertBeforeNonBlank: n < 1")
		return false
	}
	b := ed.Buf()
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
	return ed.InsertBefore(n, replay)
}

// A : Switch to insert mode after end of current line.
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

// R : Switch to replace (overwrite) mode.
func (ed *Editor) Overwrite(n int, replay bool) bool {
	if n < 1 {
		ed.Error("Overwrite: n < 1")
		return false
	}
	ed.Unimplemented("Overwrite")
	return false
}

//
// Open Line
//

// o : Open a new line below and switch to insert mode.
func (ed *Editor) OpenBelow(n int, replay bool) bool {
	if n < 1 {
		ed.Error("OpenBelow: n < 1")
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

// O : Open a new line above and switch to insert mode.
func (ed *Editor) OpenAbove(n int, replay bool) bool {
	if n < 1 {
		ed.Error("OpenAbove: n < 1")
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
