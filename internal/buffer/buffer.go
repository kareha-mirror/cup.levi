package buffer

import (
	"strings"
	"time"
	"unicode/utf8"
)

type Loc struct {
	Col int // 0-based
	Row int // 0-based
}

type Pos struct {
	X int // 0-based
	Y int // 0-based
}

type Stamp struct {
	Time time.Time
	Size int64
}

type Buffer struct {
	Loc      Loc
	ViewLoc  Loc
	VirtCol  int // 0-based
	Pos      Pos
	Lines    []string
	Path     string
	Modified bool
	Stamp    Stamp
	Marks    map[rune]Loc
}

func (b *Buffer) NumLines() int {
	return len(b.Lines)
}

func (b *Buffer) Line(row int) string {
	if len(b.Lines) < 1 {
		return ""
	}
	return b.Lines[row]
}

func (b *Buffer) SetLine(row int, line string) {
	if len(b.Lines) < 1 {
		b.Lines = []string{""}
	}
	b.Lines[row] = line
}

func (b *Buffer) CurrentLine() string {
	if len(b.Lines) < 1 {
		return ""
	}
	return b.Lines[b.Loc.Row]
}

func (b *Buffer) SetCurrentLine(line string) {
	b.SetLine(b.Loc.Row, line)
}

func (b *Buffer) Text() string {
	if len(b.Lines) < 1 {
		return ""
	}
	return strings.Join(b.Lines, "\n") + "\n"
}

func (b *Buffer) SetText(text string) {
	if len(text) < 1 {
		b.Lines = []string{}
	} else {
		// TODO should also support CRLF or not?
		if text[len(text)-1] == '\n' {
			text = text[:len(text)-1]
		}
		b.Lines = strings.Split(text, "\n")
	}
}

func (b *Buffer) MoveRow(row int) bool {
	if row < 0 {
		return false
	}
	numLines := len(b.Lines)
	if numLines == 0 && row == 0 {
		b.Loc.Row = row
		return true
	}
	if row >= numLines {
		return false
	}
	b.Loc.Row = row
	return true
}

func (b *Buffer) AdjustRow(n int) bool {
	return b.MoveRow(b.Loc.Row + n)
}

func (b *Buffer) ConfineRow() {
	n := len(b.Lines)
	if b.Loc.Row < 0 {
		b.Loc.Row = 0
	} else if b.Loc.Row >= n {
		b.Loc.Row = max(n-1, 0)
	}
}

func (b *Buffer) ConfineCol() {
	if b.Loc.Col < 0 {
		b.Loc.Col = 0
		return
	}
	rc := utf8.RuneCountInString(b.CurrentLine())
	if b.Loc.Col >= rc {
		b.Loc.Col = max(rc-1, 0)
	}
	if b.Loc.Col < b.ViewLoc.Col {
		b.ViewLoc.Col = 0
	}
}

func (b *Buffer) Confine() {
	b.ConfineRow()
	b.ConfineCol()
}

func (b *Buffer) SaveVirtCol() {
	b.VirtCol = b.Loc.Col
}

func (b *Buffer) LoadVirtCol() {
	b.Loc.Col = b.VirtCol
}

func (b *Buffer) MoveCol(col int) {
	b.Loc.Col = col
	b.ConfineCol()
	b.SaveVirtCol()
}

func (b *Buffer) AdjustCol(n int) {
	b.MoveCol(b.Loc.Col + n)
}

func (b *Buffer) Mark(r rune) {
	if b.Marks == nil {
		b.Marks = map[rune]Loc{}
	}
	b.Marks[r] = b.Loc
}

func isBlank(r rune) bool {
	return r == ' ' || r == '\t'
}

func (b *Buffer) SkipBlankLines() bool {
	for b.Loc.Row < b.NumLines() {
		line := b.CurrentLine()
		col := b.Loc.Col
		i := 0
		for _, r := range line {
			if i >= col && !isBlank(r) {
				b.Loc.Col = col
				return true
			}
			i++
		}
		if b.Loc.Row >= b.NumLines()-1 {
			b.Loc.Col = i
			return false
		}
		b.Loc.Row++
		b.Loc.Col = 0
	}
	return false
}

func isWordRune(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z')
}

func isSymbol(r rune) bool {
	return !isBlank(r) && !isWordRune(r) && r < 0x100
}

func (b *Buffer) MoveByWord(first bool) bool {
	line := b.CurrentLine()
	if len(line) < 1 {
		return false
	}
	rs := []rune(line)
	if first && isWordRune(rs[b.Loc.Col]) {
		col := b.Loc.Col + 1
		for ; col < len(rs); col++ {
			if isSymbol(rs[col]) {
				b.Loc.Col = col
				return true
			}
			if !isWordRune(rs[col]) {
				break
			}
		}
		if col >= len(rs) {
			b.Loc.Col = len(rs) - 1
			return false
		}
		for ; col < len(rs); col++ {
			if isWordRune(rs[col]) || isSymbol(rs[col]) {
				b.Loc.Col = col
				return true
			}
		}
		b.Loc.Col = len(rs) - 1
		return false
	} else {
		col := b.Loc.Col
		if first {
			col++
		}
		for ; col < len(rs); col++ {
			if isWordRune(rs[col]) || isSymbol(rs[col]) {
				b.Loc.Col = col
				return true
			}
		}
		b.Loc.Col = len(rs) - 1
		return false
	}
}
