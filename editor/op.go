package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/regs"
)

//////////////////////////////////////
// Operator Commands (Copy / Delte) //
//////////////////////////////////////

//
// Copy (Yank)
//

// Copy region from current cursor to destination of motion <mv>.
// Key: y<mv>
func (ed *Editor) CopyRegion(reg rune, start, end buf.Loc, inclusive bool) {
	copied := ed.Buf().CopyRegion(start, end, inclusive)
	ed.StoreRunes(reg, copied)
}

// Copy region from current cursor to destination of motion <mv>.
// Key: y<mv>
func (ed *Editor) CopyLineRegion(reg rune, start, end buf.Loc) {
	copied := ed.Buf().CopyLineRegion(start, end)
	if copied == nil {
		ed.Notice("Out of range")
		return
	}
	ed.StoreLines(reg, copied)
}

//
// Paste (Put)
//

// Paste after cursor from register <reg>.
// Key: "<reg>p
func (ed *Editor) Paste(reg rune, n int) bool {
	if n < 1 {
		ed.Error("Paste: n < 1")
		return false
	}
	if ed.CopyMode(reg) == regs.None {
		if reg == 0 {
			ed.Ring("The default buffer is empty")
		} else {
			ed.Ring("Buffer %c is empty", reg)
		}
		return false
	}
	copied := ed.CopiedContent(reg)
	switch ed.CopyMode(reg) {
	case regs.Runes:
		ed.Buf().Paste(copied, n)
	case regs.Lines:
		ed.Buf().PasteLines(copied, n)
	}
	return true
}

// Paste before cursor from register <reg>.
// Key: "<reg>P
func (ed *Editor) PasteBefore(reg rune, n int) bool {
	if n < 1 {
		ed.Error("PasteBefore: n < 1")
		return false
	}
	if ed.CopyMode(reg) == regs.None {
		if reg == 0 {
			ed.Ring("The default buffer is empty")
		} else {
			ed.Ring("Buffer %c is empty", reg)
		}
		return false
	}
	copied := ed.CopiedContent(reg)
	switch ed.CopyMode(reg) {
	case regs.Runes:
		ed.Buf().PasteBefore(copied, n)
	case regs.Lines:
		ed.Buf().PasteLinesBefore(copied, n)
	}
	return true
}

//
// Delete
//

// Delete character under cursor.
// Key: x
func (ed *Editor) Delete(reg rune, n int) bool {
	if n < 1 {
		ed.Error("Delete: n < 1")
		return false
	}
	deleted := ed.Buf().Delete(n, true)
	if deleted == "" {
		ed.Notice("Nothing to delete")
		return false
	}
	ed.StoreRunes(reg, []string{deleted})
	return true
}

// Delete character before cursor.
// Key: X
func (ed *Editor) DeleteBefore(reg rune, n int) bool {
	if n < 1 {
		ed.Error("DeleteBefore: n < 1")
		return false
	}
	deleted := ed.Buf().DeleteBefore(n)
	if deleted == "" {
		ed.Notice("Nothing to delete")
		return false
	}
	ed.StoreRunes(reg, []string{deleted})
	return true
}

// Delete region from current cursor to destination of motion <mv>.
// Key: d<mv>
func (ed *Editor) DeleteRegion(
	reg rune, start, end buf.Loc, inclusive bool,
) bool {
	deleted := ed.Buf().DeleteRegion(start, end, inclusive)
	ed.StoreRunes(reg, deleted)
	return true
}

// Delete region from current cursor to destination of motion <mv>.
// Key: d<mv>
func (ed *Editor) DeleteLineRegion(reg rune, start, end buf.Loc) bool {
	deleted := ed.Buf().DeleteLineRegion(start, end)
	if deleted == nil {
		ed.Notice("Out of range")
		return false
	}
	ed.StoreLines(reg, deleted)
	return true
}
