package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
)

////////////////////////////////
// Character Finding Commands //
////////////////////////////////

// f<letter> : Find character <letter> forward in current line.
func (ed *Editor) MoveFindForward(letter rune, n int) (buf.Dest, bool) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveFindForward")
	return buf.Dest{}, false
}

// F<letter> : Find character <letter> backward in current line.
func (ed *Editor) MoveFindBackward(letter rune, n int) (buf.Dest, bool) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveFindBackward")
	return buf.Dest{}, false
}

// t<letter> : Find before character <letter> forward in current line.
func (ed *Editor) MoveFindBeforeForward(letter rune, n int) (buf.Dest, bool) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveFindBeforeForward")
	return buf.Dest{}, false
}

// T<letter> : Find before character <letter> backward in current line.
func (ed *Editor) MoveFindBeforeBackward(letter rune, n int) (buf.Dest, bool) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveFindBeforeBackward")
	return buf.Dest{}, false
}

// ; : Find next match.
func (ed *Editor) MoveFindNextMatch(n int) (buf.Dest, bool) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveFindNextMatch")
	return buf.Dest{}, false
}

// , : Find previous match.
func (ed *Editor) MoveFindPrevMatch(n int) (buf.Dest, bool) {
	ed.EnsureCommand()
	ed.Unimplemented("MoveFindPrevMatch")
	return buf.Dest{}, false
}
