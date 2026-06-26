package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
)

type Find struct {
	Letter   rune
	Backward bool
	Before   bool
}

func findRuneCol(line string, start int, r rune) int {
	col := 0
	for _, ru := range line {
		if col < start {
			col++
			continue
		}
		if ru == r {
			return col
		}
		col++
	}
	return -1
}

func findBackwardRuneCol(line string, start int, r rune) int {
	rs := []rune(line)
	for col := start; col >= 0; col-- {
		if rs[col] == r {
			return col
		}
	}
	return -1
}

////////////////////////////////
// Character Finding Commands //
////////////////////////////////

func (ed *Editor) internalMoveFindForward(
	loc buf.Loc, letter rune,
) (buf.Loc, bool) {
	b := ed.Buf()
	col := findRuneCol(b.Line(loc.Row), loc.Col+1, letter)
	if col < 0 {
		return loc, false
	}
	return buf.Loc{col, loc.Row}, true
}

// f<letter> : Find character <letter> forward in current line.
func (ed *Editor) MoveFindForward(letter rune, n int) (buf.Loc, bool) {
	ed.Commit()
	if n < 1 {
		ed.Error("MoveFindForward: n < 1")
		return buf.Loc{}, false
	}
	ed.find = Find{letter, false, false}
	b := ed.Buf()
	loc := b.Loc
	var ok bool
	for i := 0; i < n; i++ {
		loc, ok = ed.internalMoveFindForward(loc, letter)
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

func (ed *Editor) internalMoveFindBackward(
	loc buf.Loc, letter rune,
) (buf.Loc, bool) {
	b := ed.Buf()
	col := findBackwardRuneCol(b.Line(loc.Row), loc.Col-1, letter)
	if col < 0 {
		return loc, false
	}
	return buf.Loc{col, loc.Row}, true
}

// F<letter> : Find character <letter> backward in current line.
func (ed *Editor) MoveFindBackward(letter rune, n int) (buf.Loc, bool) {
	ed.Commit()
	if n < 1 {
		ed.Error("MoveFindBackward: n < 1")
		return buf.Loc{}, false
	}
	ed.find = Find{letter, true, false}
	b := ed.Buf()
	loc := b.Loc
	var ok bool
	for i := 0; i < n; i++ {
		loc, ok = ed.internalMoveFindBackward(loc, letter)
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

func (ed *Editor) internalMoveFindBeforeForward(
	loc buf.Loc, letter rune,
) (buf.Loc, bool) {
	b := ed.Buf()
	col := findRuneCol(b.Line(loc.Row), loc.Col+1, letter)
	if col < 0 {
		return loc, false
	}
	col = max(col-1, 0)
	return buf.Loc{col, loc.Row}, true
}

// t<letter> : Find before character <letter> forward in current line.
func (ed *Editor) MoveFindBeforeForward(letter rune, n int) (buf.Loc, bool) {
	ed.Commit()
	if n < 1 {
		ed.Error("MoveFindBeforeForward: n < 1")
		return buf.Loc{}, false
	}
	ed.find = Find{letter, false, true}
	b := ed.Buf()
	loc := b.Loc
	var ok bool
	for i := 0; i < n; i++ {
		loc, ok = ed.internalMoveFindBeforeForward(loc, letter)
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

func (ed *Editor) internalMoveFindBeforeBackward(
	loc buf.Loc, letter rune,
) (buf.Loc, bool) {
	b := ed.Buf()
	col := findBackwardRuneCol(b.Line(loc.Row), loc.Col-1, letter)
	if col < 0 {
		return loc, false
	}
	col++
	return buf.Loc{col, loc.Row}, true
}

// T<letter> : Find before character <letter> backward in current line.
func (ed *Editor) MoveFindBeforeBackward(letter rune, n int) (buf.Loc, bool) {
	ed.Commit()
	if n < 1 {
		ed.Error("MoveFindBeforeBackward: n < 1")
		return buf.Loc{}, false
	}
	ed.find = Find{letter, true, true}
	b := ed.Buf()
	loc := b.Loc
	var ok bool
	for i := 0; i < n; i++ {
		loc, ok = ed.internalMoveFindBeforeBackward(loc, letter)
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
func (ed *Editor) MoveFindNextMatch(n int) (buf.Loc, bool) {
	ed.Commit()
	if n < 1 {
		ed.Error("MoveFindNextMatch: n < 1")
		return buf.Loc{}, false
	}
	if ed.find.Letter == 0 {
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var ok bool
	for i := 0; i < n; i++ {
		if ed.find.Backward {
			if ed.find.Before {
				loc, ok = ed.internalMoveFindBeforeBackward(
					loc, ed.find.Letter,
				)
			} else {
				loc, ok = ed.internalMoveFindBackward(
					loc, ed.find.Letter,
				)
			}
		} else {
			if ed.find.Before {
				loc, ok = ed.internalMoveFindBeforeForward(
					loc, ed.find.Letter,
				)
			} else {
				loc, ok = ed.internalMoveFindForward(
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
func (ed *Editor) MoveFindPrevMatch(n int) (buf.Loc, bool) {
	ed.Commit()
	if n < 1 {
		ed.Error("MoveFindPrevMatch: n < 1")
		return buf.Loc{}, false
	}
	if ed.find.Letter == 0 {
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var ok bool
	for i := 0; i < n; i++ {
		if ed.find.Backward {
			if ed.find.Before {
				loc, ok = ed.internalMoveFindBeforeForward(
					loc, ed.find.Letter,
				)
			} else {
				loc, ok = ed.internalMoveFindForward(
					loc, ed.find.Letter,
				)
			}
		} else {
			if ed.find.Before {
				loc, ok = ed.internalMoveFindBeforeBackward(
					loc, ed.find.Letter,
				)
			} else {
				loc, ok = ed.internalMoveFindBackward(
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
