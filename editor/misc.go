package editor

import (
	"tea.kareha.org/cup/termi/suspend"

	"tea.kareha.org/cup/levi/internal/cmd"
)

////////////////////////////
// Miscellaneous Commands //
////////////////////////////

// Show info about states of current buffer.
// Key: Ctrl-G
func (ed *Editor) ShowInfo() {
	ed.Message("%s", ed.Buf().Info())
}

// Redraw view.
// Key: Ctrl-L
func (ed *Editor) Redraw() {
	ed.redraw = true
}

// Repeat last command which is repeatable.
// Key: .
func (ed *Editor) Repeat(n int) {
	c := ed.lastCmd
	b := ed.Buf()
	prevRow := b.Loc.Row
	if _, ok := cmd.IsModifying[c.Op.Kind]; ok {
		ed.BeginUndoRecord()
	}
	if modified, ok := ed.Run(c, true); ok { // replay
		if modified {
			b.Modified = true
		}
		if ed.alive && ed.Buf() == b && b.Loc.Row != prevRow {
			b.StoreLine()
		}
		if _, ok := cmd.IsModifying[c.Op.Kind]; ok {
			if modified {
				ed.EndUndoRecord()
			} else {
				ed.CancelUndoRecord()
			}
		}
	} else {
		if _, ok := cmd.IsModifying[c.Op.Kind]; ok {
			ed.CancelUndoRecord()
		}
	}
}

// Undo last modification or redo by undoing itself.
// Key: u
func (ed *Editor) Undo(n int, replay bool) {
	if !replay {
		ed.undo = !ed.undo
	}
	b := ed.Buf()
	if ed.undo {
		if !b.Undo() {
			ed.Notice("No more history for undo")
			return
		}
	} else {
		if !b.Redo() {
			ed.Notice("No more history for redo")
			return
		}
	}
	b.Loc = b.ConfineInclusive(b.Loc)
	b.Modified = true
}

// Restore current line to previous state last visited.
// Key: U
func (ed *Editor) Restore() bool {
	return ed.Buf().RestoreLine()
}

// Save and close.
// Key: ZZ
func (ed *Editor) SaveAndClose() {
	b := ed.Buf()
	if b.Modified {
		if b.Path == "" {
			ed.Ring("File is a temporary; exit will discard modifications")
			return
		} else {
			if !ed.Save(false) {
				return
			}
		}
	}
	ed.Close(false)
	ed.CheckQuit()
}

// Suspend editor process.
// Key: Ctrl-Z
func (ed *Editor) Suspend() {
	suspend.Suspend()
}
