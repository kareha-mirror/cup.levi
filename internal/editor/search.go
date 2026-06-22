package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buffer"
)

func (ed *Editor) Locate(loc buffer.Loc) {
	b := ed.Buffer()
	b.Loc = loc
	if len(ed.vMeta) < 1 {
		return
	}
	minRow := ed.vMeta[0].Loc.Row
	maxRow := ed.vMeta[len(ed.vMeta)-1].Loc.Row
	if loc.Row >= minRow && loc.Row <= maxRow {
		// XXX col is not checked
		return
	}
	viewRow := loc.Row - (ed.h-1)/2 + 1
	if viewRow < 0 {
		viewRow = 0
	}
	b.ViewLoc.Row = viewRow
}

/////////////////////
// Search Commands //
/////////////////////

// /<pattern> Enter : Search <pattern> forward.
func (ed *Editor) SearchForward() {
	ed.EnsureCommand()
	if ed.regexp == nil {
		ed.Ring("No previous search pattern")
		return
	}
	b := ed.Buffer()
	for row := b.Loc.Row; row < b.NumLines(); row++ {
		line := b.Line(row)
		if row == b.Loc.Row {
			rs := []rune(line)
			if b.Loc.Col+1 < len(rs) {
				line = string(rs[b.Loc.Col+1:])
			} else {
				line = ""
			}
		}
		loc := ed.regexp.FindStringIndex(line)
		if loc == nil {
			continue
		}
		col := utf8.RuneCountInString(line[:loc[0]])
		if row == b.Loc.Row {
			col += b.Loc.Col + 1
		}
		ed.Locate(buffer.Loc{col, row})
		return
	}
	ed.Ring("Search wrapped")
	for row := 0; row <= b.Loc.Row; row++ {
		line := b.Line(row)
		loc := ed.regexp.FindStringIndex(line)
		if loc == nil {
			continue
		}
		col := utf8.RuneCountInString(line[:loc[0]])
		ed.Locate(buffer.Loc{col, row})
		return
	}
	ed.Ring("Pattern not found")
}

// ?<pattern> Enter : Search <pattern> backward.
func (ed *Editor) SearchBackward() {
	ed.EnsureCommand()
	if ed.regexp == nil {
		ed.Ring("No previous search pattern")
		return
	}
	b := ed.Buffer()
	rs := []rune(b.CurrentLine())
	end := len(string(rs[:b.Loc.Col]))
	for row := b.Loc.Row; row >= 0; row-- {
		line := b.Line(row)
		subLine := line
		var found []int
		for {
			loc := ed.regexp.FindStringIndex(subLine)
			if loc == nil {
				break
			}
			if row == b.Loc.Row && loc[0] >= end {
				break
			}
			subLine = subLine[loc[1]:]
			end -= loc[1]
			if found == nil {
				found = loc
			} else {
				found[0] = found[1] + loc[0]
				found[1] += loc[1]
			}
		}
		if found == nil {
			continue
		}
		col := utf8.RuneCountInString(line[:found[0]])
		ed.Locate(buffer.Loc{col, row})
		return
	}
	ed.Ring("Search wrapped")
	for row := b.NumLines() - 1; row >= b.Loc.Row; row-- {
		line := b.Line(row)
		subLine := line
		var found []int
		for {
			loc := ed.regexp.FindStringIndex(subLine)
			if loc == nil {
				break
			}
			subLine = subLine[loc[1]:]
			if found == nil {
				found = loc
			} else {
				found[0] = found[1] + loc[0]
				found[1] += loc[1]
			}
		}
		if found == nil {
			continue
		}
		col := utf8.RuneCountInString(line[:found[0]])
		ed.Locate(buffer.Loc{col, row})
		return
	}
	ed.Ring("Pattern not found")
}

// n : Search next match.
func (ed *Editor) SearchNextMatch() {
	if ed.backward {
		ed.SearchRepeatBackward()
	} else {
		ed.SearchRepeatForward()
	}
}

// N : Search previous match.
func (ed *Editor) SearchPrevMatch() {
	if ed.backward {
		ed.SearchRepeatForward()
	} else {
		ed.SearchRepeatBackward()
	}
}

// / Enter : Repeat last search forward.
func (ed *Editor) SearchRepeatForward() {
	ed.SearchForward()
}

// ? Enter : Repeat last search backward.
func (ed *Editor) SearchRepeatBackward() {
	ed.SearchBackward()
}
