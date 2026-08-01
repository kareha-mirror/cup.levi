package buf

type snaps struct {
	list [][]string
	idx  int
	temp []string
	undo bool
	redo bool
}

func (b *Buf) numSnaps() int {
	return len(b.snaps.list)
}

func (b *Buf) BeginSnapshot() {
	if b.Depth < 1 {
		return
	}

	b.snaps.temp = append([]string{}, b.Lines...)
}

func (b *Buf) EndSnapshot() {
	if b.Depth < 1 {
		return
	}

	delta := 0
	if b.snaps.undo {
		b.snaps.idx++
		delta++
	}
	if b.snaps.redo {
		b.snaps.idx--
		delta++
	}

	if b.snaps.idx+1-delta <= b.numSnaps() {
		b.snaps.list = b.snaps.list[:b.snaps.idx+1-delta]
	}

	b.snaps.list = append(b.snaps.list, b.snaps.temp)
	b.snaps.temp = nil
	b.snaps.idx = b.numSnaps() - 1

	if b.numSnaps() > b.Depth+1 {
		b.snaps.list = b.snaps.list[1:]
		b.snaps.idx = b.numSnaps() - 1
	}

	b.snaps.undo = false
	b.snaps.redo = false
}

func (b *Buf) CancelSnapshot() {
	if b.Depth < 1 {
		return
	}

	b.snaps.temp = nil
}

func (b *Buf) Undo() bool {
	if b.numSnaps() < 1 {
		return false
	}

	if b.snaps.redo {
		b.snaps.idx -= 2
		b.snaps.redo = false
	}
	if b.snaps.idx < 0 {
		return false
	}
	if b.snaps.idx > b.numSnaps()-1 {
		b.snaps.idx = b.numSnaps() - 1
		return false
	}

	if b.snaps.idx >= b.numSnaps()-1 {
		b.BeginSnapshot()
		b.EndSnapshot()
		b.snaps.idx = b.numSnaps() - 2
	}

	lines := append([]string{}, b.snaps.list[b.snaps.idx]...)
	b.Lines = lines
	b.snaps.idx--
	b.snaps.undo = true
	return true
}

func (b *Buf) Redo() bool {
	if b.numSnaps() < 1 {
		return false
	}

	if b.snaps.undo {
		b.snaps.idx += 2
		b.snaps.undo = false
	}
	if b.snaps.idx > b.numSnaps()-1 {
		return false
	}
	if b.snaps.idx < 0 {
		b.snaps.idx = 0
	}

	lines := append([]string{}, b.snaps.list[b.snaps.idx]...)
	b.Lines = lines
	b.snaps.idx++
	b.snaps.redo = true
	return true
}
