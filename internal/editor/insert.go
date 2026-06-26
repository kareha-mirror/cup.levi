package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/rkind"
)

func getIndent(s string) string {
	runes := []rune{}
	for _, r := range s {
		if !rkind.IsBlank(r) {
			break
		}
		runes = append(runes, r)
	}
	return string(runes)
}

////////////////////////
// Insertion Commands //
////////////////////////

//
// Enter Insert Mode
//

func (ed *Editor) replayInsert(line string) {
	b := ed.Buf()
	inserted := append([]string{}, ed.inserted...)
	rs := []rune(line)
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
	if b.Loc.Row+1 <= b.NumLines()-1 {
		lines = append(lines, b.Lines[b.Loc.Row+1:]...)
	}
	b.Lines = lines
	if len(inserted) < 2 {
		b.Loc.Col += utf8.RuneCountInString(ed.inserted[0])
	} else {
		b.Loc.Row += len(ed.inserted) - 1
		b.Loc.Col = utf8.RuneCountInString(ed.inserted[len(ed.inserted)-1])
	}
	ed.MoveLeft(1)
	b.Modified = true
}

// i : Switch to insert mode before cursor.
func (ed *Editor) InsertBefore(n int, replay bool) {
	ed.Commit()
	b := ed.Buf()
	if replay {
		if len(ed.inserted) < 0 {
			return
		}
		for i := 0; i < n; i++ {
			ed.replayInsert(b.CurrentLine())
		}
	} else {
		ed.inp.Init(b.CurrentLine(), b.Loc.Col, ed.cfg.AutoIndent)
		ed.inpRow = b.Loc.Row
		ed.mode = ModeInsert
	}
}

// a : Switch to insert mode after cursor.
func (ed *Editor) InsertAfter(n int, replay bool) {
	ed.Commit()
	if !replay {
		n = 1
	}
	b := ed.Buf()
	for i := 0; i < n; i++ {
		rc := utf8.RuneCountInString(b.CurrentLine())
		if b.Loc.Col >= rc-1 {
			b.Loc.Col = rc
		} else {
			ed.MoveRight(1)
		}
		ed.InsertBefore(1, replay)
	}
}

// I : Switch to insert mode before first non-blank character of current line.
func (ed *Editor) InsertBeforeNonBlank(n int, replay bool) {
	if !replay {
		n = 1
	}
	b := ed.Buf()
	for i := 0; i < n; i++ {
		b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
		ed.InsertBefore(1, replay)
	}
}

// A : Switch to insert mode after end of current line.
func (ed *Editor) InsertAfterEnd(n int, replay bool) {
	if !replay {
		n = 1
	}
	b := ed.Buf()
	for i := 0; i < n; i++ {
		loc, ok := ed.MoveToEnd()
		if !ok {
			return
		}
		b.Loc = b.ConfineInclusive(loc)
		ed.InsertAfter(1, replay)
	}
}

// R : Switch to replace (overwrite) mode.
func (ed *Editor) InsertOverwrite(n int, replay bool) {
	ed.Commit()
	ed.Unimplemented("InsertOverwrite")
}

//
// Open Line
//

// o : Open a new line below and switch to insert mode.
func (ed *Editor) InsertOpenBelow(n int, replay bool) {
	ed.Commit()
	if !replay {
		n = 1
	}
	b := ed.Buf()
	for i := 0; i < n; i++ {
		if b.NumLines() < 1 {
			ed.InsertAfter(n, replay)
			return
		}
		indent := ""
		if ed.cfg.AutoIndent {
			indent = getIndent(b.CurrentLine())
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
		ed.InsertAfter(1, replay)
	}
}

// O : Open a new line above and switch to insert mode.
func (ed *Editor) InsertOpenAbove(n int, replay bool) {
	ed.Commit()
	if !replay {
		n = 1
	}
	b := ed.Buf()
	for i := 0; i < n; i++ {
		indent := ""
		if ed.cfg.AutoIndent {
			indent = getIndent(b.CurrentLine())
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
		ed.InsertAfter(1, replay)
	}
}
