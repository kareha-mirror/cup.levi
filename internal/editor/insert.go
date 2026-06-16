package editor

import (
	"unicode/utf8"
)

func getIndent(s string) string {
	runes := []rune{}
	for _, r := range s {
		if !isBlank(r) {
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

// i : Switch to insert mode before cursor.
func (ed *Editor) InsertBefore(n int, replay bool) {
	ed.EnsureCommand()
	b := ed.Buffer()
	if replay {
		if len(ed.inserted) < 0 {
			return
		}
		inserted := append([]string{}, ed.inserted...)
		rs := []rune(ed.CurrentLine())
		head := string(rs[:b.col])
		tail := ""
		if b.col < len(rs) {
			tail = string(rs[b.col:])
		}
		inserted[0] = head + inserted[0]
		inserted[len(inserted)-1] = inserted[len(inserted)-1] + tail
		if ed.cfg.AutoIndent {
			for i := 0; i < len(inserted); i++ {
				if isBlankLine(inserted[i]) {
					inserted[i] = ""
				}
			}
		}
		lines := append([]string{}, b.lines[:b.row]...)
		lines = append(lines, inserted...)
		if b.row+1 <= len(b.lines)-1 {
			lines = append(lines, b.lines[b.row+1:]...)
		}
		b.lines = lines
		if len(inserted) < 2 {
			b.col += utf8.RuneCountInString(inserted[0])
		} else {
			b.row += len(inserted) - 1
			b.col = utf8.RuneCountInString(inserted[len(inserted)-1])
		}
		ed.MoveLeft(1)
		b.modified = true
	} else {
		ed.inp.Init(ed.CurrentLine(), b.col)
		ed.inpRow = b.row
		ed.mode = ModeInsert
	}
	// XXX n
}

// a : Switch to insert mode after cursor.
func (ed *Editor) InsertAfter(n int, replay bool) {
	ed.EnsureCommand()
	b := ed.Buffer()
	rc := ed.RuneCount()
	if b.col >= rc-1 {
		b.col = rc
	} else {
		ed.MoveRight(1)
	}
	ed.InsertBefore(n, replay)
	// XXX n
}

// I : Switch to insert mode before first non-blank character of current line.
func (ed *Editor) InsertBeforeNonBlank(n int, replay bool) {
	ed.toNonBlankCol()
	ed.InsertBefore(n, replay)
	// XXX n
}

// A : Switch to insert mode after end of current line.
func (ed *Editor) InsertAfterEnd(n int, replay bool) {
	ed.MoveToEnd()
	ed.InsertAfter(n, replay)
	// XXX n
}

// R : Switch to replace (overwrite) mode.
func (ed *Editor) InsertOverwrite(n int, replay bool) {
	ed.EnsureCommand()
	ed.Unimplemented("InsertOverwrite")
}

//
// Open Line
//

// o : Open a new line below and switch to insert mode.
func (ed *Editor) InsertOpenBelow(n int, replay bool) {
	ed.EnsureCommand()
	b := ed.Buffer()
	if len(b.lines) < 1 {
		ed.InsertAfter(n, replay)
		return
	}
	indent := ""
	if ed.cfg.AutoIndent {
		indent = getIndent(ed.CurrentLine())
	}
	lines := []string{}
	if len(b.lines) > 0 {
		lines = append(lines, b.lines[:b.row+1]...)
	}
	lines = append(lines, indent)
	if b.row+1 <= len(b.lines)-1 {
		lines = append(lines, b.lines[b.row+1:]...)
	}
	b.lines = lines
	b.row++
	ed.toNonBlankCol()
	ed.InsertAfter(n, replay)
	// XXX n
}

// O : Open a new line above and switch to insert mode.
func (ed *Editor) InsertOpenAbove(n int, replay bool) {
	ed.EnsureCommand()
	b := ed.Buffer()
	indent := ""
	if ed.cfg.AutoIndent {
		indent = getIndent(ed.CurrentLine())
	}
	lines := []string{}
	if b.row > 0 {
		lines = append(lines, b.lines[:b.row]...)
	}
	lines = append(lines, indent)
	if b.row <= len(b.lines)-1 {
		lines = append(lines, b.lines[b.row:]...)
	}
	b.lines = lines
	ed.toNonBlankCol()
	ed.InsertAfter(n, replay)
	// XXX n
}
