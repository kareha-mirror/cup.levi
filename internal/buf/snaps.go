package buf

import "maps"

type snapRecord struct {
	lines []string
	marks map[rune]Loc
}

type snaps struct {
	list []snapRecord
	idx  int
	temp snapRecord
}

func (b *Buf) numSnaps() int {
	return len(b.snaps.list)
}

func (b *Buf) BeginSnapshot() {
	if b.Depth < 1 {
		return
	}

	b.snaps.temp.lines = append([]string{}, b.Lines...)
	b.snaps.temp.marks = maps.Clone(b.Marks)
}

func (b *Buf) EndSnapshot() {
	if b.Depth < 1 {
		return
	}

	if b.snaps.idx+1 <= b.numSnaps() {
		b.snaps.list = b.snaps.list[:b.snaps.idx+1]
	}

	b.snaps.list = append(b.snaps.list, b.snaps.temp)
	b.snaps.temp = snapRecord{}
	b.snaps.idx = b.numSnaps() - 1

	if b.numSnaps() > b.Depth+1 {
		b.snaps.list = b.snaps.list[1:]
		b.snaps.idx = b.numSnaps() - 1
	}
}

func (b *Buf) CancelSnapshot() {
	if b.Depth < 1 {
		return
	}

	b.snaps.temp = snapRecord{}
}

func (b *Buf) Undo() bool {
	if b.numSnaps() < 1 {
		return false
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

	lines := append([]string{}, b.snaps.list[b.snaps.idx].lines...)
	b.Lines = lines
	b.Marks = maps.Clone(b.snaps.list[b.snaps.idx].marks)
	b.snaps.idx--
	return true
}

func (b *Buf) Redo() bool {
	if b.numSnaps() < 1 {
		return false
	}

	idx := b.snaps.idx + 2
	if idx > b.numSnaps()-1 {
		return false
	}

	lines := append([]string{}, b.snaps.list[idx].lines...)
	b.Lines = lines
	b.Marks = maps.Clone(b.snaps.list[idx].marks)
	b.snaps.idx++
	return true
}
