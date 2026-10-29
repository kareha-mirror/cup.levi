package buf

import "tea.kareha.org/cup/termi/rutil"

func (b *Buf) Delete(n int) string {
	line := b.CurrentLine()
	head, body, tail := rutil.SplitBody(line, b.Loc.Col, b.Loc.Col+n)
	b.SetCurrentLine(head + tail)
	b.Loc = b.ConfineInclusive(b.Loc)
	return body
}
