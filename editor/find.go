package editor

import (
	"tea.kareha.org/cup/termi/rutil"

	"tea.kareha.org/cup/levi/internal/buf"
)

type findState struct {
	// order matters
	active   bool
	r        rune
	backward bool
	before   bool
}

////////////////////////////////
// Character Finding Commands //
////////////////////////////////

func (ed *Editor) find(loc buf.Loc, r rune) (buf.Loc, bool) {
	col := rutil.RuneIndex(ed.Buf().Line(loc.Row), loc.Col+1, r)
	if col < 0 {
		return loc, false
	}
	return buf.Loc{col, loc.Row}, true
}

// f<char> : Find character <char> in current line and move to it.
func (ed *Editor) Find(r rune, n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("Find: n < 1")
		return buf.Loc{}, false
	}
	ed.finds = findState{true, r, false, false}
	loc := ed.Buf().Loc
	var ok bool
	for i := 0; i < n; i++ {
		loc, ok = ed.find(loc, r)
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

func (ed *Editor) findBackward(loc buf.Loc, r rune) (buf.Loc, bool) {
	col := rutil.LastRuneIndex(ed.Buf().Line(loc.Row), loc.Col-1, r)
	if col < 0 {
		return loc, false
	}
	return buf.Loc{col, loc.Row}, true
}

// F<char> : Find character <char> backward in current line and move to it.
func (ed *Editor) FindBackward(r rune, n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("FindBackward: n < 1")
		return buf.Loc{}, false
	}
	ed.finds = findState{true, r, true, false}
	loc := ed.Buf().Loc
	var ok bool
	for i := 0; i < n; i++ {
		loc, ok = ed.findBackward(loc, r)
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

func (ed *Editor) findBefore(loc buf.Loc, r rune) (buf.Loc, bool) {
	col := rutil.RuneIndex(ed.Buf().Line(loc.Row), loc.Col+1, r)
	if col < 0 {
		return loc, false
	}
	col = max(col-1, 0)
	return buf.Loc{col, loc.Row}, true
}

// t<char> : Find character <char> in current line and move to before it.
func (ed *Editor) FindBefore(r rune, n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("FindBefore: n < 1")
		return buf.Loc{}, false
	}
	ed.finds = findState{true, r, false, true}
	loc := ed.Buf().Loc
	var ok bool
	for i := 0; i < n; i++ {
		loc, ok = ed.findBefore(loc, r)
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

func (ed *Editor) findBeforeBackward(loc buf.Loc, r rune) (buf.Loc, bool) {
	col := rutil.LastRuneIndex(ed.Buf().Line(loc.Row), loc.Col-1, r)
	if col < 0 {
		return loc, false
	}
	col++
	return buf.Loc{col, loc.Row}, true
}

// T<char> : Find character <char> backward in current line and move before it.
func (ed *Editor) FindBeforeBackward(r rune, n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("FindBeforeBackward: n < 1")
		return buf.Loc{}, false
	}
	ed.finds = findState{true, r, true, true}
	loc := ed.Buf().Loc
	var ok bool
	for i := 0; i < n; i++ {
		loc, ok = ed.findBeforeBackward(loc, r)
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

// ; : Repeat find operation to find next match.
func (ed *Editor) FindNext(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("FindNext: n < 1")
		return buf.Loc{}, false
	}
	if !ed.finds.active {
		return buf.Loc{}, false
	}
	loc := ed.Buf().Loc
	var ok bool
	for i := 0; i < n; i++ {
		if ed.finds.backward {
			if ed.finds.before {
				loc, ok = ed.findBeforeBackward(loc, ed.finds.r)
			} else {
				loc, ok = ed.findBackward(loc, ed.finds.r)
			}
		} else {
			if ed.finds.before {
				loc, ok = ed.findBefore(loc, ed.finds.r)
			} else {
				loc, ok = ed.find(loc, ed.finds.r)
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

// , : Repeat find operation to find previous match.
func (ed *Editor) FindPrev(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("FindPrev: n < 1")
		return buf.Loc{}, false
	}
	if !ed.finds.active {
		return buf.Loc{}, false
	}
	loc := ed.Buf().Loc
	var ok bool
	for i := 0; i < n; i++ {
		if ed.finds.backward {
			if ed.finds.before {
				loc, ok = ed.findBefore(loc, ed.finds.r)
			} else {
				loc, ok = ed.find(loc, ed.finds.r)
			}
		} else {
			if ed.finds.before {
				loc, ok = ed.findBeforeBackward(loc, ed.finds.r)
			} else {
				loc, ok = ed.findBackward(loc, ed.finds.r)
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
