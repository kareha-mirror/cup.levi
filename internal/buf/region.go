package buf

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

func (b *Buf) ConfineRegion(start, end Loc, inclusive bool) (Loc, Loc) {
	start, end = OrderRegion(start, end)
	if inclusive {
		return b.ConfineInclusive(start), b.ConfineInclusive(end)
	} else {
		return b.Confine(start), b.Confine(end)
	}
}

func (b *Buf) RegionLines(start, end Loc) []string {
	if start.Row == end.Row {
		line := b.Line(start.Row)
		if line == "" {
			return []string{""}
		}
		rs := []rune(line)
		return []string{string(rs[start.Col:end.Col])}
	}
	lines := []string{}
	line := b.Line(start.Row)
	if line == "" {
		lines = append(lines, "")
	} else {
		rs := []rune(line)
		lines = append(lines, string(rs[start.Col:]))
	}
	for row := start.Row + 1; row < end.Row; row++ {
		lines = append(lines, b.Line(row))
	}
	line = b.Line(end.Row)
	if line == "" {
		lines = append(lines, "")
	} else {
		rs := []rune(line)
		lines = append(lines, string(rs[:end.Col]))
	}
	return lines
}
