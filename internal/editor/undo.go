package editor

func (ed *Editor) BeginRecordForUndo() {
	ed.Buf().BeginSnapshot()
}

func (ed *Editor) EndRecordForUndo() {
	ed.Buf().EndSnapshot()
}

func (ed *Editor) CancelRecordForUndo() {
	ed.Buf().CancelSnapshot()
}
