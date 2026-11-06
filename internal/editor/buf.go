package editor

///////////////////////////////////////////
// Commands for Selecting Current Buffer //
///////////////////////////////////////////

// Go to last visited buffer.
// Key: Ctrl-^, Ctrl-_
func (ed *Editor) LastBuf() {
	if ed.lastBufIdx == ed.bufIdx {
		//ed.Ring("No previous files to edit.") // nvi
		ed.Ring("No last files to edit.")
		return
	}
	ed.bufIdx, ed.lastBufIdx = ed.lastBufIdx, ed.bufIdx
	ed.undo = false
	ed.ShowFileInfo()
}

// Go to buffer specified by <num>.
// Key: <num> Ctrl-^, <num> Ctrl-_
// Not available in nvi.
// Available in Vim.
func (ed *Editor) GoToBuf(n int) bool { // n is 1-based
	if n < 1 {
		ed.Error("GoToBuf: n < 1")
		return false
	}
	n-- // to 0-based
	if n >= ed.NumBufs() {
		ed.Ring("Out of range") // levi
		return false
	}
	if !ed.bufMove {
		ed.lastBufIdx = ed.bufIdx
		ed.bufMove = true
	}
	ed.bufIdx = n
	ed.undo = false
	ed.ShowFileInfo()
	return true
}

// Go to next buffer in list.
// Key: :next Enter, :n Enter, zj (levi enhancement)
func (ed *Editor) NextBuf() {
	if ed.bufIdx+1 >= len(ed.bufs) {
		ed.Ring("No more files to edit.")
		return
	}
	if !ed.bufMove {
		ed.lastBufIdx = ed.bufIdx
		ed.bufMove = true
	}
	ed.bufIdx++
	ed.undo = false
	ed.ShowFileInfo()
}

// Go to previous buffer in list.
// Key: :prev Enter, zk (levi enhancement)
func (ed *Editor) PrevBuf() {
	if ed.bufIdx-1 < 0 {
		ed.Ring("No previous files to edit.")
		return
	}
	if !ed.bufMove {
		ed.lastBufIdx = ed.bufIdx
		ed.bufMove = true
	}
	ed.bufIdx--
	ed.undo = false
	ed.ShowFileInfo()
}
