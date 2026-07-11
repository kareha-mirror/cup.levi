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
	i := len(ed.viewMeta) - 2
	if i < 0 {
		return
	}
	b := ed.Buf()
	b.Loc = ed.viewMeta[i].Loc
	b.ViewLoc = b.Loc
	if b.Loc.Col < 1 {
		b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
	}
}

// Ctrl-B : Scroll up by view height.
func (ed *Editor) ViewUp(n int) {
	b := ed.Buf()
	viewRow := b.ViewLoc.Row - (ed.h - 3)
	if viewRow < 0 {
		viewRow = 0
	}
	viewMeta := ed.RenderMeta(buf.Loc{Col: 0, Row: viewRow})
	if len(viewMeta) < 1 {
		return
	}
	lastRow := viewMeta[len(viewMeta)-1].Loc.Row

	if len(ed.viewMeta) < 1 {
		return
	}
	topRow := ed.viewMeta[0].Loc.Row

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
		b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
	}
}

// Ctrl-D : Scroll down by half view height.
func (ed *Editor) ViewDownHalf(n int) {
	i := len(ed.viewMeta) / 2
	if i < 0 {
		return
	}
	b := ed.Buf()
	b.Loc = ed.viewMeta[i].Loc
	b.ViewLoc = b.Loc
	if b.Loc.Col < 1 {
		b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
	}
}

// Ctrl-U : Scroll up by half view height.
func (ed *Editor) ViewUpHalf(n int) {
	b := ed.Buf()
	viewRow := b.ViewLoc.Row - (ed.h-1)/2
	if viewRow < 0 {
		viewRow = 0
	}
	viewMeta := ed.RenderMeta(buf.Loc{Col: 0, Row: viewRow})
	if len(viewMeta) < 1 {
		return
	}
	lastRow := viewMeta[len(viewMeta)-1].Loc.Row

	if len(ed.viewMeta) < 1 {
		return
	}
	topRow := ed.viewMeta[0].Loc.Row

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
	b.Loc = viewMeta[len(viewMeta)-1].Loc
	if b.Loc.Col < 1 {
		b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
	}
}

// Ctrl-Y : Scroll down by line.
func (ed *Editor) ViewDownLine(n int) {
	ed.Commit()

	b := ed.Buf()
	if b.ViewLoc.Row < 1 {
		return
	}
	b.ViewLoc.Row--
	b.ViewLoc.Col = 0

	if len(ed.viewMeta) < 1 {
		return
	}
	meta := ed.viewMeta[len(ed.viewMeta)-1]
	if b.Loc.Row >= meta.Loc.Row {
		b.Loc.Row = max(b.Loc.Row-1, 0)
		b.Loc.Col = b.ConfineFreeColInclusive(b.Loc.Row)
	}
}

// Ctrl-E : Scroll up by line.
func (ed *Editor) ViewUpLine(n int) {
	ed.Commit()

	if len(ed.viewMeta) < 1 {
		return
	}
	meta := ed.viewMeta[len(ed.viewMeta)-1]
	b := ed.Buf()
	if meta.Loc.Row >= b.NumLines()-1 {
		return
	}
	b.ViewLoc.Row++
	b.ViewLoc.Col = 0

	if b.Loc.Row < b.ViewLoc.Row {
		b.Loc.Row = b.ViewLoc.Row
		b.Loc.Col = b.ConfineFreeColInclusive(b.Loc.Row)
	}
}

//
// Reposition
//

// z Enter : Reposition cursor line to top of view.
func (ed *Editor) ViewToTop() {
	b := ed.Buf()
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
	b.ViewLoc.Row = b.Loc.Row
	b.ViewLoc.Col = 0
}

// z. : Reposition cursor line middle of view.
func (ed *Editor) ViewToMiddle() {
	b := ed.Buf()
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
	b.ViewLoc.Row = max(b.Loc.Row-(ed.h-1)/2, 0)
	b.ViewLoc.Col = 0
}

// z- : Reposition cursor line bottom of view.
func (ed *Editor) ViewToBottom() {
	b := ed.Buf()
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
	b.ViewLoc.Row = max(b.Loc.Row-(ed.h-1)+1, 0)
	b.ViewLoc.Col = 0
}
