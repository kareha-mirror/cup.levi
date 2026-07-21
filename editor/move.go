package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/termi/rkind"
	"tea.kareha.org/cup/termi/rutil"

	"tea.kareha.org/cup/levi/internal/buf"
)

/////////////////////
// Motion Commands //
/////////////////////

// Note: Marking, Search, Character Finding Commands also have Motion Commands.

//
// Move by Character / Move by Line
//

// Move cursor left by character.
// Key: h
func (ed *Editor) MoveLeft(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveLeft: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Col -= n
	loc.Col = b.ConfineCol(loc)
	return loc, true
}

// Move cursor down by line.
// Key: j
func (ed *Editor) MoveDown(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveDown: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Row += n
	if !b.CheckRowInclusive(loc.Row) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	return loc, true
}

// Move cursor here.
// internal use
func (ed *Editor) MoveHere(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveHere: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Row += n - 1
	if !b.CheckRowInclusive(loc.Row) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	return loc, true
}

// Move cursor up by line.
// Key: k
func (ed *Editor) MoveUp(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveUp: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Row -= n
	if !b.CheckRowInclusive(loc.Row) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	return loc, true
}

// Move cursor right by character.
// Key: l
func (ed *Editor) MoveRight(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveRight: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Col += n
	loc.Col = b.ConfineCol(loc)
	return loc, true
}

//
// Move in Line
//

// Move cursor to start of current line.
// Key: 0
func (ed *Editor) MoveToStart() (buf.Loc, bool) {
	loc := ed.Buf().Loc
	loc.Col = 0
	return loc, true
}

// Move cursor to end of current line.
// Key: $
func (ed *Editor) MoveToEnd(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveToEnd: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Row += n - 1
	if !b.CheckRowInclusive(loc.Row) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	loc.Col = utf8.RuneCountInString(b.Line(loc.Row))
	return loc, true
}

// Move cursor to first non-blank character of current line.
// Key: ^
func (ed *Editor) MoveToAfterIndent() (buf.Loc, bool) {
	b := ed.Buf()
	loc := b.Loc
	loc.Col = b.NonBlankColOfLine(loc.Row)
	return loc, true
}

// Move cursor to column <num> of current line.
// Key: <num>|
// Proper vi's column number is visual-based, but levi's is rune-based.
func (ed *Editor) MoveToColumn(n int) (buf.Loc, bool) { // n: 1-based
	if n < 1 {
		ed.Error("MoveToColumn: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Col = n - 1
	loc.Col = b.ConfineCol(loc)
	return loc, true
}

//
// Move by Word / Move by Loose Word
//

// Move cursor forward by word.
// Key: w
func (ed *Editor) MoveByWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByWord: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 0; i < n; i++ {
		if loc, found = b.MoveByWord(loc); found {
			continue
		}
		loc.Row++
		loc.Col = 0
		if loc, found = b.SkipBlanks(loc); !found {
			return loc, true
		}
	}
	return loc, true
}

// Move cursor forward by word used by cw.
// internal use
func (ed *Editor) MoveByChangeWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByChangeWord: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 1; i < n; i++ {
		if loc, found = b.MoveByWord(loc); found {
			continue
		}
		loc.Row++
		loc.Col = 0
		if loc, found = b.SkipBlanks(loc); !found {
			return loc, true
		}
	}
	if loc, found = b.MoveByWordAlt(loc); found {
		return loc, true
	}
	return loc, true
}

// Move cursor forward by word used by dw.
// internal use
func (ed *Editor) MoveByDeleteWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByDeleteWord: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 0; i < n; i++ {
		if loc, found = b.MoveByWord(loc); found {
			continue
		}
		if i == n-1 && b.Line(loc.Row) != "" {
			break
		}
		loc.Row++
		loc.Col = 0
		if loc, found = b.SkipBlanks(loc); !found {
			return loc, true
		}
	}
	return loc, true
}

// Move cursor backward by word.
// Key: b
func (ed *Editor) MoveBackwardByWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBackwardByWord: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 0; i < n; i++ {
		if loc.Col > 0 {
			loc.Col--
		} else {
			if loc.Row < 1 {
				return loc, true
			}
			loc.Row--
			line := b.Line(loc.Row)
			rc := utf8.RuneCountInString(line)
			loc.Col = max(rc-1, 0)
		}
		if loc, found = b.SkipBackwardBlanks(loc); !found {
			return loc, true
		}
		if loc, found = b.MoveBackwardByWord(loc); !found {
			return loc, true
		}
	}
	return loc, true
}

