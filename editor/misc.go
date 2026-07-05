package editor

import (
	"tea.kareha.org/cup/termi/suspend"

	"tea.kareha.org/cup/levi/internal/cmd"
)

////////////////////////////
// Miscellaneous Commands //
////////////////////////////

// Ctrl-G : Show info about states of current buffer.
func (ed *Editor) ShowInfo() {
	ed.Message("%s", ed.Buf().Info())
}

// Ctrl-L : Redraw view.
func (ed *Editor) Redraw() {
	ed.redraw = true
}

// . : Repeat last command which is repeatable.
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

// u : Undo last modification or redo by undoing itself.
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

// U : Restore current line to previous state last visited.
func (ed *Editor) Restore() bool {
	return ed.Buf().RestoreLine()
}

// ZZ : Save and close.
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

// Ctrl-Z : Suspend editor process.
func (ed *Editor) Suspend() {
	suspend.Suspend()
}
