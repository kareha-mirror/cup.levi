package buf

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/rkind"
)

func (b *Buf) SkipBlankLines(loc Loc) (Loc, bool) {
	for loc.Row < b.NumLines() {
		line := b.Line(loc.Row)
		col := 0
		for _, r := range line {
			if col >= loc.Col && !rkind.IsBlank(r) {
				loc.Col = col
				return loc, true
			}
			col++
		}
		loc.Row++
		loc.Col = 0
	}
	loc.Row = max(b.NumLines()-1, 0)
	loc.Col = max(utf8.RuneCountInString(b.Line(loc.Row))-1, 0)
	return loc, false
}

func (b *Buf) SkipBackwardBlankLines(loc Loc) (Loc, bool) {
	for loc.Row >= 0 {
		line := b.Line(loc.Row)
		if line != "" {
			rs := []rune(line)
			col := loc.Col
			for ; col >= 0; col-- {
				r := rs[col]
				if !rkind.IsBlank(r) {
					loc.Col = col
					return loc, true
				}
			}
		}
		loc.Row--
		line = b.Line(loc.Row)
		loc.Col = max(utf8.RuneCountInString(line)-1, 0)
	}
	loc.Row = 0
	loc.Col = 0
	return loc, false
}

func (b *Buf) MoveByWord(loc Loc) (Loc, bool) {
	line := b.Line(loc.Row)
	if len(line) < 1 {
		return loc, false
	}
	rs := []rune(line)
	col := loc.Col
	kind := rkind.Kind(rs[col])
	col++
	k := kind
	for ; col < len(rs); col++ {
		k = rkind.Kind(rs[col])
		if k != kind {
			break
		}
	}
	if col >= len(rs) {
		return loc, false
	}
	if kind == rkind.Blank || k != rkind.Blank {
		loc.Col = col
		return loc, true
	}
	col++
	for ; col < len(rs); col++ {
		k = rkind.Kind(rs[col])
		if k != rkind.Blank {
			break
		}
	}
	if col < len(rs) {
		loc.Col = col
		return loc, true
	}
	return loc, false
}

func (b *Buf) MoveByWordEx(loc Loc) (Loc, bool) {
	line := b.Line(loc.Row)
	if len(line) < 1 {
		return loc, false
	}
	rs := []rune(line)
	col := loc.Col
	kind := rkind.Kind(rs[col])
	col++
	k := kind
	for ; col < len(rs); col++ {
		k = rkind.Kind(rs[col])
		if k != kind {
			break
		}
	}
	if col >= len(rs) {
		loc.Col = len(rs)
		return loc, false
	}
	loc.Col = col
	return loc, true
}

func (b *Buf) MoveBackwardByWord(loc Loc) (Loc, bool) {
	line := b.Line(loc.Row)
	if len(line) < 1 {
		return loc, false
	}
	rs := []rune(line)
	kind := rkind.Kind(rs[loc.Col])
	if kind == rkind.Blank {
		return loc, false
	}
	col := loc.Col
	for ; col >= 0; col-- {
		k := rkind.Kind(rs[col])
		if k != kind {
			break
		}
	}
	loc.Col = col + 1
	return loc, true
}
