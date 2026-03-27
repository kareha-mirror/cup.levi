package editor

// i : Switch to insert mode before cursor.
func (ed *Editor) InsertBefore(n int) {
	if ed.mode == ModeInsert {
		panic("invalid state")
	}
	ed.inp.Init(ed.CurrentLine(), ed.col)
	ed.mode = ModeInsert
}

// a : Switch to insert mode after cursor.
func (ed *Editor) InsertAfter(n int) {
	if ed.mode == ModeInsert {
		panic("invalid state")
	}
	rc := ed.RuneCount()
	if ed.col >= rc-1 {
		ed.col = rc
	} else {
		ed.MoveRight(1)
	}
	ed.InsertBefore(n)
}
