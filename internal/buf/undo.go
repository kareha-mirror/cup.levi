package buf

type History struct {
	snapshots [][]string
	idx       int
	undo      bool
	redo      bool
}

func (b *Buf) NumSnapshots() int {
	return len(b.History.snapshots)
}

func (b *Buf) TakeSnapshot() {
	if b.Depth < 1 {
		return
	}

	delta := 0
	if b.History.undo {
		b.History.idx++
		delta++
	}
	if b.History.redo {
		b.History.idx--
		delta++
	}

	if b.History.idx+1-delta <= b.NumSnapshots() {
		b.History.snapshots = b.History.snapshots[:b.History.idx+1-delta]
	}

	lines := append([]string{}, b.Lines...)
	b.History.snapshots = append(b.History.snapshots, lines)
	b.History.idx = b.NumSnapshots() - 1

	if b.NumSnapshots() > b.Depth+1 {
		b.History.snapshots = b.History.snapshots[1:]
		b.History.idx = b.NumSnapshots() - 1
	}

	b.History.undo = false
	b.History.redo = false
}

func (b *Buf) CancelSnapshot() {
	if b.Depth < 1 {
		return
	}

	if b.NumSnapshots() < 1 {
		return
	}
	b.History.snapshots = b.History.snapshots[:b.NumSnapshots()-1]
	b.History.idx = b.NumSnapshots() - 1
}

func (b *Buf) Undo() bool {
	if b.NumSnapshots() < 1 {
		return false
	}

	if b.History.redo {
		b.History.idx -= 2
		b.History.redo = false
	}
	if b.History.idx < 0 {
		return false
	}
	if b.History.idx > b.NumSnapshots()-1 {
		b.History.idx = b.NumSnapshots() - 1
		return false
	}

	if b.History.idx >= b.NumSnapshots()-1 {
		b.TakeSnapshot()
		b.History.idx = b.NumSnapshots() - 2
	}

	lines := append([]string{}, b.History.snapshots[b.History.idx]...)
	b.Lines = lines
	b.History.idx--
	b.History.undo = true
	return true
}

func (b *Buf) Redo() bool {
	if b.NumSnapshots() < 1 {
		return false
	}

	if b.History.undo {
		b.History.idx += 2
		b.History.undo = false
	}
	if b.History.idx > b.NumSnapshots()-1 {
		return false
	}
	if b.History.idx < 0 {
		b.History.idx = 0
	}

	lines := append([]string{}, b.History.snapshots[b.History.idx]...)
	b.Lines = lines
	b.History.idx++
	b.History.redo = true
	return true
}
