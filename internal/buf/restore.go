package buf

func (b *Buf) StoreLine() {
	b.current = b.CurrentLine()
}

func (b *Buf) RestoreLine() bool {
	if b.CurrentLine() == b.current {
		return false
	}
	b.SetCurrentLine(b.current)
	b.Loc.Col = 0
	return true
}
