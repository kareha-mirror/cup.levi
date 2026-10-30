package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
)

//////////////////////
// Editing Commands //
//////////////////////

// Replace single character under cursor.
// Key: r<char>
func (ed *Editor) Replace(r rune, n int) bool {
	if n < 1 {
		ed.Error("Replace: n < 1")
		return false
	}
	if !ed.Buf().Replace(r, n) {
		ed.Notice("Out of range")
		return false
	}
	return true
}

// Join current line with next line.
// Key: J
func (ed *Editor) Join(n int) bool {
	if n < 1 {
		ed.Error("Join: n < 1")
		return false
	}
	if !ed.Buf().Join(n) {
		ed.Ring("No following lines to join")
		return false
	}
	return true
}

// Indent region from current cursor to destination of motion <mv>.
// Key: > <mv>
func (ed *Editor) IndentRegion(start, end buf.Loc) bool {
	indent := "\t"
	b := ed.Buf()
	if ed.cfg.DetectIndent {
		indent = b.Indent
	}
	if !b.IndentRegion(start, end, indent) {
		ed.Notice("No lines indented")
		return false
	}
	return true
}

// Outdent region from current cursor to destination of motion <mv>.
// Key: < <mv>
func (ed *Editor) OutdentRegion(start, end buf.Loc) bool {
	indent := "\t"
	b := ed.Buf()
	if ed.cfg.DetectIndent {
		indent = b.Indent
	}
	if !b.OutdentRegion(start, end, indent) {
		ed.Notice("No lines outdented")
		return false
	}
	return true
}
