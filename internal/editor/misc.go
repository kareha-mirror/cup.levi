package editor

import (
	"fmt"
)

////////////////////////////
// Miscellaneous Commands //
////////////////////////////

// Ctrl-g : Show info such as current cursor position.
func (ed *Editor) MiscShowInfo() {
	ed.EnsureCommand()
	path := ed.path
	if path == "" {
		path = "(memory)"
	}
	modified := "unmodified"
	if ed.modified {
		modified = "modified"
	}
	info := "empty file"
	linesLen := len(ed.lines)
	if linesLen > 0 {
		info = fmt.Sprintf(
			"line %d of %d [%d%%]",
			ed.row+1, linesLen, 100*(ed.row+1)/linesLen,
		)
	}
	ed.Message("%s: %s: %s", path, modified, info)
}

// . : Repeat last edit.
func (ed *Editor) MiscRepeat(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MiscRepeat")
}

// u : Undo.
func (ed *Editor) MiscUndo(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("MiscUndo")
}

// U : Restore current line to previous state.
func (ed *Editor) MiscRestore() {
	ed.EnsureCommand()
	ed.Unimplemented("MiscRestore")
}

// ZZ : Save and quit.
func (ed *Editor) MiscSaveAndQuit() {
	ed.EnsureCommand()
	if ed.modified && ed.path == "" {
		ed.Ring("File is a temporary; exit will discard modifications")
		return
	}
	if ed.modified && ed.path != "" {
		err := ed.Save(false)
		if err != nil {
			return
		}
	}
	ed.alive = false
}
