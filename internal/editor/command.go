package editor

// key: i
func (ed *Editor) Insert() {
	if ed.mode == ModeInsert {
		panic("invalid state")
	}
	ed.ins.Enter(ed.lines[ed.row], ed.col)
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
	ed.Insert()
}

// key: h
func (ed *Editor) MoveLeft(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.col = max(ed.col-n, 0)
}

// key: l
func (ed *Editor) MoveRight(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.col = min(ed.col+n, max(ed.RuneCount()-1, 0))
}

// key: j
func (ed *Editor) MoveDown(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.row = min(ed.row+n, max(len(ed.lines)-1, 0))
}

// key: k
func (ed *Editor) MoveUp(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.row = max(ed.row-n, 0)
}

// key: x
func (ed *Editor) DeleteRune(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	if len(ed.lines[ed.row]) < 1 {
		ed.Ring()
		return
	}
	rs := []rune(ed.lines[ed.row])
	if ed.col < 1 {
		ed.lines[ed.row] = string(rs[1:])
	} else {
		head := string(rs[:ed.col])
		tail := string(rs[ed.col+1:])
		ed.lines[ed.row] = head + tail
	}
	rc := ed.RuneCount()
	if ed.col >= rc {
		ed.col = max(rc-1, 0)
	}
}
