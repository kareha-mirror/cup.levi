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
	letter rune, n int,
) (buf.Loc, bool) {
	b := ed.Buf()
	col := findRuneCol(b.CurrentLine(), b.Loc.Col+1, letter)
	if col < 0 {
		return buf.Loc{}, false
	}
	return buf.Loc{col, b.Loc.Row}, true
}

// f<letter> : Find character <letter> forward in current line.
func (ed *Editor) MoveFindForward(letter rune, n int) (buf.Loc, bool) {
	ed.Commit()
	ed.find = Find{letter, false, false}
	return ed.internalMoveFindForward(letter, n)
}

func (ed *Editor) internalMoveFindBackward(
	letter rune, n int,
) (buf.Loc, bool) {
	b := ed.Buf()
	col := findBackwardRuneCol(b.CurrentLine(), b.Loc.Col-1, letter)
	if col < 0 {
		return buf.Loc{}, false
	}
	return buf.Loc{col, b.Loc.Row}, true
}

// F<letter> : Find character <letter> backward in current line.
func (ed *Editor) MoveFindBackward(letter rune, n int) (buf.Loc, bool) {
	ed.Commit()
	ed.find = Find{letter, true, false}
	return ed.internalMoveFindBackward(letter, n)
}

func (ed *Editor) internalMoveFindBeforeForward(
	letter rune, n int,
) (buf.Loc, bool) {
	b := ed.Buf()
	col := findRuneCol(b.CurrentLine(), b.Loc.Col+1, letter)
	if col < 0 {
		return buf.Loc{}, false
	}
	col = max(col-1, 0)
	return buf.Loc{col, b.Loc.Row}, true
}

// t<letter> : Find before character <letter> forward in current line.
func (ed *Editor) MoveFindBeforeForward(letter rune, n int) (buf.Loc, bool) {
	ed.Commit()
	ed.find = Find{letter, false, true}
	return ed.internalMoveFindBeforeForward(letter, n)
}

func (ed *Editor) internalMoveFindBeforeBackward(
	letter rune, n int,
) (buf.Loc, bool) {
	b := ed.Buf()
	col := findBackwardRuneCol(b.CurrentLine(), b.Loc.Col-1, letter)
	if col < 0 {
		return buf.Loc{}, false
	}
	col++
	return buf.Loc{col, b.Loc.Row}, true
}

// T<letter> : Find before character <letter> backward in current line.
func (ed *Editor) MoveFindBeforeBackward(letter rune, n int) (buf.Loc, bool) {
	ed.Commit()
	ed.find = Find{letter, true, true}
	return ed.internalMoveFindBeforeBackward(letter, n)
}

// ; : Find next match.
func (ed *Editor) MoveFindNextMatch(n int) (buf.Loc, bool) {
	ed.Commit()
	if ed.find.Letter == 0 {
		return buf.Loc{}, false
	}
	if ed.find.Backward {
		if ed.find.Before {
			return ed.internalMoveFindBeforeBackward(ed.find.Letter, n)
		} else {
			return ed.internalMoveFindBackward(ed.find.Letter, n)
		}
	} else {
		if ed.find.Before {
			return ed.internalMoveFindBeforeForward(ed.find.Letter, n)
		} else {
			return ed.internalMoveFindForward(ed.find.Letter, n)
		}
	}
}

// , : Find previous match.
func (ed *Editor) MoveFindPrevMatch(n int) (buf.Loc, bool) {
	ed.Commit()
	if ed.find.Letter == 0 {
		return buf.Loc{}, false
	}
	if ed.find.Backward {
		if ed.find.Before {
			return ed.internalMoveFindBeforeForward(ed.find.Letter, n)
		} else {
			return ed.internalMoveFindForward(ed.find.Letter, n)
		}
	} else {
		if ed.find.Before {
			return ed.internalMoveFindBeforeBackward(ed.find.Letter, n)
		} else {
			return ed.internalMoveFindBackward(ed.find.Letter, n)
		}
	}
}
