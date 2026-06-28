package editor

func (ed *Editor) BeginRecordForUndo() {
	ed.Buf().TakeSnapshot()
}

func (ed *Editor) EndRecordForUndo() {
	// do nothing
}

func (ed *Editor) CancelRecordForUndo() {
	ed.Buf().CancelSnapshot()
}
