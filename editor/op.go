package editor

import (
	"strings"
	"unicode/utf8"

	"tea.kareha.org/cup/termi/rutil"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/kill"
)

//////////////////////////////////////
// Operator Commands (Copy / Delte) //
//////////////////////////////////////

//
// Copy (Yank)
//

// y<mv> : Copy region from current cursor to destination of motion <mv>.
func (ed *Editor) CopyRegion(
	reg rune, start buf.Loc, end buf.Loc, inclusive bool,
) {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, inclusive, false)
	lines := b.RegionRunewise(start, end)
	ed.StoreRunes(reg, lines)
	b.Loc = start
}

// y<mv> : Copy region from current cursor to destination of motion <mv>.
func (ed *Editor) CopyLineRegion(
	reg rune, start buf.Loc, end buf.Loc,
) {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, true, true)
	if end.Row+1 > b.NumLines() {
		ed.Notice("Out of range")
		return
	}
	ed.StoreLines(reg, b.Lines[start.Row:end.Row+1])
	b.Loc = start
}

//
// Paste (Put)
//

// "<reg>p : Paste after cursor from register <reg>.
func (ed *Editor) Paste(reg rune, n int) bool {
	if n < 1 {
		ed.Error("Paste: n < 1")
		return false
	}
	if ed.KillMode(reg) == kill.None {
		if reg == 0 {
			ed.Ring("The default buffer is empty")
		} else {
			ed.Ring("Buffer %c is empty", reg)
		}
		return false
	}
	killed := ed.KilledContent(reg)
	b := ed.Buf()
	switch ed.KillMode(reg) {
	case kill.Runes:
		if len(killed) < 2 {
			line := b.CurrentLine()
			sb := strings.Builder{}
			head, tail := rutil.Split(line, b.Loc.Col+1)
			sb.WriteString(head)
			for i := 0; i < n; i++ {
				sb.WriteString(killed[0])
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
			lines = append(lines, head+killed[0])
			if len(killed) > 2 {
				lines = append(lines, killed[1:len(killed)-1]...)
			}
			lines = append(lines, killed[len(killed)-1]+tail)

			if b.Loc.Row+1 <= b.NumLines()-1 {
				lines = append(lines, b.Lines[b.Loc.Row+1:]...)
			}

			b.Lines = lines
		}
	case kill.Lines:
		lines := []string{}
		if b.Loc.Row+1 <= b.NumLines() {
			lines = append(lines, b.Lines[:b.Loc.Row+1]...)
		}
		for i := 0; i < n; i++ {
			lines = append(lines, killed...)
		}
		if b.Loc.Row+1 <= b.NumLines()-1 {
			lines = append(lines, b.Lines[b.Loc.Row+1:]...)
		}
		move := b.NumLines() > 0
		b.Lines = lines
		if move {
			b.Loc.Row++
			b.Loc = b.Confine(b.Loc)
			b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
			b.VirtCol = b.Loc.Col
		}
	}
	return true
}

// "<reg>p : Paste before cursor from register <reg>.
func (ed *Editor) PasteBefore(reg rune, n int) bool {
	if n < 1 {
		ed.Error("PasteBefore: n < 1")
		return false
	}
	b := ed.Buf()
	if ed.KillMode(reg) == kill.None {
		if reg == 0 {
			ed.Ring("The default buffer is empty")
		} else {
			ed.Ring("Buffer %c is empty", reg)
		}
		return false
	}
	killed := ed.KilledContent(reg)
	switch ed.KillMode(reg) {
	case kill.Runes:
		if len(killed) < 2 {
			sb := strings.Builder{}
			head, tail := rutil.Split(b.CurrentLine(), b.Loc.Col)
			sb.WriteString(head)
			for i := 0; i < n; i++ {
				sb.WriteString(killed[0])
			}
			sb.WriteString(tail)
			b.SetCurrentLine(sb.String())
		} else {
			lines := append([]string{}, b.Lines[:b.Loc.Row]...)

			head, tail := rutil.Split(b.CurrentLine(), b.Loc.Col)
			lines = append(lines, head+killed[0])

			if len(killed) > 2 {
				lines = append(
					lines, killed[1:len(killed)-1]...,
				)
			}

			lines = append(lines, killed[len(killed)-1]+tail)

			if b.Loc.Row+1 < b.NumLines() {
				lines = append(lines, b.Lines[b.Loc.Row+1:]...)
			}

			b.Lines = lines
		}
	case kill.Lines:
		lines := append([]string{}, b.Lines[:b.Loc.Row]...)
		for i := 0; i < n; i++ {
			lines = append(lines, killed...)
		}
		lines = append(lines, b.Lines[b.Loc.Row:]...)
		b.Lines = lines
		b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
		b.VirtCol = b.Loc.Col
	}
	return true
}

//
// Delete
//

func (ed *Editor) internalDelete(reg rune, n int) bool {
	b := ed.Buf()
	if len(b.CurrentLine()) < 1 {
		return false
	}
	line := b.CurrentLine()
	rc := utf8.RuneCountInString(line)
	n = min(n, rc-b.Loc.Col)
	head, body, tail := rutil.SplitBody(line, b.Loc.Col, b.Loc.Col+n)
	ed.StoreRunes(reg, []string{body})
	b.SetCurrentLine(head + tail)
	return true
}

// x : Delete character under cursor.
func (ed *Editor) Delete(reg rune, n int) bool {
	if n < 1 {
		ed.Error("Delete: n < 1")
		return false
	}
	if !ed.internalDelete(reg, n) {
		ed.Notice("Nothing to delete")
		return false
	}
	b := ed.Buf()
	b.Loc = b.ConfineInclusive(b.Loc)
	return true
}

// X : Delete character before cursor.
func (ed *Editor) DeleteBefore(reg rune, n int) bool {
	if n < 1 {
		ed.Error("DeleteBefore: n < 1")
		return false
	}
	b := ed.Buf()
	if n > b.Loc.Col {
		n = b.Loc.Col
	}
	if n < 1 {
		ed.Notice("Nothing to delete")
		return false
	}
	b.Loc.Col -= n
	if !ed.internalDelete(reg, n) {
		ed.Notice("Nothing to delete")
		return false
	}
	b.Loc = b.ConfineInclusive(b.Loc)
	return true
}

// d<mv> : Delete region from current cursor to destination of motion <mv>.
func (ed *Editor) DeleteRegion(
	reg rune, start buf.Loc, end buf.Loc, inclusive bool,
) bool {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, inclusive, false)
	if inclusive {
		end.Col++
	}
	lines := b.RegionRunewise(start, end)
	ed.StoreRunes(reg, lines)

	lines = append([]string{}, b.Lines[:start.Row]...)

	head := rutil.Head(b.Line(start.Row), start.Col)
	tail := rutil.Tail(b.Line(end.Row), end.Col)
	lines = append(lines, head+tail)

	if end.Row+1 < b.NumLines() {
		lines = append(lines, b.Lines[end.Row+1:]...)
	}
	b.Lines = lines
	b.Loc = start
	b.Loc = b.ConfineInclusive(b.Loc)
	return true
}

// d<mv> : Delete region from current cursor to destination of motion <mv>.
func (ed *Editor) DeleteLineRegion(
	reg rune, start buf.Loc, end buf.Loc,
) bool {
	b := ed.Buf()
	start, end = b.ConfineRegion(start, end, true, true)
	if end.Row+1 > b.NumLines() {
		ed.Notice("Out of range")
		return false
	}
	lines := append([]string{}, b.Lines[:start.Row]...)
	ed.StoreLines(reg, b.Lines[start.Row:end.Row+1])
	if end.Row+1 < b.NumLines() {
		lines = append(lines, b.Lines[end.Row+1:]...)
	}
	b.Lines = lines
	b.Loc = start
	b.Loc = b.ConfineInclusive(b.Loc)
	return true
}
