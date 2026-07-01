package editor

///////////////////////////////////////////
// Commands for Selecting Current Buffer //
///////////////////////////////////////////

// :next, :n : Go to next buffer in list.
// Also available as 'zj' (levi enhancement).
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

// :prev : Go to previous buffer in list.
// Also available as 'zk' (levi enhancement).
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

// Ctrl-^, Ctrl-_ : Go to last visited buffer.
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

// <num> Ctrl-^, <num> Ctrl-_ : Go to buffer specified by <num>.
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
