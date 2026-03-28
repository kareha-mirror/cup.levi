package editor

/////////////////////
// Search Commands //
/////////////////////

// /<pattern> Enter : Search <pattern> forward.
func (ed *Editor) SearchForward(s string) {
	ed.EnsureCommand()
	ed.Unimplemented("SearchForward")
}

// ?<pattern> Enter : Search <pattern> backward.
func (ed *Editor) SearchBackward(s string) {
	ed.EnsureCommand()
	ed.Unimplemented("SearchBackward")
}

// n : Search next match.
func (ed *Editor) SearchNextMatch() {
	ed.EnsureCommand()
	ed.Unimplemented("SearchNextMatch")
}

// N : Search previous match.
func (ed *Editor) SearchPrevMatch() {
	ed.EnsureCommand()
	ed.Unimplemented("SearchPrevMatch")
}

// / Enter : Repeat last search forward.
func (ed *Editor) SearchRepeatForward() {
	ed.EnsureCommand()
	ed.Unimplemented("SearchRepeatForward")
}

// ? Enter : Repeat last search backward.
func (ed *Editor) SearchRepeatBackward() {
	ed.EnsureCommand()
	ed.Unimplemented("SearchRepeatBackward")
}
