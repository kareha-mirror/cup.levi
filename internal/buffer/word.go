package buffer

import (
	"unicode/utf8"
)

func isBlankRune(r rune) bool {
	return r == ' ' || r == '\t'
}

func (b *Buffer) SkipBlankLines() bool {
	for b.Loc.Row < b.NumLines() {
		line := b.CurrentLine()
		col := b.Loc.Col
		i := 0
		for _, r := range line {
			if i >= col && !isBlankRune(r) {
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

func (b *Buffer) SkipBackwardBlankLines() bool {
	for b.Loc.Row >= 0 {
		line := b.CurrentLine()
		col := b.Loc.Col
		rs := []rune(line)
		i := len(rs) - 1
		for ; i >= 0; i-- {
			r := rs[i]
			if i <= col && !isBlankRune(r) {
				b.Loc.Col = i
				return true
			}
		}
		if b.Loc.Row < 1 {
			b.Loc.Col = max(i, 0)
			return false
		}
		b.Loc.Row--
		line = b.CurrentLine()
		b.Loc.Col = max(utf8.RuneCountInString(line)-1, 0)
	}
	b.Loc.Row = 0
	b.Loc.Col = 0
	return false
}

func isWordRune(r rune) bool {
	return (r >= '0' && r <= '9') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z')
}

func isSymbolRune(r rune) bool {
	return !isBlankRune(r) && !isWordRune(r) && r < 0x100
}

func (b *Buffer) MoveByWord(advance bool) bool {
	line := b.CurrentLine()
	if len(line) < 1 {
		return false
	}
	rs := []rune(line)
	col := b.Loc.Col
	if advance {
		col++
	}
	if advance && isWordRune(rs[b.Loc.Col]) {
		for ; col < len(rs); col++ {
			r := rs[col]
			if isSymbolRune(r) {
				b.Loc.Col = col
				return true
			}
			if !isWordRune(r) {
				break
			}
		}
		if col >= len(rs) {
			b.Loc.Col = len(rs) - 1
			return false
		}
		for ; col < len(rs); col++ {
			r := rs[col]
			if isWordRune(r) || isSymbolRune(r) {
				b.Loc.Col = col
				return true
			}
		}
		b.Loc.Col = len(rs) - 1
		return false
	} else if advance && isSymbolRune(rs[b.Loc.Col]) {
		for ; col < len(rs); col++ {
			r := rs[col]
			if isWordRune(r) {
				b.Loc.Col = col
				return true
			}
			if !isSymbolRune(r) {
				break
			}
		}
		if col >= len(rs) {
			b.Loc.Col = len(rs) - 1
			return false
		}
		for ; col < len(rs); col++ {
			r := rs[col]
			if isWordRune(r) || isSymbolRune(r) {
				b.Loc.Col = col
				return true
			}
		}
		b.Loc.Col = len(rs) - 1
		return false
	} else {
		for ; col < len(rs); col++ {
			r := rs[col]
			if isWordRune(r) || isSymbolRune(r) {
				b.Loc.Col = col
				return true
			}
		}
		b.Loc.Col = len(rs) - 1
		return false
	}
}

func (b *Buffer) MoveBackwardByWord() bool {
	line := b.CurrentLine()
	if len(line) < 1 {
		return false
	}
	rs := []rune(line)
	col := b.Loc.Col

	for ; col >= 0; col-- {
		r := rs[col]
		if isWordRune(r) || isSymbolRune(r) {
			b.Loc.Col = col
			break
		}
	}
	if col < 0 {
		b.Loc.Col = 0
		return false
	}

	if isWordRune(rs[col]) {
		if col < 1 {
			return false
		}
		col--
		for ; col >= 0; col-- {
			r := rs[col]
			if !isWordRune(r) {
				break
			}
		}
		b.Loc.Col = col + 1
		return true
	} else {
		if col < 1 {
			return false
		}
		col--
		for ; col >= 0; col-- {
			r := rs[col]
			if !isSymbolRune(r) {
				break
			}
		}
		b.Loc.Col = col + 1
		return true
	}
}
