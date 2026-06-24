package buf

import (
	"unicode/utf8"
)

type Dest struct {
	Loc       Loc
	Linewise  bool
	FreeCol   bool
	Inclusive bool
}

func (b *Buf) CheckRow(row int) bool {
	if row < 0 {
		return false
	}
	numLines := len(b.Lines)
	if numLines == 0 && row == 0 {
		return true
	}
	if row >= numLines {
		return false
	}
	return true
}

func (b *Buf) ConfineRow(row int) int {
	if row < 0 {
		return 0
	}
	n := len(b.Lines)
	if row >= n {
		return max(n-1, 0)
	}
	return row
}

func (b *Buf) ConfineCol(loc Loc) int {
	if loc.Col < 0 {
		return 0
	}
	rc := utf8.RuneCountInString(b.Line(loc.Row))
	if loc.Col >= rc {
		return max(rc-1, 0)
	}
	return loc.Col
}

func (b *Buf) Confine(loc Loc) Loc {
	loc.Row = b.ConfineRow(loc.Row)
	loc.Col = b.ConfineCol(loc)
	return loc
}
