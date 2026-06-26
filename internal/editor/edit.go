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
func (ed *Editor) EditReplace(letter rune, n int, replay bool) {
	ed.Commit()
	if n < 1 {
		ed.Error("EditReplace: n < 1")
		return
	}
	ed.Unimplemented("EditReplace")
}

// J : Join current line with next line.
func (ed *Editor) EditJoin(n int) {
	ed.Commit()
	if n < 1 {
		ed.Error("EditJoin: n < 1")
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
func (ed *Editor) EditIndent(n int) {
	ed.Commit()
	if n < 1 {
		ed.Error("EditIndent: n < 1")
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
func (ed *Editor) EditOutdent(n int) {
	ed.Commit()
	if n < 1 {
		ed.Error("EditOutdent: n < 1")
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
func (ed *Editor) EditIndentRegion(start buf.Loc, end buf.Loc) {
	ed.Commit()
	start, end = ed.confineRegion(start, end, true)
	b := ed.Buf()
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
func (ed *Editor) EditOutdentRegion(start buf.Loc, end buf.Loc) {
	ed.Commit()
	start, end = ed.confineRegion(start, end, true)
	b := ed.Buf()
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
