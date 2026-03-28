package editor

//////////////////////
// Editing Commands //
//////////////////////

// r : Replace single character under cursor.
func (ed *Editor) EditReplace(letter rune, n int) {
	ed.EnsureCommand()
	ed.Unimplemented("EditReplace")
}

// J : Join current line with next line.
func (ed *Editor) EditJoin(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("EditJoin")
}

// >> : Indent current line.
func (ed *Editor) EditIndent(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("EditIndent")
}

// << : Outdent current line.
func (ed *Editor) EditOutdent(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("EditOutdent")
}

// > <mv> : Indent region from current cursor to destination of motion <mv>.
func (ed *Editor) EditIndentRegion(start Loc, end Loc) {
	ed.EnsureCommand()
	ed.Unimplemented("EditIndentRegion")
}

// < <mv> : Outdent region from current cursor to destination of motion <mv>.
func (ed *Editor) EditOutdentRegion(start Loc, end Loc) {
	ed.EnsureCommand()
	ed.Unimplemented("EditOutdentRegion")
}
