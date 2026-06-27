package buf

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/rkind"
)

func (b *Buf) SkipBlankLines(loc Loc) (Loc, bool) {
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
	loc.Row = max(numLines-1, 0)
	loc.Col = utf8.RuneCountInString(b.Line(loc.Row))
	return loc, false
}

func (b *Buf) SkipBackwardBlankLines(loc Loc) (Loc, bool) {
	for loc.Row >= 0 {
		rs := []rune(b.Line(loc.Row))
		if len(rs) > 0 {
			for col := loc.Col; col >= 0; col-- {
				r := rs[col]
				if !rkind.IsBlank(r) {
					loc.Col = col
					return loc, true
				}
			}
		}
		loc.Row--
		rc := utf8.RuneCountInString(b.Line(loc.Row))
		loc.Col = max(rc-1, 0)
	}
	loc.Row = 0
	loc.Col = 0
	return loc, false
}

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

func (b *Buf) MoveByWordEx(loc Loc) (Loc, bool) {
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
	return loc, loc.Col < len(rs)
}

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
