package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
)

////////////////////////////////
// Character Finding Commands //
////////////////////////////////

// f<letter> : Find character <letter> forward in current line.
func (ed *Editor) MoveFindForward(letter rune, n int) (buf.Loc, bool) {
	ed.Commit()
	ed.Unimplemented("MoveFindForward")
	return buf.Loc{}, false
}

// F<letter> : Find character <letter> backward in current line.
func (ed *Editor) MoveFindBackward(letter rune, n int) (buf.Loc, bool) {
	ed.Commit()
	ed.Unimplemented("MoveFindBackward")
	return buf.Loc{}, false
}

// t<letter> : Find before character <letter> forward in current line.
func (ed *Editor) MoveFindBeforeForward(letter rune, n int) (buf.Loc, bool) {
	ed.Commit()
	ed.Unimplemented("MoveFindBeforeForward")
	return buf.Loc{}, false
}

// T<letter> : Find before character <letter> backward in current line.
func (ed *Editor) MoveFindBeforeBackward(letter rune, n int) (buf.Loc, bool) {
	ed.Commit()
	ed.Unimplemented("MoveFindBeforeBackward")
	return buf.Loc{}, false
}

// ; : Find next match.
func (ed *Editor) MoveFindNextMatch(n int) (buf.Loc, bool) {
	ed.Commit()
	ed.Unimplemented("MoveFindNextMatch")
	return buf.Loc{}, false
}

// , : Find previous match.
func (ed *Editor) MoveFindPrevMatch(n int) (buf.Loc, bool) {
	ed.Commit()
	ed.Unimplemented("MoveFindPrevMatch")
	return buf.Loc{}, false
}
