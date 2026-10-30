package buf

import (
	"strings"

	"tea.kareha.org/cup/termi/rutil"
)

func (b *Buf) CopyRegion(start, end Loc, inclusive bool) []string {
	start, end = b.ConfineRegion(start, end, inclusive, false)
	copied := b.RunewiseContent(start, end)

	b.Loc = start

	return copied
}

func (b *Buf) CopyLineRegion(start, end Loc) []string {
	start, end = b.ConfineRegion(start, end, true, true)
	if end.Row+1 > b.NumLines() {
		return nil
	}
	copied := append([]string{}, b.Lines[start.Row:end.Row+1]...)

	b.Loc = start

	return copied
}

func (b *Buf) Paste(copied []string, n int) {
	if len(copied) < 2 {
		line := b.CurrentLine()
		sb := strings.Builder{}
		head, tail := rutil.Split(line, b.Loc.Col+1)
		sb.WriteString(head)
		for i := 0; i < n; i++ {
			sb.WriteString(copied[0])
		}
		sb.WriteString(tail)
		b.SetCurrentLine(sb.String())

		if len(line) > 0 {
			b.Loc.Col++
			b.VirtCol = b.Loc.Col
		}
	} else {
		lines := []string{}
		lines = append(lines, b.Lines[:b.Loc.Row]...)

		head, tail := rutil.Split(b.CurrentLine(), b.Loc.Col+1)
		lines = append(lines, head+copied[0])
		if len(copied) > 2 {
			lines = append(lines, copied[1:len(copied)-1]...)
		}
		lines = append(lines, copied[len(copied)-1]+tail)

		if b.Loc.Row+1 <= b.NumLines()-1 {
			lines = append(lines, b.Lines[b.Loc.Row+1:]...)
		}

		b.Lines = lines
	}
}

func (b *Buf) PasteLines(copied []string, n int) {
	move := b.NumLines() > 0

	lines := []string{}
	if b.Loc.Row+1 <= b.NumLines() {
		lines = append(lines, b.Lines[:b.Loc.Row+1]...)
	}
	for i := 0; i < n; i++ {
		lines = append(lines, copied...)
	}
	if b.Loc.Row+1 <= b.NumLines()-1 {
		lines = append(lines, b.Lines[b.Loc.Row+1:]...)
	}
	b.Lines = lines

	if move {
		b.Loc.Row++
		b.Loc = b.Confine(b.Loc)
		b.Loc.Col = b.FirstNonBlankCol(b.Loc.Row)
		b.VirtCol = b.Loc.Col
	}
}

func (b *Buf) PasteBefore(copied []string, n int) {
	if len(copied) < 2 {
		sb := strings.Builder{}
		head, tail := rutil.Split(b.CurrentLine(), b.Loc.Col)
		sb.WriteString(head)
		for i := 0; i < n; i++ {
			sb.WriteString(copied[0])
		}
		sb.WriteString(tail)
		b.SetCurrentLine(sb.String())
	} else {
		lines := append([]string{}, b.Lines[:b.Loc.Row]...)

		head, tail := rutil.Split(b.CurrentLine(), b.Loc.Col)
		lines = append(lines, head+copied[0])

		if len(copied) > 2 {
			lines = append(
				lines, copied[1:len(copied)-1]...,
			)
		}

		lines = append(lines, copied[len(copied)-1]+tail)

		if b.Loc.Row+1 < b.NumLines() {
			lines = append(lines, b.Lines[b.Loc.Row+1:]...)
		}

		b.Lines = lines
	}
}

func (b *Buf) PasteLinesBefore(copied []string, n int) {
	lines := append([]string{}, b.Lines[:b.Loc.Row]...)
	for i := 0; i < n; i++ {
		lines = append(lines, copied...)
	}
	lines = append(lines, b.Lines[b.Loc.Row:]...)
	b.Lines = lines

	b.Loc.Col = b.FirstNonBlankCol(b.Loc.Row)
	b.VirtCol = b.Loc.Col
}

func (b *Buf) Delete(n int) string {
	line := b.CurrentLine()
	head, body, tail := rutil.SplitBody(line, b.Loc.Col, b.Loc.Col+n)
	b.SetCurrentLine(head + tail)

	b.Loc = b.ConfineInclusive(b.Loc)

	return body
}

func (b *Buf) DeleteBefore(n int) string {
	if n > b.Loc.Col {
		n = b.Loc.Col
	}
	b.Loc.Col -= n

	return b.Delete(n)
}

func (b *Buf) DeleteRegion(start, end Loc, inclusive bool) []string {
	start, end = b.ConfineRegion(start, end, inclusive, false)
	if inclusive {
		end.Col++
	}
	deleted := b.RunewiseContent(start, end)

	lines := append([]string{}, b.Lines[:start.Row]...)
	head := rutil.Head(b.Line(start.Row), start.Col)
	tail := rutil.Tail(b.Line(end.Row), end.Col)
	lines = append(lines, head+tail)
	if end.Row+1 < b.NumLines() {
		lines = append(lines, b.Lines[end.Row+1:]...)
	}
	b.Lines = lines

	b.Loc = start
	b.Loc = b.ConfineInclusive(b.Loc)

	return deleted
}

func (b *Buf) DeleteLineRegion(start, end Loc) []string {
	start, end = b.ConfineRegion(start, end, true, true)
	if end.Row+1 > b.NumLines() {
		return nil
	}
	deleted := append([]string{}, b.Lines[start.Row:end.Row+1]...)

	lines := append([]string{}, b.Lines[:start.Row]...)
	if end.Row+1 < b.NumLines() {
		lines = append(lines, b.Lines[end.Row+1:]...)
	}
	b.Lines = lines

	b.Loc = start
	b.Loc = b.ConfineInclusive(b.Loc)

	return deleted
}
