package editor

import (
	"strings"
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/rkind"
	"tea.kareha.org/cup/levi/internal/rutil"
)

//////////////////////
// Editing Commands //
//////////////////////

// r : Replace single character under cursor.
func (ed *Editor) Replace(char rune, n int) bool {
	if n < 1 {
		ed.Error("Replace: n < 1")
		return false
	}
	b := ed.Buf()
	line := b.CurrentLine()
	rc := utf8.RuneCountInString(line)
	if b.Loc.Col+n > rc {
		ed.Notice("Out of range")
		return false
	}

	head, _, tail := rutil.SplitBody(line, b.Loc.Col, b.Loc.Col+n)
	body := make([]rune, n)
	for i := 0; i < n; i++ {
		body[i] = char
	}
	b.SetCurrentLine(head + string(body) + tail)

	b.Loc.Col += n - 1
	b.VirtCol = b.Loc.Col
	return true
}

// J : Join current line with next line.
func (ed *Editor) Join(n int) bool {
	if n < 1 {
		ed.Error("Join: n < 1")
		return false
	}
	b := ed.Buf()
	if b.Loc.Row+1 >= b.NumLines() {
		ed.Ring("No following lines to join")
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

// >> : Indent current line.
func (ed *Editor) Indent(n int) bool {
	if n < 1 {
		ed.Error("Indent: n < 1")
		return false
	}
	b := ed.Buf()
	if b.Loc.Row+n > b.NumLines() {
		ed.Notice("Out of range")
		return false
	}
	for row := b.Loc.Row; row < b.Loc.Row+n; row++ {
		line := b.Line(row)
		b.SetLine(row, "\t"+line)
	}
	b.Loc.Col++
	b.Loc = b.ConfineInclusive(b.Loc)
	return true
}

// << : Outdent current line.
func (ed *Editor) Outdent(n int) bool {
	if n < 1 {
		ed.Error("Outdent: n < 1")
		return false
	}
	b := ed.Buf()
	if b.Loc.Row+n > b.NumLines() {
		ed.Notice("Out of range")
		return false
	}
	outdented := false
	for row := b.Loc.Row; row < b.Loc.Row+n; row++ {
		line := b.Line(row)
		if strings.HasPrefix(line, "\t") {
			b.SetLine(row, line[1:])
			if row == b.Loc.Row {
				outdented = true
			}
		}
	}
	if outdented {
		b.Loc.Col--
		b.Loc = b.ConfineInclusive(b.Loc)
	}
	return true
}

// > <mv> : Indent region from current cursor to destination of motion <mv>.
func (ed *Editor) IndentRegion(start buf.Loc, end buf.Loc) bool {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, true)
	indented := false
	for row := start.Row; row <= end.Row; row++ {
		line := b.Line(row)
		b.SetLine(row, "\t"+line)
		if row == b.Loc.Row {
			indented = true
		}
	}
	if indented {
		b.Loc.Col++
	}
	b.Loc = b.ConfineInclusive(b.Loc)
	return true
}

// < <mv> : Outdent region from current cursor to destination of motion <mv>.
func (ed *Editor) OutdentRegion(start buf.Loc, end buf.Loc) bool {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, true)
	outdented := false
	for row := start.Row; row <= end.Row; row++ {
		line := b.Line(row)
		if strings.HasPrefix(line, "\t") {
			b.SetLine(row, line[1:])
			if row == b.Loc.Row {
				outdented = true
			}
		}
	}
	if outdented {
		b.Loc.Col--
		b.Loc = b.ConfineInclusive(b.Loc)
	}
	return true
}
