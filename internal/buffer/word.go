package buffer

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/runekind"
)

func (b *Buffer) SkipBlankLines() bool {
	for b.Loc.Row < b.NumLines() {
		line := b.CurrentLine()
		col := 0
		for _, r := range line {
			if col >= b.Loc.Col && !runekind.IsBlank(r) {
				b.Loc.Col = col
				return true
			}
			col++
		}
		b.Loc.Row++
		b.Loc.Col = 0
	}
	b.Loc.Row = max(b.NumLines()-1, 0)
	b.Loc.Col = max(utf8.RuneCountInString(b.CurrentLine())-1, 0)
	return false
}

func (b *Buffer) SkipBackwardBlankLines() bool {
	for b.Loc.Row >= 0 {
		line := b.CurrentLine()
		if line != "" {
			rs := []rune(line)
			col := b.Loc.Col
			for ; col >= 0; col-- {
				r := rs[col]
				if !runekind.IsBlank(r) {
					b.Loc.Col = col
					return true
				}
			}
		}
		b.Loc.Row--
		line = b.CurrentLine()
		b.Loc.Col = max(utf8.RuneCountInString(line)-1, 0)
	}
	b.Loc.Row = 0
	b.Loc.Col = 0
	return false
}

func (b *Buffer) MoveByWord() bool {
	line := b.CurrentLine()
	if len(line) < 1 {
		return false
	}
	rs := []rune(line)
	col := b.Loc.Col
	kind := runekind.Kind(rs[col])
	col++
	k := kind
	for ; col < len(rs); col++ {
		k = runekind.Kind(rs[col])
		if k != kind {
			break
		}
	}
	if col >= len(rs) {
		return false
	}
	if kind == runekind.Blank || k != runekind.Blank {
		b.Loc.Col = col
		return true
	}
	col++
	for ; col < len(rs); col++ {
		k = runekind.Kind(rs[col])
		if k != runekind.Blank {
			break
		}
	}
	if col < len(rs) {
		b.Loc.Col = col
		return true
	}
	return false
}

func (b *Buffer) MoveBackwardByWord() bool {
	line := b.CurrentLine()
	if len(line) < 1 {
		return false
	}
	rs := []rune(line)
	kind := runekind.Kind(rs[b.Loc.Col])
	if kind == runekind.Blank {
		return false
	}
	col := b.Loc.Col
	for ; col >= 0; col-- {
		k := runekind.Kind(rs[col])
		if k != kind {
			break
		}
	}
	b.Loc.Col = col + 1
	return true
}
