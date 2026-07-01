package editor

func (ed *Editor) NextBuf() {
	if ed.bufIdx+1 >= len(ed.bufs) {
		ed.Ring("No more files to edit.")
		return
	}
	if !ed.bufMove {
		ed.lastBufIdx = ed.bufIdx
	}
	ed.bufIdx++
	ed.undo = false
	ed.redraw = true
	ed.ShowFileInfo()
	ed.bufMove = true
}

func (ed *Editor) PrevBuf() {
	if ed.bufIdx-1 < 0 {
		ed.Ring("No previous files to edit.")
		return
	}
	if !ed.bufMove {
		ed.lastBufIdx = ed.bufIdx
	}
	ed.bufIdx--
	ed.undo = false
	ed.redraw = true
	ed.ShowFileInfo()
	ed.bufMove = true
}

func (ed *Editor) LastBuf() {
	if ed.lastBufIdx == ed.bufIdx {
		ed.Ring("No need to change buffer.")
		return
	}
	ed.bufIdx, ed.lastBufIdx = ed.lastBufIdx, ed.bufIdx
	ed.undo = false
	ed.redraw = true
	ed.ShowFileInfo()
}

func (ed *Editor) GoToBuf(n int) bool {
	if n < 1 {
		ed.Error("GoToBuf: n < 1")
		return false
	}
	n--
	if n >= ed.NumBufs() {
		ed.Ring("Out of range")
		return false
	}
	if !ed.bufMove {
		ed.lastBufIdx = ed.bufIdx
	}
	ed.bufIdx = n
	ed.undo = false
	ed.redraw = true
	ed.ShowFileInfo()
	ed.bufMove = true
	return true
}
