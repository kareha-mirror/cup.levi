package buf

type History struct {
	snapshots [][]string
	idx       int
	undo      bool
}

func (b *Buf) NumSnapshots() int {
	return len(b.History.snapshots)
}

func (b *Buf) TakeSnapshot() {
	if b.NumSnapshots() > 0 {
		b.History.snapshots = b.History.snapshots[:b.History.idx+1]
	}

	lines := append([]string{}, b.Lines...)
	b.History.snapshots = append(b.History.snapshots, lines)
	b.History.idx = b.NumSnapshots() - 1

	if b.NumSnapshots() > b.Depth {
		b.History.snapshots = b.History.snapshots[1:]
		b.History.idx = b.NumSnapshots() - 1
	}

	b.History.undo = false
}

func (b *Buf) CancelSnapshot() {
	if b.NumSnapshots() < 1 {
		return
	}

	b.History.snapshots = b.History.snapshots[:b.NumSnapshots()-1]
	b.History.idx = max(b.NumSnapshots()-1, 0)
}

func (b *Buf) Undo() bool {
	if b.NumSnapshots() < 1 {
		return false
	}
	if !b.History.undo && b.History.idx >= b.NumSnapshots()-1 {
		b.TakeSnapshot()
		if b.NumSnapshots() < 1 {
			return false
		}
	}
	if b.History.idx < 1 {
		return false
	}
	b.History.idx--

	lines := append([]string{}, b.History.snapshots[b.History.idx]...)
	b.Lines = lines

	b.History.undo = true
	return true
}

func (b *Buf) Redo() bool {
	if b.History.idx >= b.NumSnapshots()-1 {
		return false
	}
	b.History.idx++

	lines := append([]string{}, b.History.snapshots[b.History.idx]...)
	b.Lines = lines
	return true
}
