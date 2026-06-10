package editor

/////////////////////
// Prompt Commands //
/////////////////////

// :+<num> Enter : Move cursor to first non-blank character of next line.
func (ed *Editor) PromptMoveByLine(n int) {
	ed.EnsureCommand()
	if !ed.UpdateRow(n) {
		ed.Ring("Illegal address: only %d lines in the file.", len(ed.lines))
		return
	}
	ed.MoveToNonBlank()
}

// :-<num> Enter : Move cursor to first non-blank character of previous line.
func (ed *Editor) PromptMoveBackwardByLine(n int) {
	ed.EnsureCommand()
	if ed.row-n == -1 {
		n++
	}
	if !ed.UpdateRow(-n) {
		ed.Ring("Reference to a line number less than 0.")
		return
	}
	ed.MoveToNonBlank()
}

// :<num> Enter : Move cursor to first non-blank character of line specifined by <num>.
func (ed *Editor) PromptMoveToLine(n int) { // n: 1-based
	ed.EnsureCommand()
	if n == 0 {
		n = 1
	}
	n--
	if !ed.UpdateRow(n - len(ed.lines)) {
		ed.Ring("Illegal address: only %d lines in the file.", len(ed.lines))
		return
	}
	ed.MoveToNonBlank()
}

// :wq Enter : Save current file and quit.
func (ed *Editor) PromptSaveAndQuit() {
	ed.EnsureCommand()
	if ed.modified && ed.path == "" {
		ed.Ring("File is a temporary; exit will discard modifications.")
		return
	}
	ed.Save()
	ed.alive = false
}

// :w Enter : Save current file.
func (ed *Editor) PromptSave() {
	ed.EnsureCommand()
	ed.Save()
}

// :w! Enter : Force save current file.
func (ed *Editor) PromptForceSave() {
	ed.EnsureCommand()
	ed.Save()
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
