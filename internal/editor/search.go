package editor

import (
	"regexp"
	"unicode/utf8"

	"tea.kareha.org/cup/termi"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/rutil"
)

type searchState struct {
	backward bool
	pattern  termi.RuneBuf
	regexp   *regexp.Regexp
}

func (ed *Editor) Locate() {
	b := ed.Buf()
	if len(ed.viewMeta) < 1 {
		return
	}
	minRow := ed.viewMeta[0].Loc.Row
	maxRow := ed.viewMeta[len(ed.viewMeta)-1].Loc.Row
	if b.Loc.Row >= minRow && b.Loc.Row <= maxRow {
		// XXX col is not checked
		return
	}
	viewRow := b.Loc.Row - (ed.h-1)/2 + 1
	if viewRow < 0 {
		viewRow = 0
	}
	b.ViewLoc.Row = viewRow
}

/////////////////////
// Search Commands //
/////////////////////

// /<pattern> Enter : Search <pattern> and move to it.
//func (ed *Editor) Search(pattern string) (buf.Loc, bool) {
func (ed *Editor) Search() (buf.Loc, bool) { // XXX
	if ed.searchs.regexp == nil {
		ed.Ring("No previous search pattern")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	for row := b.Loc.Row; row < b.NumLines(); row++ {
		line := b.Line(row)
		if row == b.Loc.Row {
			line = rutil.Tail(line, b.Loc.Col+1)
		}
		loc := ed.searchs.regexp.FindStringIndex(line)
		if loc == nil {
			continue
		}
		col := utf8.RuneCountInString(line[:loc[0]])
		if row == b.Loc.Row {
			col += b.Loc.Col + 1
		}
		return buf.Loc{col, row}, true
	}
	for row := 0; row <= b.Loc.Row; row++ {
		line := b.Line(row)
		loc := ed.searchs.regexp.FindStringIndex(line)
		if loc == nil {
			continue
		}
		col := utf8.RuneCountInString(line[:loc[0]])
		ed.Ring("Search wrapped")
		return buf.Loc{col, row}, true
	}
	ed.Ring("Pattern not found")
	return buf.Loc{}, false
}

// ?<pattern> Enter : Search <pattern> backward and move to it.
//func (ed *Editor) SearchBackward(pattern string) (buf.Loc, bool) {
func (ed *Editor) SearchBackward() (buf.Loc, bool) { // XXX
	if ed.searchs.regexp == nil {
		ed.Ring("No previous search pattern")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	end := len(rutil.Head(b.CurrentLine(), b.Loc.Col))
	for row := b.Loc.Row; row >= 0; row-- {
		line := b.Line(row)
		subLine := line
		var found []int
		for {
			loc := ed.searchs.regexp.FindStringIndex(subLine)
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
		return buf.Loc{col, row}, true
	}
	for row := b.NumLines() - 1; row >= b.Loc.Row; row-- {
		line := b.Line(row)
		subLine := line
		var found []int
		for {
			loc := ed.searchs.regexp.FindStringIndex(subLine)
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
		ed.Ring("Search wrapped")
		return buf.Loc{col, row}, true
	}
	ed.Ring("Pattern not found")
	return buf.Loc{}, false
}

// n : Repeat last search operation to search next match.
func (ed *Editor) SearchNext() (buf.Loc, bool) {
	if ed.searchs.backward {
		return ed.RepeatBackwardSearch()
	} else {
		return ed.RepeatSearch()
	}
}

// N : Repeat last search operation to search previous match.
func (ed *Editor) SearchPrev() (buf.Loc, bool) {
	if ed.searchs.backward {
		return ed.RepeatSearch()
	} else {
		return ed.RepeatBackwardSearch()
	}
}

// / Enter : Repeat last search.
func (ed *Editor) RepeatSearch() (buf.Loc, bool) {
	//return ed.Search(ed.searchs.pattern)
	return ed.Search() // XXX
}

// ? Enter : Repeat last backward search.
func (ed *Editor) RepeatBackwardSearch() (buf.Loc, bool) {
	//return ed.SearchBackward(ed.searchs.pattern)
	return ed.SearchBackward() // XXX
}
