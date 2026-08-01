package buf

func (b *Buf) StoreLine() {
	b.stored = b.CurrentLine()
}

func (b *Buf) RestoreLine() bool {
	if b.CurrentLine() == b.stored {
		return false
	}
	b.SetCurrentLine(b.stored)
	b.Loc.Col = 0
	return true
}
