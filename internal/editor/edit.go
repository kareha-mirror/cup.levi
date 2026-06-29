package editor

import (
	"strings"
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/rkind"
)

//////////////////////
// Editing Commands //
//////////////////////

// r : Replace single character under cursor.
func (ed *Editor) Replace(letter rune, n int, replay bool) {
	if n < 1 {
		ed.Error("Replace: n < 1")
		return
	}
	ed.Unimplemented("Replace")
}

// J : Join current line with next line.
func (ed *Editor) Join(n int) {
	if n < 1 {
		ed.Error("Join: n < 1")
		return
	}
	b := ed.Buf()
	if b.Loc.Row+1 >= b.NumLines() {
		ed.Ring("No following lines to join")
		return
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
	b.Modified = true
}

// >> : Indent current line.
func (ed *Editor) Indent(n int) {
	if n < 1 {
		ed.Error("Indent: n < 1")
		return
	}
	b := ed.Buf()
	if b.Loc.Row+n > b.NumLines() {
		ed.Notice("Out of range")
		return
	}
	for row := b.Loc.Row; row < b.Loc.Row+n; row++ {
		line := b.Line(row)
		b.SetLine(row, "\t"+line)
	}
	b.Loc.Col++
	b.Loc = b.ConfineInclusive(b.Loc)
	b.Modified = true
}

// << : Outdent current line.
func (ed *Editor) Outdent(n int) {
	if n < 1 {
		ed.Error("Outdent: n < 1")
		return
	}
	b := ed.Buf()
	if b.Loc.Row+n > b.NumLines() {
		ed.Notice("Out of range")
		return
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
	b.Modified = true
}

// > <mv> : Indent region from current cursor to destination of motion <mv>.
func (ed *Editor) IndentRegion(start buf.Loc, end buf.Loc) {
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
	b.Modified = true
}

// < <mv> : Outdent region from current cursor to destination of motion <mv>.
func (ed *Editor) OutdentRegion(start buf.Loc, end buf.Loc) {
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
	b.Modified = true
}