// Move cursor to end of word.
// Key: e
func (ed *Editor) MoveToEndOfWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveToEndOfWord: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 1; i < n; i++ {
		if loc, found = b.MoveByWord(loc); found {
			continue
		}
		loc.Row++
		loc.Col = 0
		if loc, found = b.SkipBlanks(loc); !found {
			return loc, true
		}
	}
	loc.Col++
	rc := utf8.RuneCountInString(b.Line(loc.Row))
	if loc.Col >= rc {
		if loc.Row < b.NumLines()-1 {
			loc.Col = 0
			loc.Row++
		} else {
			loc.Col = rc
			return loc, true
		}
	}
	if loc, found = b.SkipBlanks(loc); !found {
		return loc, true
	}
	if loc, found = b.MoveByWordAlt(loc); found {
		loc.Col = max(loc.Col-1, 0)
		return loc, true
	}
	return loc, true
}

// Move cursor forward by loose word.
// Key: W
func (ed *Editor) MoveByLooseWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByLooseWord: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 0; i < n; i++ {
		if loc, found = b.MoveByLooseWord(loc); found {
			continue
		}
		loc.Row++
		loc.Col = 0
		if loc, found = b.SkipBlanks(loc); !found {
			return loc, true
		}
	}
	return loc, true
}

// Move cursor forward by loose word used by cW.
// internal use
func (ed *Editor) MoveByChangeLooseWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByChangeLooseWord: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 1; i < n; i++ {
		if loc, found = b.MoveByLooseWord(loc); found {
			continue
		}
		loc.Row++
		loc.Col = 0
		if loc, found = b.SkipBlanks(loc); !found {
			return loc, true
		}
	}
	if loc, found = b.MoveByLooseWordAlt(loc); found {
		return loc, true
	}
	return loc, true
}

// Move cursor forward by word used by dW.
// internal use
func (ed *Editor) MoveByDeleteLooseWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByDeleteLooseWord: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 0; i < n; i++ {
		if loc, found = b.MoveByLooseWord(loc); found {
			continue
		}
		if i == n-1 && b.Line(loc.Row) != "" {
			break
		}
		loc.Row++
		loc.Col = 0
		if loc, found = b.SkipBlanks(loc); !found {
			return loc, true
		}
	}
	return loc, true
}

// Move cursor backward by loose word.
// Key: B
func (ed *Editor) MoveBackwardByLooseWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBackwardByLooseWord: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 0; i < n; i++ {
		if loc.Col > 0 {
			loc.Col--
		} else {
			if loc.Row < 1 {
				return loc, true
			}
			loc.Row--
			line := b.Line(loc.Row)
			rc := utf8.RuneCountInString(line)
			loc.Col = max(rc-1, 0)
		}
		if loc, found = b.SkipBackwardBlanks(loc); !found {
			return loc, true
		}
		if loc, found = b.MoveBackwardByLooseWord(loc); !found {
			return loc, true
		}
	}
	return loc, true
}

// Move cursor to end of loose word.
// Key: E
func (ed *Editor) MoveToEndOfLooseWord(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveToEndOfLooseWord: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	var found bool
	for i := 1; i < n; i++ {
		if loc, found = b.MoveByLooseWord(loc); found {
			continue
		}
		loc.Row++
		loc.Col = 0
		if loc, found = b.SkipBlanks(loc); !found {
			return loc, true
		}
	}
	loc.Col++
	rc := utf8.RuneCountInString(b.Line(loc.Row))
	if loc.Col >= rc {
		if loc.Row < b.NumLines()-1 {
			loc.Col = 0
			loc.Row++
		} else {
			loc.Col = rc
			return loc, true
		}
	}
	if loc, found = b.SkipBlanks(loc); !found {
		return loc, true
	}
	if loc, found = b.MoveByLooseWordAlt(loc); found {
		loc.Col = max(loc.Col-1, 0)
		return loc, true
	}
	return loc, true
}

//
// Move by Line
//

// Move cursor to first non-blank character of next line.
// Key: Enter, +
func (ed *Editor) MoveByLine(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByLine: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Row += n
	if !b.CheckRowInclusive(loc.Row) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	loc.Col = b.NonBlankColOfLine(loc.Row)
	return loc, true
}

// Move cursor to first non-blank character of previous line.
// Key: -
func (ed *Editor) MoveBackwardByLine(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBackwardByLine: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	loc.Row -= n
	if !b.CheckRowInclusive(loc.Row) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	loc.Col = b.NonBlankColOfLine(loc.Row)
	return loc, true
}

