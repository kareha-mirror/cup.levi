package buf

import (
	"unicode/utf8"

	"tea.kareha.org/cup/termi/rkind"
)

func (b *Buf) IsRowIncluded(row int) bool {
	if row < 0 {
		return false
	}
	numLines := b.NumLines()
	if row < numLines {
		return true
	}
	// empty case
	if numLines == 0 && row == 0 {
		return true
	}
	return false
}

// not inclusive
func (b *Buf) confineRow(loc Loc) Loc {
	if loc.Row < 0 {
		return Loc{Col: loc.Col, Row: 0}
	}
	numLines := b.NumLines()
	if loc.Row > numLines {
		return Loc{Col: loc.Col, Row: numLines}
	}
	return loc
}

// not inclusive
func (b *Buf) confineCol(loc Loc) Loc {
	if loc.Col < 0 {
		return Loc{Col: 0, Row: loc.Row}
	}
	rc := utf8.RuneCountInString(b.Line(loc.Row))
	if loc.Col > rc {
		return Loc{Col: rc, Row: loc.Row}
	}
	return loc
}

// not inclusive
func (b *Buf) Confine(loc Loc) Loc {
	loc = b.confineRow(loc)
	loc = b.confineCol(loc)
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

func (b *Buf) FirstNonBlankCol(row int) int {
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
func (b *Buf) ConfineFreeCol(loc Loc) Loc {
	rc := utf8.RuneCountInString(b.Line(loc.Row))
	if b.VirtCol < rc {
		return Loc{Col: b.VirtCol, Row: loc.Row}
	}
	return Loc{Col: max(rc-1, 0), Row: loc.Row}
}
