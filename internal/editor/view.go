package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
)

///////////////////
// View Commands //
///////////////////

//
// Scroll by View Height / Scroll by Line
//

// Ctrl-F : Scroll down by view height.
func (ed *Editor) ViewDown(n int) {
	ed.EnsureCommand()
	i := len(ed.vMeta) - 2
	if i < 0 {
		return
	}
	b := ed.Buf()
	b.Loc = ed.vMeta[i].Loc
	b.ViewLoc = b.Loc
	if b.Loc.Col < 1 {
		ed.toNonBlankCol()
	}
}

// Ctrl-B : Scroll up by view height.
func (ed *Editor) ViewUp(n int) {
	ed.EnsureCommand()
	b := ed.Buf()
	viewRow := b.ViewLoc.Row - (ed.h - 3)
	if viewRow < 0 {
		viewRow = 0
	}
	viewMeta := ed.RenderMeta(buf.Loc{0, viewRow})
	if len(viewMeta) < 1 {
		return
	}
	lastRow := viewMeta[len(viewMeta)-1].Loc.Row

	if len(ed.vMeta) < 1 {
		return
	}
	topRow := ed.vMeta[0].Loc.Row

	deltaRow := topRow - lastRow - 1
	if deltaRow < 0 {
		deltaRow = 0
	}
	if deltaRow >= len(viewMeta) {
		return
	}
	newViewLoc := viewMeta[deltaRow].Loc
	viewMeta = ed.RenderMeta(newViewLoc)
	if len(viewMeta) < 2 {
		return
	}
	b.ViewLoc = newViewLoc
	b.Loc = viewMeta[len(viewMeta)-2].Loc
	if b.Loc.Col < 1 {
		ed.toNonBlankCol()
	}
	// TODO buggy
}

// Ctrl-D : Scroll down by half view height.
func (ed *Editor) ViewDownHalf(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("ViewDownHalf")
}

// Ctrl-U : Scroll up by half view height.
func (ed *Editor) ViewUpHalf(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("ViewUpHalf")
}

// Ctrl-Y : Scroll down by line.
func (ed *Editor) ViewDownLine(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("ViewDownLine")
}

// Ctrl-E : Scroll up by line.
func (ed *Editor) ViewUpLine(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("ViewUpLine")
}

//
// Reposition
//

// z Enter : Reposition cursor line to top of view.
func (ed *Editor) ViewToTop() {
	ed.EnsureCommand()
	ed.Unimplemented("ViewToTop")
}

// z. : Reposition cursor line middle of view.
func (ed *Editor) ViewToMiddle() {
	ed.EnsureCommand()
	ed.Unimplemented("ViewToMiddle")
}

// z- : Reposition cursor line bottom of view.
func (ed *Editor) ViewToBottom() {
	ed.EnsureCommand()
	ed.Unimplemented("ViewToBottom")
}

//
// Redraw
//

// Ctrl-L : Redraw view.
func (ed *Editor) ViewRedraw() {
	ed.EnsureCommand()
	ed.redraw = true
}
