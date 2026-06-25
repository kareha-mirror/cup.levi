package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/rkind"
)

//////////////////////
// Editing Commands //
//////////////////////

// r : Replace single character under cursor.
func (ed *Editor) EditReplace(letter rune, n int, replay bool) {
	ed.EnsureCommand()
	ed.Unimplemented("EditReplace")
}

func trimLeftBlanks(s string) string {
	for i, r := range s {
		if !rkind.IsBlank(r) {
			return s[i:]
		}
	}
	return ""
}

// J : Join current line with next line.
func (ed *Editor) EditJoin(n int) {
	if n < 1 {
		ed.Error("EditJoin: n < 1")
		return
	}
	ed.EnsureCommand()
	b := ed.Buf()
	if b.Loc.Row+1 >= b.NumLines() {
		ed.Ring("No following lines to join")
		return
	}
	if n > 1 {
		n--
	}

	current := b.CurrentLine()
	col := b.Loc.Col

	for i := 1; i <= n; i++ {
		if b.Loc.Row+i >= b.NumLines() {
			break
		}
		next := trimLeftBlanks(b.Line(b.Loc.Row + i))
		link := ""
		if len(next) > 0 {
			r, _ := utf8.DecodeLastRuneInString(current)
			if r != utf8.RuneError && !rkind.IsBlank(r) {
				link = " "
			}
		}
		col = utf8.RuneCountInString(current)
		current = current + link + next
	}

	lines := append([]string{}, b.Lines[:b.Loc.Row]...)
	lines = append(lines, current)
	if b.Loc.Row+1+n < b.NumLines() {
		lines = append(lines, b.Lines[b.Loc.Row+1+n:]...)
	}
	b.Lines = lines

	b.Loc.Col = col
	b.Loc.Col = b.ConfineCol(b.Loc)
}

// >> : Indent current line.
func (ed *Editor) EditIndent(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("EditIndent")
}

// << : Outdent current line.
func (ed *Editor) EditOutdent(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("EditOutdent")
}

// > <mv> : Indent region from current cursor to destination of motion <mv>.
func (ed *Editor) EditIndentRegion(start buf.Loc, end buf.Loc) {
	ed.EnsureCommand()
	ed.Unimplemented("EditIndentRegion")
}

// < <mv> : Outdent region from current cursor to destination of motion <mv>.
func (ed *Editor) EditOutdentRegion(start buf.Loc, end buf.Loc) {
	ed.EnsureCommand()
	ed.Unimplemented("EditOutdentRegion")
}
