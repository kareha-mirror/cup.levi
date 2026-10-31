package buf

import (
	"unicode/utf8"

	"tea.kareha.org/cup/termi/rutil"
)

// not care if inclusive or not
func normalizeRegion(start, end Loc) (Loc, Loc) {
	if start.Row < end.Row {
		return start, end
	} else if end.Row < start.Row {
		return end, start
	} else if start.Col < end.Col { // start.Row == end.Row
		return start, end
	} else {
		return end, start
	}
}

func (b *Buf) ConfineRegion(
	start, end Loc, inclusive bool, linewise bool,
) (Loc, Loc) {
	start, end = normalizeRegion(start, end)
	if inclusive {
		return b.ConfineInclusive(start), b.ConfineInclusive(end)
		// caller may adjust as end.Col++ to use as if not inclusive
	}

	// start is virtually inclusive
	start, end = b.Confine(start), b.Confine(end)
	if linewise {
		// make row inclusive
		if start.Row < end.Row && end.Col == 0 {
			end.Row--
			end.Col = utf8.RuneCountInString(b.Line(end.Row))
		}
	}
	return start, end
}

// row is inclusive
// col is not inclusive
func (b *Buf) runewiseContent(start, end Loc) []string {
	if start.Row == end.Row {
		body := rutil.Body(b.Line(start.Row), start.Col, end.Col)
		return []string{body}
	}

	tail := rutil.Tail(b.Line(start.Row), start.Col)
	lines := append([]string{}, tail)
	for row := start.Row + 1; row < end.Row; row++ {
		lines = append(lines, b.Line(row))
	}
	head := rutil.Head(b.Line(end.Row), end.Col)
	lines = append(lines, head)
	return lines
}
