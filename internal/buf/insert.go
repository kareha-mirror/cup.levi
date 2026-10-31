package buf

import (
	"unicode/utf8"

	"tea.kareha.org/cup/termi/rkind"
	"tea.kareha.org/cup/termi/rutil"
)

func (b *Buf) internalInsert(inserted []string, ai bool) {
	ins := append([]string{}, inserted...)
	line := b.CurrentLine()
	head, tail := rutil.Split(line, b.Loc.Col)
	ins[0] = head + ins[0]
	ins[len(ins)-1] = ins[len(ins)-1] + tail
	if ai {
		for i := 0; i < len(ins); i++ {
			if rkind.IsBlankLine(ins[i]) {
				ins[i] = ""
			}
		}
	}
	lines := append([]string{}, b.Lines[:b.Loc.Row]...)
	lines = append(lines, ins...)
	if b.Loc.Row+1 < b.NumLines() {
		lines = append(lines, b.Lines[b.Loc.Row+1:]...)
	}
	b.Lines = lines
	if len(ins) < 2 {
		b.Loc.Col += utf8.RuneCountInString(inserted[0])
	} else {
		b.Loc.Row += len(inserted) - 1
		b.Loc.Col = utf8.RuneCountInString(inserted[len(inserted)-1])
	}
}

func (b *Buf) Insert(inserted []string, ai bool, n int) {
	for i := 0; i < n; i++ {
		b.internalInsert(inserted, ai)
	}
	b.Loc.Col--
	b.Loc = b.ConfineInclusive(b.Loc)
	b.VirtCol = b.Loc.Col
}
