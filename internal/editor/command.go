package editor

func (ed *Editor) enterInsert() {
	rs := []rune(ed.lines[ed.row])
	ed.head = string(rs[:ed.col])
	ed.tail = string(rs[ed.col:])
	ed.mode = modeInsert
}

func (ed *Editor) enterInsertAfter() {
	len := ed.runeCount()
	if ed.col >= len-1 {
		ed.col = len
		ed.head = ed.lines[ed.row]
		ed.tail = ""
		ed.mode = modeInsert
		return
	}

	ed.moveRight(1)
	ed.enterInsert()
}

func (ed *Editor) moveLeft(n int) {
	ed.col = max(ed.col-n, 0)
}

func (ed *Editor) moveRight(n int) {
	ed.col = min(ed.col+n, max(ed.runeCount()-1, 0))
}

func (ed *Editor) moveDown(n int) {
	ed.row = min(ed.row+n, max(len(ed.lines)-1, 0))
}

func (ed *Editor) moveUp(n int) {
	ed.row = max(ed.row-n, 0)
}
