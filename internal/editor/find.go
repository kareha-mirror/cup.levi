package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/rutil"
)

type Find struct {
	Letter   rune
	Backward bool
	Before   bool
}

////////////////////////////////
// Character Finding Commands //
////////////////////////////////

func (ed *Editor) internalFindForward(
	loc buf.Loc, letter rune,
) (buf.Loc, bool) {
	col := rutil.RuneIndex(ed.Buf().Line(loc.Row), loc.Col+1, letter)
	if col < 0 {
		return loc, false
	}
	return buf.Loc{col, loc.Row}, true
}

// f<letter> : Find character <letter> forward in current line.
func (ed *Editor) FindForward(letter rune, n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("FindForward: n < 1")
		return buf.Loc{}, false
	}
	ed.find = Find{letter, false, false}
	loc := ed.Buf().Loc
	var ok bool
	for i := 0; i < n; i++ {
		loc, ok = ed.internalFindForward(loc, letter)
		if !ok {
			if i == 0 {
				ed.Notice("Not found")
				return buf.Loc{}, false
			}
			break
		}
	}
	return loc, true
}

func (ed *Editor) internalFindBackward(
	loc buf.Loc, letter rune,
) (buf.Loc, bool) {
	col := rutil.LastRuneIndex(ed.Buf().Line(loc.Row), loc.Col-1, letter)
	if col < 0 {
		return loc, false
	}
	return buf.Loc{col, loc.Row}, true
}

// F<letter> : Find character <letter> backward in current line.
func (ed *Editor) FindBackward(letter rune, n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("FindBackward: n < 1")
		return buf.Loc{}, false
	}
	ed.find = Find{letter, true, false}
	loc := ed.Buf().Loc
	var ok bool
	for i := 0; i < n; i++ {
		loc, ok = ed.internalFindBackward(loc, letter)
		if !ok {
			if i == 0 {
				ed.Notice("Not found")
				return buf.Loc{}, false
			}
			break
		}
	}
	return loc, true
}

func (ed *Editor) internalFindBeforeForward(
	loc buf.Loc, letter rune,
) (buf.Loc, bool) {
	col := rutil.RuneIndex(ed.Buf().Line(loc.Row), loc.Col+1, letter)
	if col < 0 {
		return loc, false
	}
	col = max(col-1, 0)
	return buf.Loc{col, loc.Row}, true
}

// t<letter> : Find before character <letter> forward in current line.
func (ed *Editor) FindBeforeForward(letter rune, n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("FindBeforeForward: n < 1")
		return buf.Loc{}, false
	}
	ed.find = Find{letter, false, true}
	loc := ed.Buf().Loc
	var ok bool
	for i := 0; i < n; i++ {
		loc, ok = ed.internalFindBeforeForward(loc, letter)
		if !ok {
			if i == 0 {
				ed.Notice("Not found")
				return buf.Loc{}, false
			}
			break
		}
	}
	return loc, true
}

func (ed *Editor) internalFindBeforeBackward(
	loc buf.Loc, letter rune,
) (buf.Loc, bool) {
	col := rutil.LastRuneIndex(ed.Buf().Line(loc.Row), loc.Col-1, letter)
	if col < 0 {
		return loc, false
	}
	col++
	return buf.Loc{col, loc.Row}, true
}

// T<letter> : Find before character <letter> backward in current line.
func (ed *Editor) FindBeforeBackward(letter rune, n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("FindBeforeBackward: n < 1")
		return buf.Loc{}, false
	}
	ed.find = Find{letter, true, true}
	loc := ed.Buf().Loc
	var ok bool
	for i := 0; i < n; i++ {
		loc, ok = ed.internalFindBeforeBackward(loc, letter)
		if !ok {
			if i == 0 {
				ed.Notice("Not found")
				return buf.Loc{}, false
			}
			break
		}
	}
	return loc, true
}

// ; : Find next match.
func (ed *Editor) FindNext(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("FindNext: n < 1")
		return buf.Loc{}, false
	}
	if ed.find.Letter == 0 {
		return buf.Loc{}, false
	}
	loc := ed.Buf().Loc
	var ok bool
	for i := 0; i < n; i++ {
		if ed.find.Backward {
			if ed.find.Before {
				loc, ok = ed.internalFindBeforeBackward(
					loc, ed.find.Letter,
				)
			} else {
				loc, ok = ed.internalFindBackward(
					loc, ed.find.Letter,
				)
			}
		} else {
			if ed.find.Before {
				loc, ok = ed.internalFindBeforeForward(
					loc, ed.find.Letter,
				)
			} else {
				loc, ok = ed.internalFindForward(
					loc, ed.find.Letter,
				)
			}
		}
		if !ok {
			if i == 0 {
				ed.Notice("Not found")
				return buf.Loc{}, false
			}
			break
		}
	}
	return loc, true
}

// , : Find previous match.
func (ed *Editor) FindPrev(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("FindPrev: n < 1")
		return buf.Loc{}, false
	}
	if ed.find.Letter == 0 {
		return buf.Loc{}, false
	}
	loc := ed.Buf().Loc
	var ok bool
	for i := 0; i < n; i++ {
		if ed.find.Backward {
			if ed.find.Before {
				loc, ok = ed.internalFindBeforeForward(
					loc, ed.find.Letter,
				)
			} else {
				loc, ok = ed.internalFindForward(
					loc, ed.find.Letter,
				)
			}
		} else {
			if ed.find.Before {
				loc, ok = ed.internalFindBeforeBackward(
					loc, ed.find.Letter,
				)
			} else {
				loc, ok = ed.internalFindBackward(
					loc, ed.find.Letter,
				)
			}
		}
		if !ok {
			if i == 0 {
				ed.Notice("Not found")
				return buf.Loc{}, false
			}
			break
		}
	}
	return loc, true
}
