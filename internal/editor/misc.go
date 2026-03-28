package editor

////////////////////////////
// Miscellaneous Commands //
////////////////////////////

// Ctrl-g : Show info such as current cursor position.
func (ed *Editor) MiscShowInfo() {
	ed.EnsureCommand()
	ed.Unimplemented("MiscShowInfo")
}

// . : Repeat last edit.
func (ed *Editor) MiscRepeat(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MiscRepeat")
}

// u : Undo.
func (ed *Editor) MiscUndo(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MiscUndo")
}

// U : Restore current line to previous state.
func (ed *Editor) MiscRestore() {
	ed.EnsureCommand()
	ed.Unimplemented("MiscRestore")
}

// ZZ : Save and quit.
func (ed *Editor) MiscSaveAndQuit() {
	ed.EnsureCommand()
	ed.quit = true
}
