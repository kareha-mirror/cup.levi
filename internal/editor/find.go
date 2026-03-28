package editor

////////////////////////////////
// Character Finding Commands //
////////////////////////////////

// f<letter> : Find character <letter> forward in current line.
func (ed *Editor) FindForward(letter rune, n int) {
	ed.EnsureCommand()
	ed.Unimplemented("FindForward")
}

// F<letter> : Find character <letter> backward in current line.
func (ed *Editor) FindBackward(letter rune, n int) {
	ed.EnsureCommand()
	ed.Unimplemented("FindBackward")
}

// t<letter> : Find before character <letter> forward in current line.
func (ed *Editor) FindBeforeForward(letter rune, n int) {
	ed.EnsureCommand()
	ed.Unimplemented("FindBeforeForward")
}

// T<letter> : Find before character <letter> backward in current line.
func (ed *Editor) FindBeforeBackward(letter rune, n int) {
	ed.EnsureCommand()
	ed.Unimplemented("FindBeforeBackward")
}

// ; : Find next match.
func (ed *Editor) FindNextMatch(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("FindNextMatch")
}

// , : Find previous match.
func (ed *Editor) FindPrevMatch(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("FindPrevMatch")
}
