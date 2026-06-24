package editor

func (ed *Editor) BeginMemory() {
	b := ed.Buf()
	b.Snapshot = append([]string{}, b.Lines...)
}

func (ed *Editor) EndMemory() {
	// do nothing
}

func (ed *Editor) CancelMemory() {
	// unimplemented
}
