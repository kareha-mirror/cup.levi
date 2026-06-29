package editor

import (
	"fmt"
	"unicode/utf8"

	"tea.kareha.org/cup/termi"
)

////////////////////////////
// Miscellaneous Commands //
////////////////////////////

// Ctrl-G : Show info such as current cursor position.
func (ed *Editor) ShowInfo() {
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
		numBytes := numLines
		numRunes := numLines
		if b.CRLF {
			numBytes *= 2
			numRunes *= 2
		}
		for _, line := range b.Lines {
			numBytes += len(line)
			numRunes += utf8.RuneCountInString(line)
		}
		info = fmt.Sprintf(
			"line %d of %d [%d%%] %d bytes, %d runes.",
			b.Loc.Row+1, numLines, 100*(b.Loc.Row+1)/numLines,
			numBytes, numRunes,
		)
	}
	ed.Message("%s: %s: %s", path, modified, info)
}

// . : Repeat last edit.
func (ed *Editor) Repeat(n int) {
	c := ed.lastCmd
	if _, ok := IsInsertCmd[c.Op.Kind]; ok {
		ed.BeginRecordForUndo()
	} else if _, ok := IsEditCmd[c.Op.Kind]; ok {
		ed.BeginRecordForUndo()
	}
	if ed.Run(c, true) {
		if _, ok := IsInsertCmd[c.Op.Kind]; ok {
			ed.EndRecordForUndo()
		} else if _, ok := IsEditCmd[c.Op.Kind]; ok {
			ed.EndRecordForUndo()
		}
	} else {
		if _, ok := IsInsertCmd[c.Op.Kind]; ok {
			ed.CancelRecordForUndo()
		} else if _, ok := IsEditCmd[c.Op.Kind]; ok {
			ed.CancelRecordForUndo()
		}
	}
}

// u : Undo.
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

// U : Restore current line to previous state.
func (ed *Editor) Restore() {
	ed.Unimplemented("Restore")
}

// ZZ : Save and quit.
func (ed *Editor) SaveAndQuit() {
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

// Ctrl-Z : Suspend.
func (ed *Editor) Suspend() {
	termi.Suspend()
}
