package buf

import (
	"unicode/utf8"
)

// not care if inclusive or not
func OrderRegion(start, end Loc) (Loc, Loc) {
	if start.Row < end.Row {
		return start, end
	}
	if end.Row < start.Row {
		return end, start
	}
	// start.Row == end.Row
	if start.Col < end.Col {
		return start, end
	}
	return end, start
}

// inclusive or not selectable
func (b *Buf) ConfineRegion(start, end Loc, inclusive bool) (Loc, Loc) {
	start, end = OrderRegion(start, end)
	if inclusive {
		return b.ConfineInclusive(start), b.ConfineInclusive(end)
	}
	// XXX start should be inclusive or not?
	start, end = b.Confine(start), b.Confine(end)
	if start.Row < end.Row && end.Col < 1 {
		end.Row--
		end.Col = utf8.RuneCountInString(b.Line(end.Row))
	}
	return start, end
}

// not inclusive
func (b *Buf) RegionRunewise(start, end Loc) []string {
	if start.Row == end.Row {
		rs := []rune(b.Line(start.Row))
		return []string{string(rs[start.Col:end.Col])}
	}
	rs := []rune(b.Line(start.Row))
	lines := append([]string{}, string(rs[start.Col:]))
	for row := start.Row + 1; row < end.Row; row++ {
		lines = append(lines, b.Line(row))
	}
	rs = []rune(b.Line(end.Row))
	lines = append(lines, string(rs[:end.Col]))
	return lines
}
