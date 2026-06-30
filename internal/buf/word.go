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

// input row is inclusive
// output is not inclusive
func (b *Buf) MoveByWord(loc Loc) (Loc, bool) {
	line := b.Line(loc.Row)
	if line == "" {
		return loc, false
	}
	for col := 0; line != ""; col++ {
		if col >= loc.Col {
			break
		}
		_, size := utf8.DecodeRuneInString(line)
		line = line[size:]
	}
	r, size := utf8.DecodeRuneInString(line)
	line = line[size:]
	kind := rkind.Kind(r)
	loc.Col++
	k := kind
	for ; line != ""; loc.Col++ {
		r, size = utf8.DecodeRuneInString(line)
		k = rkind.Kind(r)
		if k != kind {
			break
		}
		line = line[size:]
	}
	if line == "" {
		return loc, false
	}
	if kind == rkind.Blank || k != rkind.Blank {
		return loc, true
	}
	line = line[size:]
	loc.Col++
	for ; line != ""; loc.Col++ {
		r, size = utf8.DecodeRuneInString(line)
		k = rkind.Kind(r)
		if k != rkind.Blank {
			break
		}
		line = line[size:]
	}
	return loc, line != ""
}

// input row is inclusive
// output is not inclusive
func (b *Buf) MoveByWordAlt(loc Loc) (Loc, bool) {
	line := b.Line(loc.Row)
	if line == "" {
		return loc, false
	}
	for col := 0; line != ""; col++ {
		if col >= loc.Col {
			break
		}
		_, size := utf8.DecodeRuneInString(line)
		line = line[size:]
	}
	if line == "" {
		return loc, false
	}
	r, size := utf8.DecodeRuneInString(line)
	line = line[size:]
	kind := rkind.Kind(r)
	loc.Col++
	k := kind
	for ; line != ""; loc.Col++ {
		r, size = utf8.DecodeRuneInString(line)
		k = rkind.Kind(r)
		if k != kind {
			break
		}
		line = line[size:]
	}
	return loc, line != ""
}

// input row inclusive
// output is not inclusive
func (b *Buf) MoveBackwardByWord(loc Loc) (Loc, bool) {
	line := b.Line(loc.Row)
	if line == "" {
		return loc, false
	}
	col := utf8.RuneCountInString(line)
	for line != "" {
		col--
		if loc.Col >= col {
			break
		}
		_, size := utf8.DecodeLastRuneInString(line)
		line = line[:len(line)-size]
	}
	r, size := utf8.DecodeLastRuneInString(line)
	line = line[:len(line)-size]
	kind := rkind.Kind(r)
	if kind == rkind.Blank {
		return loc, false
	}
	for ; line != ""; loc.Col-- {
		r, size = utf8.DecodeLastRuneInString(line)
		k := rkind.Kind(r)
		if k != kind {
			break
		}
		line = line[:len(line)-size]
	}
	return loc, true
}
