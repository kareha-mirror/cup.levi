package buffer

import (
	"unicode/utf8"
)

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
