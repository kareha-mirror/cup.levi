package editor

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
func (ed *Editor) InsertBefore(n int) {
	ed.EnsureCommand()
	b := ed.Buffer()
	ed.inp.Init(ed.CurrentLine(), b.col)
	ed.inpRow = b.row
	ed.mode = ModeInsert
	// XXX n
}

// a : Switch to insert mode after cursor.
func (ed *Editor) InsertAfter(n int) {
	ed.EnsureCommand()
	b := ed.Buffer()
	rc := ed.RuneCount()
	if b.col >= rc-1 {
		b.col = rc
	} else {
		ed.MoveRight(1)
	}
	ed.InsertBefore(n)
	// XXX n
}

// I : Switch to insert mode before first non-blank character of current line.
func (ed *Editor) InsertBeforeNonBlank(n int) {
	ed.toNonBlankCol()
	ed.InsertBefore(n)
	// XXX n
}

// A : Switch to insert mode after end of current line.
func (ed *Editor) InsertAfterEnd(n int) {
	ed.MoveToEnd()
	ed.InsertAfter(n)
	// XXX n
}

// R : Switch to replace (overwrite) mode.
func (ed *Editor) InsertOverwrite(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("InsertOverwrite")
}

//
// Open Line
//

// o : Open a new line below and switch to insert mode.
func (ed *Editor) InsertOpenBelow(n int) {
	ed.EnsureCommand()
	b := ed.Buffer()
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
	ed.InsertAfter(n)
	// XXX n
}

// O : Open a new line above and switch to insert mode.
func (ed *Editor) InsertOpenAbove(n int) {
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
	ed.InsertAfter(n)
	// XXX n
}
