package editor

import (
	"fmt"

	"tea.kareha.org/cup/termi"
)

////////////////////////////
// Miscellaneous Commands //
////////////////////////////

// Ctrl-G : Show info such as current cursor position.
func (ed *Editor) MiscShowInfo() {
	ed.Commit()
	b := ed.Buf()
	path := b.Path
	if path == "" {
		path = "(memory)"
	}
	modified := "unmodified"
	if b.Modified {
		modified = "modified"
	}
	info := "empty file"
	numLines := b.NumLines()
	if numLines > 0 {
		info = fmt.Sprintf(
			"line %d of %d [%d%%]",
			b.Loc.Row+1, numLines, 100*(b.Loc.Row+1)/numLines,
		)
	}
	ed.Message("%s: %s: %s", path, modified, info)
}

// . : Repeat last edit.
func (ed *Editor) MiscRepeat(n int) {
	ed.Commit()
	c := ed.lastCmd
	if _, ok := InsertCmds[c.Kind]; ok {
		ed.BeginMemory()
	} else if _, ok := EditCmds[c.Kind]; ok {
		ed.BeginMemory()
	}
	if ed.Run(c, true) {
		if _, ok := InsertCmds[c.Kind]; ok {
			ed.EndMemory()
		} else if _, ok := EditCmds[c.Kind]; ok {
			ed.EndMemory()
		}
	} else {
		if _, ok := InsertCmds[c.Kind]; ok {
			ed.CancelMemory()
		} else if _, ok := EditCmds[c.Kind]; ok {
			ed.CancelMemory()
		}
	}
}

// u : Undo.
func (ed *Editor) MiscUndo(n int, replay bool) {
	ed.Commit()
	b := ed.Buf()
	if b.Snapshot == nil {
		return
	}
	b.Lines = b.Snapshot
	b.Snapshot = nil
	b.Loc = b.ConfineInclusive(b.Loc)
}

// U : Restore current line to previous state.
func (ed *Editor) MiscRestore() {
	ed.Commit()
	ed.Unimplemented("MiscRestore")
}

// ZZ : Save and quit.
func (ed *Editor) MiscSaveAndQuit() {
	ed.Commit()
	b := ed.Buf()
	if b.Modified && b.Path == "" {
		ed.Ring("File is a temporary; exit will discard modifications")
		return
	}
	if b.Modified && b.Path != "" {
		err := ed.Save(false)
		if err != nil {
			return
		}
	}
	ed.Close(false)
}

// Ctrl-Z : Suspend.
func (ed *Editor) MiscSuspend() {
	ed.Commit()
	termi.Suspend()
}
