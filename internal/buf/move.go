package buf

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/rkind"
)

func (b *Buf) CheckRowInclusive(row int) bool {
	if row < 0 {
		return false
	}
	numLines := len(b.Lines)
	if row < numLines {
		return true
	}
	if numLines == 0 && row == 0 {
		return true
	}
	return false
}

func (b *Buf) ConfineRow(row int) int {
	if row < 0 {
		return 0
	}
	n := len(b.Lines)
	if row > n {
		return n
	}
	return row
}

func (b *Buf) ConfineCol(loc Loc) int {
	if loc.Col < 0 {
		return 0
	}
	rc := utf8.RuneCountInString(b.Line(loc.Row))
	if loc.Col > rc {
		return rc
	}
	return loc.Col
}

func (b *Buf) Confine(loc Loc) Loc {
	loc.Row = b.ConfineRow(loc.Row)
	loc.Col = b.ConfineCol(loc)
	return loc
}

func (b *Buf) ConfineInclusive(loc Loc) Loc {
	if loc.Row >= b.NumLines() {
		loc.Row = max(b.NumLines()-1, 0)
		line := b.Line(loc.Row)
		rc := utf8.RuneCountInString(line)
		loc.Col = max(rc-1, 0)
		return loc
	}
	line := b.Line(loc.Row)
	rc := utf8.RuneCountInString(line)
	if loc.Col >= rc {
		loc.Col = rc - 1
	}
	loc.Col = max(loc.Col, 0)
	return loc
}

func (b *Buf) NonBlankColOfLine(row int) int {
	line := b.Line(row)
	col := 0
	for _, r := range line {
		if !rkind.IsBlank(r) {
			break
		}
		col++
	}
	return col
}