// Move cursor to first non-blank character of last line.
// Key: G
func (ed *Editor) MoveToLastLine() (buf.Loc, bool) {
	var loc buf.Loc
	b := ed.Buf()
	loc.Row = b.ConfineRow(b.NumLines() - 1)
	loc.Col = b.NonBlankColOfLine(loc.Row)
	return loc, true
}

// Move cursor to first non-blank character of line specified by <num>.
// Key: <num>G
func (ed *Editor) MoveToLine(n int) (buf.Loc, bool) { // n: 1-based
	if n < 1 {
		ed.Error("MoveToLine: n < 1")
		return buf.Loc{}, false
	}
	var loc buf.Loc
	loc.Row = n - 1
	b := ed.Buf()
	if !b.CheckRowInclusive(loc.Row) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	loc.Col = b.NonBlankColOfLine(loc.Row)
	return loc, true
}

//
// Move by Block
//

func (ed *Editor) moveBySentence(loc buf.Loc) buf.Loc {
	b := ed.Buf()
	first := true
	for loc.Row < b.NumLines() {
		line := b.Line(loc.Row)
		if first {
			if line == "" || rkind.IsBlank(rutil.RuneAt(line, loc.Col)) {
				loc, _ = b.SkipBlanks(loc)
				return loc
			}
		} else {
			if rkind.IsBlankLine(line) {
				return buf.Loc{Col: 0, Row: loc.Row}
			}
		}
		first = false

		col := 0
		found := false
	loop:
		for _, r := range line {
			if col < loc.Col {
				col++
				continue
			}
			if found {
				col++
				switch r {
				case ' ', '\t':
					break loop
				case ')', ']', '}', '"', '\'':
					continue
				default:
					found = false
					continue
				}
			}
			found = r == '.' || r == '?' || r == '!'
			col++
		}
		if found {
			loc, _ = b.SkipBlanks(buf.Loc{Col: col, Row: loc.Row})
			return loc
		}
		loc.Col = 0
		loc.Row++
	}
	return loc
}

// Move cursor forward by sentence.
// Key: )
func (ed *Editor) MoveBySentence(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBySentence: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	for i := 0; i < n; i++ {
		loc = ed.moveBySentence(loc)
	}
	return loc, true
}

func (ed *Editor) moveBackwardBySentence(loc buf.Loc) buf.Loc {
	b := ed.Buf()
	first := true
	line := b.Line(loc.Row)
	head := false
	var headLoc buf.Loc
	for {
		if first {
			if line == "" || rkind.IsBlank(rutil.RuneAt(line, loc.Col)) {
				loc, _ = b.SkipBackwardBlanks(loc)
			}
		} else {
			if rkind.IsBlankLine(line) {
				if head {
					return headLoc
				}
				return buf.Loc{Col: 0, Row: loc.Row}
			}
		}
		first = false
		col := 0
		found := false
		orig := b.Loc
		nbCol := b.NonBlankColOfLine(loc.Row)
		list := []int{nbCol}
		for _, r := range line {
			if col >= loc.Col {
				break
			}
			if found {
				col++
				switch r {
				case ' ', '\t':
					list = append(list, col-2)
					found = false
					continue
				case ')', ']', '}', '"', '\'':
					continue
				default:
					found = false
					continue
				}
			}
			if r == '.' || r == '?' || r == '!' {
				found = true
			}
			col++
		}
		if found {
			list = append(list, col-1)
		}
		if found && head {
			return headLoc
		}
		for i := len(list) - 1; i >= 0; i-- {
			f := list[i]
			l := buf.Loc{Col: f, Row: loc.Row}
			if f > nbCol {
				l = ed.moveBySentence(l)
			}
			if l != orig {
				if f <= nbCol {
					head = true
					headLoc = l
				} else {
					return l
				}
			}
		}
		loc.Row--
		if loc.Row < 0 {
			break
		}
		line = b.Line(loc.Row)
		loc.Col = utf8.RuneCountInString(line)
	}
	if loc.Row < 0 {
		loc.Row = 0
	}
	col := ed.Buf().NonBlankColOfLine(loc.Row)
	return buf.Loc{Col: col, Row: loc.Row}
}

// Move cursor backward by sentence.
// Key: (
func (ed *Editor) MoveBackwardBySentence(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBackwardBySentence: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	for i := 0; i < n; i++ {
		loc = ed.moveBackwardBySentence(loc)
	}
	return loc, true
}

