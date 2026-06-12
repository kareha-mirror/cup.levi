package editor

import (
	"tea.kareha.org/cup/termi"
)

/////////////////////
// Prompt Commands //
/////////////////////

// :+<num> Enter : Move cursor to first non-blank character of next line.
func (ed *Editor) PromptMoveByLine(n int) {
	if n < 0 {
		ed.Error("PromptMoveByLine: n < 0")
		return
	}
	ed.EnsureCommand()
	if !ed.adjustRow(n) {
		ed.Ring("Illegal address: only %d lines in the file.", len(ed.lines))
		return
	}
	ed.toNonBlankCol()
}

// :-<num> Enter : Move cursor to first non-blank character of previous line.
func (ed *Editor) PromptMoveBackwardByLine(n int) {
	if n < 0 {
		ed.Error("PromptMoveBackwardByLine: n < 0")
		return
	}
	ed.EnsureCommand()
	if ed.row-n == -1 {
		n++
	}
	if !ed.adjustRow(-n) {
		ed.Ring("Reference to a line number less than 0.")
		return
	}
	ed.toNonBlankCol()
}

// :<num> Enter : Move cursor to first non-blank character of line specifined by <num>.
func (ed *Editor) PromptMoveToLine(n int) { // n: 1-based
	if n < 0 {
		ed.Error("PromptMoveToLine: n < 0")
		return
	}
	ed.EnsureCommand()
	if n == 0 {
		n = 1
	}
	if !ed.setRow(n - 1) {
		ed.Ring("Illegal address: only %d lines in the file.", len(ed.lines))
		return
	}
	ed.toNonBlankCol()
}

// :wq Enter : Save current file and quit.
func (ed *Editor) PromptSaveAndQuit() {
	ed.EnsureCommand()
	if ed.modified && ed.path == "" {
		ed.Ring("File is a temporary; exit will discard modifications.")
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

// :w Enter : Save current file.
func (ed *Editor) PromptSave(name string) {
	ed.EnsureCommand()
	if name == "" {
		ed.Save(false)
		return
	}
	ed.SaveAs(name, false)
}

// :w! Enter : Force save current file.
func (ed *Editor) PromptForceSave(name string) {
	ed.EnsureCommand()
	if name == "" {
		ed.Save(true)
		return
	}
	ed.SaveAs(name, true)
}

// :q Enter : Quit editor.
func (ed *Editor) PromptQuit() {
	ed.EnsureCommand()
	if ed.modified {
		if ed.path == "" {
			ed.Ring("File is a temporary; exit will discard modifications.")
			return
		}
		ed.Ring("File modified since last complete write; write or use ! to override.")
		return
	}
	ed.alive = false
}

// :q! Enter : Force quit editor.
func (ed *Editor) PromptForceQuit() {
	ed.EnsureCommand()
	ed.alive = false
}

// :e Enter : Open file.
func (ed *Editor) PromptOpen() {
	ed.EnsureCommand()
	ed.Unimplemented("PromptOpen")
}

// :e! Enter : Force open file.
func (ed *Editor) PromptForceOpen() {
	ed.EnsureCommand()
	ed.Unimplemented("PromptForceOpen")
}

// :r Enter : Read file and insert to current buffer.
func (ed *Editor) PromptRead() {
	ed.EnsureCommand()
	ed.Unimplemented("PromptRead")
}

// :n Enter : Switch to next buffer (tab).
func (ed *Editor) PromptNext() {
	ed.EnsureCommand()
	ed.Unimplemented("PromptNext")
}

// :prev Enter : Switch to previous buffer (tab).
func (ed *Editor) PromptPrev() {
	ed.EnsureCommand()
	ed.Unimplemented("PromptPrev")
}

// :sh Enter : Execute shell.
func (ed *Editor) PromptShell() {
	ed.EnsureCommand()
	ed.Unimplemented("PromptShell")
}

// :wa Enter : Save all files.
func (ed *Editor) PromptSaveAll() {
	ed.EnsureCommand()
	ed.Unimplemented("PromptSaveAll")
}

// :qa Enter : Close all files and quit editor.
func (ed *Editor) PromptQuitAll() {
	ed.EnsureCommand()
	ed.Unimplemented("PromptQuitAll")
}

// :qa! Enter : Force close all files and quit editor.
func (ed *Editor) PromptForceQuitAll() {
	ed.EnsureCommand()
	ed.alive = false
}

// :set ts=<num> Enter
func (ed *Editor) PromptTabStop(n int) {
	ed.EnsureCommand()
	if n < 1 {
		ed.Ring("set: the ts option may never be set to 0.")
		return
	}
	ed.cfg.TabWidth = n
	termi.TabWidth = n
	ed.redraw = true
}

// :set ai Enter
func (ed *Editor) PromptAutoIndent() {
	ed.EnsureCommand()
	ed.cfg.AutoIndent = true
}

// :set noai Enter
func (ed *Editor) PromptNoAutoIndent() {
	ed.EnsureCommand()
	ed.cfg.AutoIndent = false
}
