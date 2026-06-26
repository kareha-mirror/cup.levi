package editor

import (
	"strings"

	"tea.kareha.org/cup/termi"

	"tea.kareha.org/cup/levi/internal/colors"
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
	ed.Commit()
	b := ed.Buf()
	if !b.CheckRowInclusive(b.Loc.Row + n) {
		ed.Ring("Illegal address: only %d lines in the file.", b.NumLines())
		return
	}
	b.Loc.Row += n
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
}

// :-<num> Enter : Move cursor to first non-blank character of previous line.
func (ed *Editor) PromptMoveBackwardByLine(n int) {
	if n < 0 {
		ed.Error("PromptMoveBackwardByLine: n < 0")
		return
	}
	ed.Commit()
	b := ed.Buf()
	if b.Loc.Row-n == -1 {
		n++
	}
	if !b.CheckRowInclusive(b.Loc.Row - n) {
		ed.Ring("Reference to a line number less than 0.")
		return
	}
	b.Loc.Row -= n
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
}

// :<num> Enter : Move cursor to first non-blank character of line specifined by <num>.
func (ed *Editor) PromptMoveToLine(n int) { // n: 1-based
	if n < 0 {
		ed.Error("PromptMoveToLine: n < 0")
		return
	}
	ed.Commit()
	b := ed.Buf()
	if n == 0 {
		n = 1
	}
	if !b.CheckRowInclusive(n - 1) {
		ed.Ring("Illegal address: only %d lines in the file.", b.NumLines())
		return
	}
	b.Loc.Row = n - 1
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
}

// :wq Enter : Save current file and quit.
func (ed *Editor) PromptSaveAndQuit() {
	ed.Commit()
	b := ed.Buf()
	if b.Modified && b.Path == "" {
		ed.Ring("File is a temporary; exit will discard modifications.")
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

// :w Enter : Save current file.
func (ed *Editor) PromptSave(name string) {
	ed.Commit()
	if name == "" {
		ed.Save(false)
		return
	}
	ed.SaveAs(name, false)
}

// :w! Enter : Force save current file.
func (ed *Editor) PromptForceSave(name string) {
	ed.Commit()
	if name == "" {
		ed.Save(true)
		return
	}
	ed.SaveAs(name, true)
}

// :q Enter : Quit editor.
func (ed *Editor) PromptQuit() {
	ed.Commit()
	b := ed.Buf()
	if b.Modified {
		if b.Path == "" {
			ed.Ring("File is a temporary; exit will discard modifications.")
			return
		}
		ed.Ring("File modified since last complete write; write or use ! to override.")
		return
	}
	ed.Close(false)
}

// :q! Enter : Force quit editor.
func (ed *Editor) PromptForceQuit() {
	ed.Commit()
	ed.Close(true)
}

// :e Enter : Open file.
func (ed *Editor) PromptOpen(name string) {
	ed.Commit()
	ed.Load(name, false)
	ed.InitialInfo()
}

// :e! Enter : Force open file.
func (ed *Editor) PromptForceOpen(name string) {
	ed.Commit()
	ed.Load(name, true)
	ed.InitialInfo()
}

// :r Enter : Read file and insert to current buffer.
func (ed *Editor) PromptRead() {
	ed.Commit()
	ed.Unimplemented("PromptRead")
}

// :n Enter : Switch to next buffer (tab).
func (ed *Editor) PromptNext() {
	ed.Commit()
	if ed.bufIdx+1 >= len(ed.bufs) {
		ed.Ring("No more files to edit.")
		return
	}
	ed.bufIdx++
	ed.redraw = true
	ed.InitialInfo()
}

// :prev Enter : Switch to previous buffer (tab).
func (ed *Editor) PromptPrev() {
	ed.Commit()
	if ed.bufIdx-1 < 0 {
		ed.Ring("No previous files to edit.")
		return
	}
	ed.bufIdx--
	ed.redraw = true
	ed.InitialInfo()
}

// :sh Enter : Execute shell.
func (ed *Editor) PromptShell() {
	ed.Commit()
	ed.Unimplemented("PromptShell")
}

// :wa Enter : Save all files.
func (ed *Editor) PromptSaveAll() {
	ed.Commit()
	ed.Unimplemented("PromptSaveAll")
}

// :qa Enter : Close all files and quit editor.
func (ed *Editor) PromptQuitAll() {
	ed.Commit()
	ed.Unimplemented("PromptQuitAll")
}

// :qa! Enter : Force close all files and quit editor.
func (ed *Editor) PromptForceQuitAll() {
	ed.Commit()
	ed.alive = false
}

// :set ts=<num> Enter
func (ed *Editor) PromptTabStop(n int) {
	ed.Commit()
	if n < 1 {
		ed.Ring("set: the ts option may never be set to 0.")
		return
	}
	ed.cfg.TabStop = n
	termi.TabWidth = n
	ed.redraw = true
}

// :set ai Enter
func (ed *Editor) PromptAutoIndent() {
	ed.Commit()
	ed.cfg.AutoIndent = true
}

// :set noai Enter
func (ed *Editor) PromptNoAutoIndent() {
	ed.Commit()
	ed.cfg.AutoIndent = false
}

// :colors Enter
func (ed *Editor) PromptColors(name string) {
	ed.Commit()

	if name == "." {
		colors, err := colors.Parse(ed.Buf().Text())
		if err != nil {
			ed.Error("%v", err)
			return
		}
		ed.colors = colors
		ed.redraw = true
		return
	}

	list, err := colors.LoadList(ed.dir)
	if err != nil {
		ed.Error("%v", err)
		return
	}

	if name == "" {
		ed.Message(strings.Join(list.Names, " "))
		return
	}

	colors, err := list.Load(name)
	if err != nil {
		ed.Error("%v", err)
		return
	}
	ed.colors = colors
	ed.redraw = true
}