func (ed *Editor) moveByParagraph(loc buf.Loc) buf.Loc {
	b := ed.Buf()
	if loc.Row >= b.NumLines()-1 {
		loc.Col = max(utf8.RuneCountInString(b.Line(loc.Row))-1, 0)
		return loc
	}
	loc, ok := b.SkipBlanks(loc)
	if !ok {
		return loc
	}
	loc.Col = 0
	for ; loc.Row < b.NumLines(); loc.Row++ {
		line := b.Line(loc.Row)
		if rkind.IsBlankLine(line) {
			return loc
		}
	}
	if loc.Row >= b.NumLines() {
		loc.Row = max(b.NumLines()-1, 0)
		loc.Col = 0
	}
	return loc
}

// Move cursor forward by paragraph.
// Key: }
func (ed *Editor) MoveByParagraph(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveByParagraph: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	for i := 0; i < n; i++ {
		loc = ed.moveByParagraph(loc)
	}
	return loc, true
}

func (ed *Editor) moveBackwardByParagraph(loc buf.Loc) buf.Loc {
	b := ed.Buf()
	loc, ok := b.SkipBackwardBlanks(loc)
	if !ok {
		return loc
	}
	loc.Col = 0
	for ; loc.Row >= 0; loc.Row-- {
		line := b.Line(loc.Row)
		if rkind.IsBlankLine(line) {
			return loc
		}
	}
	if loc.Row < 0 {
		loc.Row = 0
		loc.Col = 0
	}
	return loc
}

// Move cursor backward by paragraph.
// Key: {
func (ed *Editor) MoveBackwardByParagraph(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBackwardByParagraph: n < 1")
		return buf.Loc{}, false
	}
	b := ed.Buf()
	loc := b.Loc
	for i := 0; i < n; i++ {
		loc = ed.moveBackwardByParagraph(loc)
	}
	return loc, true
}

// Move cursor forward by section.
// Key: ]]
func (ed *Editor) MoveBySection(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBySection: n < 1")
		return buf.Loc{}, false
	}
	ed.Unimplemented("MoveBySection")
	return buf.Loc{}, false
}

// Move cursor backward by section.
// Key: [[
func (ed *Editor) MoveBackwardBySection(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveBackwardBySection: n < 1")
		return buf.Loc{}, false
	}
	ed.Unimplemented("MoveBackwardBySection")
	return buf.Loc{}, false
}

//
// Move in View
//

// Move cursor to top of view.
// Key: H
func (ed *Editor) MoveToTopOfView() (buf.Loc, bool) {
	if len(ed.viewMeta) < 1 {
		return buf.Loc{}, false
	}
	loc := ed.viewMeta[0].Loc
	if loc.Col < 1 {
		loc.Col = ed.Buf().NonBlankColOfLine(loc.Row)
	}
	return loc, true
}

// Move cursor to middle of view.
// Key: M
func (ed *Editor) MoveToMiddleOfView() (buf.Loc, bool) {
	if len(ed.viewMeta) < 1 {
		return buf.Loc{}, false
	}
	i := len(ed.viewMeta)/2 - 1
	loc := ed.viewMeta[i].Loc
	if loc.Col < 1 {
		loc.Col = ed.Buf().NonBlankColOfLine(loc.Row)
	}
	return loc, true
}

// Move cursor to bottom of view.
// Key: L
func (ed *Editor) MoveToBottomOfView() (buf.Loc, bool) {
	if len(ed.viewMeta) < 1 {
		return buf.Loc{}, false
	}
	i := len(ed.viewMeta) - 1
	loc := ed.viewMeta[i].Loc
	if loc.Col < 1 {
		loc.Col = ed.Buf().NonBlankColOfLine(loc.Row)
	}
	return loc, true
}

// Move cursor below <num> lines from top of view.
// Key: <num>H
func (ed *Editor) MoveToBelowTopOfView(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveToBelowTopOfView: n < 1")
		return buf.Loc{}, false
	}
	i := n - 1
	if i >= len(ed.viewMeta) {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	loc := ed.viewMeta[i].Loc
	if loc.Col < 1 {
		loc.Col = ed.Buf().NonBlankColOfLine(loc.Row)
	}
	return loc, true
}

// Move cursor above <num> lines from bottom of view.
// Key: <num>L
func (ed *Editor) MoveToAboveBottomOfView(n int) (buf.Loc, bool) {
	if n < 1 {
		ed.Error("MoveToAboveBottomOfView: n < 1")
		return buf.Loc{}, false
	}
	i := len(ed.viewMeta) - n
	if i < 0 {
		ed.Notice("Out of range")
		return buf.Loc{}, false
	}
	loc := ed.viewMeta[i].Loc
	if loc.Col < 1 {
		loc.Col = ed.Buf().NonBlankColOfLine(loc.Row)
	}
	return loc, true
}
