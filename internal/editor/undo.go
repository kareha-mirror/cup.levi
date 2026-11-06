package editor

func (ed *Editor) BeginUndoRecord() {
	ed.Buf().BeginSnapshot()
}

func (ed *Editor) EndUndoRecord() {
	ed.Buf().EndSnapshot()
}

func (ed *Editor) CancelUndoRecord() {
	ed.Buf().CancelSnapshot()
}
