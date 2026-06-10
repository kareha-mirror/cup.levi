package editor

////////////////////////
// Insertion Commands //
////////////////////////

//
// Enter Insert Mode
//

// i : Switch to insert mode before cursor.
func (ed *Editor) InsertBefore(n int) {
	ed.EnsureCommand()
	ed.inp.Init(ed.CurrentLine(), ed.col)
	ed.inpRow = ed.row
	ed.mode = ModeInsert
	// XXX n
}

// a : Switch to insert mode after cursor.
func (ed *Editor) InsertAfter(n int) {
	ed.EnsureCommand()
	rc := ed.RuneCount()
	if ed.col >= rc-1 {
		ed.col = rc
	} else {
		ed.MoveRight(1)
	}
	ed.InsertBefore(n)
	// XXX n
}

// I : Switch to insert mode before first non-blank character of current line.
func (ed *Editor) InsertBeforeNonBlank(n int) {
	ed.MoveToNonBlank()
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
	lines := []string{}
	if len(ed.lines) > 0 {
		lines = append(lines, ed.lines[:ed.row+1]...)
	}
	lines = append(lines, "")
	if ed.row+1 <= len(ed.lines)-1 {
		lines = append(lines, ed.lines[ed.row+1:]...)
	}
	ed.lines = lines
	ed.row++
	ed.Confine()
	ed.InsertAfter(n)
	// XXX n
}

// O : Open a new line above and switch to insert mode.
func (ed *Editor) InsertOpenAbove(n int) {
	ed.EnsureCommand()
	lines := []string{}
	if ed.row > 0 {
		lines = append(lines, ed.lines[:ed.row]...)
	}
	lines = append(lines, "")
	if ed.row <= len(ed.lines)-1 {
		lines = append(lines, ed.lines[ed.row:]...)
	}
	ed.lines = lines
	ed.Confine()
	ed.InsertAfter(n)
	// XXX n
}
