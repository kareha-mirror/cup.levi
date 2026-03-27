package editor

// key: i
func (ed *Editor) Insert() {
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
	ed.Insert()
}

// key: h
func (ed *Editor) MoveLeft(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.col -= n
	ed.Confine()
}

// key: l
func (ed *Editor) MoveRight(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.col += n
	ed.Confine()
}

// key: j
func (ed *Editor) MoveDown(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.row += n
	ed.Confine()
}

// key: k
func (ed *Editor) MoveUp(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.row -= n
	ed.Confine()
}

// key: x
func (ed *Editor) DeleteRune(n int) {
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
