package buf

import (
	"strings"
	"unicode/utf8"

	"tea.kareha.org/cup/termi/rkind"
	"tea.kareha.org/cup/termi/rutil"
)

func (b *Buf) Replace(r rune, n int) bool {
	line := b.CurrentLine()
	rc := utf8.RuneCountInString(line)
	if b.Loc.Col+n > rc {
		return false
	}

	head, _, tail := rutil.SplitBody(line, b.Loc.Col, b.Loc.Col+n)
	body := strings.Repeat(string(r), n)
	b.SetCurrentLine(head + body + tail)

	b.Loc.Col += n - 1
	b.VirtCol = b.Loc.Col
	return true
}

func (b *Buf) Join(n int) bool {
	if b.Loc.Row+1 >= b.NumLines() {
		return false
	}
	if n > 1 {
		n--
	}

	sb := strings.Builder{}
	sb.WriteString(b.CurrentLine())
	for i := 1; i <= n; i++ {
		if b.Loc.Row+i >= b.NumLines() {
			break
		}
		next := rkind.TrimPrefixBlanks(b.Line(b.Loc.Row + i))
		b.Loc.Col = utf8.RuneCountInString(sb.String())
		if len(next) > 0 {
			r, _ := utf8.DecodeLastRuneInString(sb.String())
			if r != utf8.RuneError && !rkind.IsBlank(r) {
				sb.WriteString(" ")
			}
			sb.WriteString(next)
		}
	}

	lines := append([]string{}, b.Lines[:b.Loc.Row]...)
	lines = append(lines, sb.String())
	if b.Loc.Row+1+n < b.NumLines() {
		lines = append(lines, b.Lines[b.Loc.Row+1+n:]...)
	}
	b.Lines = lines

	b.Loc = b.ConfineInclusive(b.Loc)
	return true
}

func (b *Buf) IndentRegion(start, end Loc, indent string) bool {
	start, end = b.ConfineRegion(start, end, true, true)
	indented := false
	for row := start.Row; row <= end.Row; row++ {
		line := b.Line(row)
		if line == "" {
			continue
		}
		b.SetLine(row, indent+line)
		if row == b.Loc.Row {
			indented = true
		}
	}
	if indented {
		b.Loc.Col += len(indent)
		b.Loc = b.ConfineInclusive(b.Loc)
	}
	return indented
}

func (b *Buf) OutdentRegion(start, end Loc, indent string) bool {
	start, end = b.ConfineRegion(start, end, true, true)
	outdented := false
	for row := start.Row; row <= end.Row; row++ {
		line := b.Line(row)
		if strings.HasPrefix(line, indent) {
			b.SetLine(row, line[len(indent):])
			if row == b.Loc.Row {
				outdented = true
			}
		}
	}
	if outdented {
		b.Loc.Col -= len(indent)
		b.Loc = b.ConfineInclusive(b.Loc)
	}
	return outdented
}
