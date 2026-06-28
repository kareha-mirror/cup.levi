package buf

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/rkind"
)

// not inclusive
func (b *Buf) SkipBlanks(loc Loc) (Loc, bool) {
	numLines := b.NumLines()
	for loc.Row < numLines {
		col := 0
		for _, r := range b.Line(loc.Row) {
			if col >= loc.Col && !rkind.IsBlank(r) {
				loc.Col = col
				return loc, true
			}
			col++
		}
		loc.Row++
		loc.Col = 0
	}
	// confine row inclusive
	loc.Row = max(numLines-1, 0)
	loc.Col = utf8.RuneCountInString(b.Line(loc.Row))
	return loc, false
}

// row is inclusive
// col is not inclusive
func (b *Buf) SkipBackwardBlanks(loc Loc) (Loc, bool) {
	first := true
	for loc.Row >= 0 {
		line := b.Line(loc.Row)
		col := utf8.RuneCountInString(line)
		if !first {
			loc.Col = col
		}
		for len(line) > 0 {
			col--
			r, size := utf8.DecodeLastRuneInString(line)
			line = line[:len(line)-size]
			if col > loc.Col {
				continue
			}
			if !rkind.IsBlank(r) {
				loc.Col = col
				return loc, true
			}
		}
		loc.Row--
		first = false
	}
	loc.Row = 0
	loc.Col = 0
	return loc, false
}

// input is inclusive
// output is not inclusive
func (b *Buf) MoveByWord(loc Loc) (Loc, bool) {
	rs := []rune(b.Line(loc.Row))
	if len(rs) < 1 {
		return loc, false
	}
	kind := rkind.Kind(rs[loc.Col])
	loc.Col++
	k := kind
	for ; loc.Col < len(rs); loc.Col++ {
		k = rkind.Kind(rs[loc.Col])
		if k != kind {
			break
		}
	}
	if loc.Col >= len(rs) {
		return loc, false
	}
	if kind == rkind.Blank || k != rkind.Blank {
		return loc, true
	}
	loc.Col++
	for ; loc.Col < len(rs); loc.Col++ {
		k = rkind.Kind(rs[loc.Col])
		if k != rkind.Blank {
			break
		}
	}
	return loc, loc.Col < len(rs)
}

// input is inclusive
// output is not inclusive
func (b *Buf) MoveByWordEx(loc Loc) (Loc, bool) {
	rs := []rune(b.Line(loc.Row))
	if len(rs) < 1 || loc.Col >= len(rs) {
		return loc, false
	}
	kind := rkind.Kind(rs[loc.Col])
	loc.Col++
	k := kind
	for ; loc.Col < len(rs); loc.Col++ {
		k = rkind.Kind(rs[loc.Col])
		if k != kind {
			break
		}
	}
	return loc, loc.Col < len(rs)
}

// inclusive
func (b *Buf) MoveBackwardByWord(loc Loc) (Loc, bool) {
	rs := []rune(b.Line(loc.Row))
	if len(rs) < 1 {
		return loc, false
	}
	kind := rkind.Kind(rs[loc.Col])
	if kind == rkind.Blank {
		return loc, false
	}
	for ; loc.Col >= 0; loc.Col-- {
		k := rkind.Kind(rs[loc.Col])
		if k != kind {
			break
		}
	}
	loc.Col++
	return loc, true
}
