package editor

import (
	"fmt"
	"runtime"
	"strings"

	"tea.kareha.org/cup/termi"

	"tea.kareha.org/cup/levi/internal/color"
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
	b := ed.Buf()
	row := b.Loc.Row + n
	if !b.CheckRowInclusive(row) {
		ed.Ring("Illegal address: only %d lines in the file.", b.NumLines())
		return
	}
	b.Loc.Row = row
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
}

// :-<num> Enter : Move cursor to first non-blank character of previous line.
func (ed *Editor) PromptMoveBackwardByLine(n int) {
	if n < 0 {
		ed.Error("PromptMoveBackwardByLine: n < 0")
		return
	}
	b := ed.Buf()
	row := b.Loc.Row - n
	if row == -1 {
		row++
	}
	if !b.CheckRowInclusive(row) {
		ed.Ring("Reference to a line number less than 0.")
		return
	}
	b.Loc.Row = row
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
}

// :<num> Enter : Move cursor to first non-blank character of line specifined by <num>.
func (ed *Editor) PromptMoveToLine(n int) { // n: 1-based
	if n < 0 {
		ed.Error("PromptMoveToLine: n < 0")
		return
	}
	if n == 0 {
		n = 1
	}
	b := ed.Buf()
	row := n - 1
	if !b.CheckRowInclusive(row) {
		ed.Ring("Illegal address: only %d lines in the file.", b.NumLines())
		return
	}
	b.Loc.Row = row
	b.Loc.Col = b.NonBlankColOfLine(b.Loc.Row)
}

// :wq Enter : Save current file and quit.
func (ed *Editor) PromptSaveAndQuit() {
	b := ed.Buf()
	if b.Modified && b.Path == "" {
		ed.Ring("File is a temporary; exit will discard modifications.")
		return
	}
	if b.Modified && b.Path != "" {
		if !ed.Save(false) {
			return
		}
	}
	ed.Close(false)
	ed.CheckQuit()
}

// :w Enter : Save current file.
func (ed *Editor) PromptSave(name string) {
	if name == "" {
		if !ed.Save(false) {
			return
		}
		return
	}
	if !ed.SaveAs(name, false) {
		return
	}
}

// :w! Enter : Force save current file.
func (ed *Editor) PromptForceSave(name string) {
	if name == "" {
		if !ed.Save(true) {
			return
		}
		return
	}
	if !ed.SaveAs(name, true) {
		return
	}
}

// :q Enter : Quit editor.
func (ed *Editor) PromptQuit() {
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
	ed.CheckQuit()
}

// :q! Enter : Force quit editor.
func (ed *Editor) PromptForceQuit() {
	ed.Close(true)
	ed.CheckQuit()
}

// :e Enter : Open file.
func (ed *Editor) PromptOpen(name string) {
	if !ed.Load(name, false) {
		return
	}
	ed.ShowFileInfo()
}

// :e! Enter : Force open file.
func (ed *Editor) PromptForceOpen(name string) {
	if !ed.Load(name, true) {
		return
	}
	ed.ShowFileInfo()
}

// :r Enter : Read file and insert to current buffer.
func (ed *Editor) PromptRead() {
	ed.Unimplemented("PromptRead")
}

// :n Enter : Switch to next buffer (tab).
func (ed *Editor) PromptNext() {
	if ed.bufIdx+1 >= len(ed.bufs) {
		ed.Ring("No more files to edit.")
		return
	}
	ed.bufIdx++
	ed.redraw = true
	ed.ShowFileInfo()
}

// :prev Enter : Switch to previous buffer (tab).
func (ed *Editor) PromptPrev() {
	if ed.bufIdx-1 < 0 {
		ed.Ring("No previous files to edit.")
		return
	}
	ed.bufIdx--
	ed.redraw = true
	ed.ShowFileInfo()
}

// :sh Enter : Execute shell.
func (ed *Editor) PromptShell() {
	ed.Unimplemented("PromptShell")
}

// :wa Enter : Save all files.
func (ed *Editor) PromptSaveAll() {
	ed.Unimplemented("PromptSaveAll")
}

// :qa Enter : Close all files and quit editor.
func (ed *Editor) PromptQuitAll() {
	ed.Unimplemented("PromptQuitAll")
}

// :qa! Enter : Force close all files and quit editor.
func (ed *Editor) PromptForceQuitAll() {
	ed.alive = false
}

// :set ts=<num> Enter
func (ed *Editor) PromptTabStop(n int) {
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
	ed.cfg.AutoIndent = true
}

// :set noai Enter
func (ed *Editor) PromptNoAutoIndent() {
	ed.cfg.AutoIndent = false
}

// :newline Enter
func (ed *Editor) PromptNewline(name string) {
	switch name {
	case "":
		if ed.Buf().CRLF {
			ed.Message("Newline is CRLF (Windows)")
		} else {
			ed.Message("Newline is LF (Unix)")
		}
	case "unix", "u", "linux", "lin", "l", "bsd", "b", "mac", "m":
		ed.Buf().CRLF = false
		ed.PromptNewline("")
	case "windows", "win", "w", "dos", "d":
		ed.Buf().CRLF = true
		ed.PromptNewline("")
	default:
		ed.Error("Please specify unix or windows")
	}
}

// :colors Enter
func (ed *Editor) PromptColors(name string) {
	// colors . : parse and load colorscheme from current buffer
	if name == "." {
		colors, err := color.ParseScheme(ed.Buf().Text(false))
		if err != nil {
			ed.Error("%v", err)
			return
		}
		ed.colors = colors
		ed.redraw = true
		return
	}

	list, err := color.LoadSchemeList(ed.cfgDir)
	if err != nil {
		ed.Error("%v", err)
		return
	}

	// colors : list registered colorschemes
	if name == "" {
		ed.Message(strings.Join(list.Names, " "))
		return
	}

	// colors <name> : locad colorscheme from list
	colors, err := list.Load(name)
	if err != nil {
		ed.Error("%v", err)
		return
	}
	ed.colors = colors
	ed.redraw = true
}

// :mem Enter
func (ed *Editor) PromptMem() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("Alloc     = % 8d KiB\n", m.Alloc/1024))
	sb.WriteString(fmt.Sprintf("HeapAlloc = % 8d KiB\n", m.HeapAlloc/1024))
	sb.WriteString(fmt.Sprintf("HeapSys   = % 8d KiB\n", m.HeapSys/1024))
	sb.WriteString(fmt.Sprintf("Sys       = % 8d KiB\n", m.Sys/1024))
	sb.WriteRune('\n')
	sb.WriteString(fmt.Sprintf("NumGC = %d\n", m.NumGC))

	ed.Message(sb.String())
}

// :hello Enter
func (ed *Editor) PromptHello(n int) {
	// hello : show Hello, World!
	if n < 1 {
		ed.Message("Hello, World!")
		return
	}

	// hello <num> : show list of numbers
	for i := 1; i <= n; i++ {
		ed.Message("%d", i)
	}
}
