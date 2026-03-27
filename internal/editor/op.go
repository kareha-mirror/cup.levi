package editor

// x : Delete character under cursor.
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
