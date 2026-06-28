package buf

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/rkind"
)

func (b *Buf) CheckRowInclusive(row int) bool {
	if row < 0 {
		return false
	}
	numLines := b.NumLines()
	if row < numLines {
		return true
	}
	if numLines == 0 && row == 0 {
		return true
	}
	return false
}

// not inclusive
func (b *Buf) ConfineRow(row int) int {
	if row < 0 {
		return 0
	}
	numLines := b.NumLines()
	if row > numLines {
		return numLines
	}
	return row
}

// not inclusive
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

// not inclusive
func (b *Buf) Confine(loc Loc) Loc {
	loc.Row = b.ConfineRow(loc.Row)
	loc.Col = b.ConfineCol(loc)
	return loc
}

func (b *Buf) ConfineInclusive(loc Loc) Loc {
	if loc.Row >= b.NumLines() {
		loc.Row = max(b.NumLines()-1, 0)
		rc := utf8.RuneCountInString(b.Line(loc.Row))
		loc.Col = max(rc-1, 0)
		return loc
	}
	rc := utf8.RuneCountInString(b.Line(loc.Row))
	if loc.Col >= rc {
		loc.Col = rc - 1
	}
	loc.Col = max(loc.Col, 0)
	return loc
}

func (b *Buf) NonBlankColOfLine(row int) int {
	col := 0
	for _, r := range b.Line(row) {
		if !rkind.IsBlank(r) {
			break
		}
		col++
	}
	return col
}

// inclusive
func (b *Buf) ConfineColVirtInclusive(row int) int {
	rc := utf8.RuneCountInString(b.Line(row))
	if b.VirtCol < rc {
		return b.VirtCol
	} else {
		return max(rc-1, 0)
	}
}
