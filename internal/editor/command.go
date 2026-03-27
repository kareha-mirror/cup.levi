package editor

// key: i
func (ed *Editor) InsertBefore() {
	if ed.mode == ModeInsert {
		panic("invalid state")
	}
	ed.ins.Init(ed.CurrentLine(), ed.col)
	ed.mode = ModeInsert
}

// key: a
func (ed *Editor) InsertAfter() {
	if ed.mode == ModeInsert {
		panic("invalid state")
	}
	rc := ed.RuneCount()
	if ed.col >= rc-1 {
		ed.col = rc
	} else {
		ed.MoveRight(1)
	}
	ed.InsertBefore()
}

// key: x
func (ed *Editor) OpDelete(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	if len(ed.CurrentLine()) < 1 {
		ed.Ring()
		return
	}
	rs := []rune(ed.CurrentLine())
	if ed.col < 1 {
		ed.lines[ed.row] = string(rs[1:])
	} else {
		head := string(rs[:ed.col])
		tail := string(rs[ed.col+1:])
		ed.lines[ed.row] = head + tail
	}
	ed.Confine()
}

// key: ZZ
func (ed *Editor) MiscSaveAndQuit() {
	ed.quit = true
}
