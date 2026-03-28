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
	ed.mode = ModeInsert
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
}

// I : Switch to insert mode before first non-blank character of current line.
func (ed *Editor) InsertBeforeNonBlank(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("InsertBeforeNonBlank")
}

// A : Switch to insert mode after end of current line.
func (ed *Editor) InsertAfterEnd(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("InsertAfterEnd")
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
	ed.Unimplemented("InsertOpenBelow")
}

// O : Open a new line above and switch to insert mode.
func (ed *Editor) InsertOpenAbove(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("InsertOpenAbove")
}
