package buf

type snaps struct {
	list [][]string
	idx  int
	temp []string
	undo bool
	redo bool
}

func (b *Buf) numSnaps() int {
	return len(b.ss.list)
}

func (b *Buf) BeginSnapshot() {
	if b.Depth < 1 {
		return
	}

	b.ss.temp = append([]string{}, b.Lines...)
}

func (b *Buf) EndSnapshot() {
	if b.Depth < 1 {
		return
	}

	delta := 0
	if b.ss.undo {
		b.ss.idx++
		delta++
	}
	if b.ss.redo {
		b.ss.idx--
		delta++
	}

	if b.ss.idx+1-delta <= b.numSnaps() {
		b.ss.list = b.ss.list[:b.ss.idx+1-delta]
	}

	b.ss.list = append(b.ss.list, b.ss.temp)
	b.ss.temp = nil
	b.ss.idx = b.numSnaps() - 1

	if b.numSnaps() > b.Depth+1 {
		b.ss.list = b.ss.list[1:]
		b.ss.idx = b.numSnaps() - 1
	}

	b.ss.undo = false
	b.ss.redo = false
}

func (b *Buf) CancelSnapshot() {
	if b.Depth < 1 {
		return
	}

	b.ss.temp = nil
}

func (b *Buf) Undo() bool {
	if b.numSnaps() < 1 {
		return false
	}

	if b.ss.redo {
		b.ss.idx -= 2
		b.ss.redo = false
	}
	if b.ss.idx < 0 {
		return false
	}
	if b.ss.idx > b.numSnaps()-1 {
		b.ss.idx = b.numSnaps() - 1
		return false
	}

	if b.ss.idx >= b.numSnaps()-1 {
		b.BeginSnapshot()
		b.EndSnapshot()
		b.ss.idx = b.numSnaps() - 2
	}

	lines := append([]string{}, b.ss.list[b.ss.idx]...)
	b.Lines = lines
	b.ss.idx--
	b.ss.undo = true
	return true
}

func (b *Buf) Redo() bool {
	if b.numSnaps() < 1 {
		return false
	}

	if b.ss.undo {
		b.ss.idx += 2
		b.ss.undo = false
	}
	if b.ss.idx > b.numSnaps()-1 {
		return false
	}
	if b.ss.idx < 0 {
		b.ss.idx = 0
	}

	lines := append([]string{}, b.ss.list[b.ss.idx]...)
	b.Lines = lines
	b.ss.idx++
	b.ss.redo = true
	return true
}
