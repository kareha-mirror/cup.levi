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
	b := ed.Buffer()
	path := b.path
	if path == "" {
		path = "(memory)"
	}
	modified := "unmodified"
	if b.modified {
		modified = "modified"
	}
	info := "empty file"
	linesLen := len(b.lines)
	if linesLen > 0 {
		info = fmt.Sprintf(
			"line %d of %d [%d%%]",
			b.row+1, linesLen, 100*(b.row+1)/linesLen,
		)
	}
	ed.Message("%s: %s: %s", path, modified, info)
}

// . : Repeat last edit.
func (ed *Editor) MiscRepeat(n int) {
	ed.EnsureCommand()
	ed.Run(ed.lastCmd, true)
}

// u : Undo.
func (ed *Editor) MiscUndo(n int, replay bool) {
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
	b := ed.Buffer()
	if b.modified && b.path == "" {
		ed.Ring("File is a temporary; exit will discard modifications")
		return
	}
	if b.modified && b.path != "" {
		err := ed.Save(false)
		if err != nil {
			return
		}
	}
	ed.Close(false)
}
