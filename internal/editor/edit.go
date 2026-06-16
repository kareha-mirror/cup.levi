package editor

import (
	"unicode/utf8"
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
		if !isBlank(r) {
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
	b := ed.Buffer()
	if b.row+1 >= len(b.lines) {
		ed.Ring("No following lines to join")
		return
	}
	if n > 1 {
		n--
	}

	current := b.lines[b.row]
	col := b.col

	for i := 1; i <= n; i++ {
		if b.row+i >= len(b.lines) {
			break
		}
		next := trimLeftBlanks(b.lines[b.row+i])
		link := ""
		if len(next) > 0 {
			r, _ := utf8.DecodeLastRuneInString(current)
			if r != utf8.RuneError && !isBlank(r) {
				link = " "
			}
		}
		col = utf8.RuneCountInString(current)
		current = current + link + next
	}

	lines := append([]string{}, b.lines[:b.row]...)
	lines = append(lines, current)
	if b.row+1+n < len(b.lines) {
		lines = append(lines, b.lines[b.row+1+n:]...)
	}
	b.lines = lines

	b.col = col
	ed.confineCol()
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
func (ed *Editor) EditIndentRegion(start Loc, end Loc) {
	ed.EnsureCommand()
	ed.Unimplemented("EditIndentRegion")
}

// < <mv> : Outdent region from current cursor to destination of motion <mv>.
func (ed *Editor) EditOutdentRegion(start Loc, end Loc) {
	ed.EnsureCommand()
	ed.Unimplemented("EditOutdentRegion")
}
