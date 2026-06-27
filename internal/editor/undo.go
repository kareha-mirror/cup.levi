package editor

func (ed *Editor) BeginRecordForUndo() {
	b := ed.Buf()
	b.Snapshot = append([]string{}, b.Lines...)
}

func (ed *Editor) EndRecordForUndo() {
	// do nothing
}

func (ed *Editor) CancelRecordForUndo() {
	// unimplemented
}
