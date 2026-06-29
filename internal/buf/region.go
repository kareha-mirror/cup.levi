package buf

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/rutil"
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
		// caller may adjust as end.Col++ to use as if not inclusive
	}
	// start is virtually inclusive
	start, end = b.Confine(start), b.Confine(end)
	// make row inclusive
	if start.Row < end.Row && end.Col == 0 {
		end.Row--
		end.Col = utf8.RuneCountInString(b.Line(end.Row))
	}
	return start, end
}

// row is inclusive
// col is not inclusive
func (b *Buf) RegionRunewise(start, end Loc) []string {
	if start.Row == end.Row {
		s := rutil.Body(b.Line(start.Row), start.Col, end.Col)
		return []string{s}
	}
	s := rutil.Tail(b.Line(start.Row), start.Col)
	lines := append([]string{}, s)
	for row := start.Row + 1; row < end.Row; row++ {
		lines = append(lines, b.Line(row))
	}
	s = rutil.Head(b.Line(end.Row), end.Col)
	lines = append(lines, s)
	return lines
}
