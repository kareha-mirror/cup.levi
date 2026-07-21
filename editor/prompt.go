package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"tea.kareha.org/cup/termi"
	"tea.kareha.org/cup/termi/shutil"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/color"
)

/////////////////////
// Prompt Commands //
/////////////////////

// Move cursor to first non-blank character of next line.
// Key: :+<num> Enter
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

// Move cursor to first non-blank character of previous line.
// Key: :-<num> Enter
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

// Move cursor to first non-blank character of line specifined by <num>.
// Key: :<num> Enter
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

// Save current file and quit.
// Key: :wq Enter
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

// Save current file.
// Key: :w Enter
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

// Force save current file.
// Key: :w! Enter
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

// Quit editor.
// Key: :q Enter
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

// Force quit editor.
// Key: :q! Enter
func (ed *Editor) PromptForceQuit() {
	ed.Close(true)
	ed.CheckQuit()
}

// Load file.
// Key :e Enter
func (ed *Editor) PromptLoad(name string) {
	b := ed.Buf()
	b.BeginSnapshot()
	if !ed.Load(name, false) {
		b.CancelSnapshot()
		return
	}
	b.EndSnapshot()
	ed.ShowFileInfo()
}

// Force load file.
// Key: :e! Enter
func (ed *Editor) PromptForceLoad(name string) {
	b := ed.Buf()
	b.BeginSnapshot()
	if !ed.Load(name, true) {
		b.CancelSnapshot()
		return
	}
	b.EndSnapshot()
	ed.ShowFileInfo()
}

// Read file and insert to current buffer.
// Key :r Enter
func (ed *Editor) PromptRead(name string) {
	b := ed.Buf()
	if name == "" {
		name = b.Path
	}
	data, err := ed.hooks.ReadFile(name)
	if err != nil {
		ed.Error("%v", err)
		return
	}
	ed.BeginUndoRecord()
	inserts, _ := buf.TextToLines(string(data))
	lines := append([]string(nil), b.Lines[:b.Loc.Row+1]...)
	lines = append(lines, inserts...)
	if b.Loc.Row+1 <= b.NumLines() {
		lines = append(lines, b.Lines[b.Loc.Row+1:]...)
	}
	b.Lines = lines
	b.Loc.Row++
	b.Loc = b.ConfineInclusive(b.Loc)
	ed.EndUndoRecord()
	ed.ShowTextInfo(name, inserts, b.CRLF)
}

// Execute shell.
// Key: :sh Enter
func (ed *Editor) PromptShell() {
	if ed.hooks.Shell != nil {
		termi.FinishKey()
		fmt.Print(termi.Clear)
		fmt.Print(termi.HomeCursor)
		fmt.Printf(termi.ResetAlternate)
		termi.Cooked()
		fmt.Print(termi.ShowCursor)

		err := ed.hooks.Shell()

		fmt.Print(termi.HideCursor)
		termi.Raw()
		fmt.Printf(termi.SetAlternate)
		termi.InitKey()
		ed.redraw = true

		if err != nil {
			ed.Error("%v", err)
		}
		return
	}

	shell := shutil.Path()
	cmd := exec.Command(shell)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	termi.FinishKey()
	fmt.Print(termi.Clear)
	fmt.Print(termi.HomeCursor)
	fmt.Printf(termi.ResetAlternate)
	termi.Cooked()
	fmt.Print(termi.ShowCursor)

	err := cmd.Run()

	fmt.Print(termi.HideCursor)
	termi.Raw()
	fmt.Printf(termi.SetAlternate)
	termi.InitKey()
	ed.redraw = true

	if err != nil {
		ed.Error("%v", err)
	}
}

// Save all files.
// Key: :wa Enter
func (ed *Editor) PromptSaveAll() {
	bufIdx := ed.bufIdx
	for i := 0; i < ed.NumBufs(); i++ {
		ed.bufIdx = i
		b := ed.Buf()
		if b.Modified {
			if !ed.Save(false) {
				break
			}
		} else {
			if b.Path == "" {
				ed.Message("(memory) is not modified")
			} else {
				ed.Message("%s is not modified", b.Path)
			}
		}
	}
	if ed.bufIdx < ed.NumBufs()-1 {
		ed.Ring("Not all files processed") // levi
	}
	ed.bufIdx = bufIdx
}

// Force save all files.
// Key: :wa! Enter
func (ed *Editor) PromptForceSaveAll() {
	bufIdx := ed.bufIdx
	for i := 0; i < ed.NumBufs(); i++ {
		ed.bufIdx = i
		if !ed.Save(true) {
			break
		}
	}
	if ed.bufIdx < ed.NumBufs()-1 {
		ed.Ring("Not all files processed") // levi
	}
	ed.bufIdx = bufIdx
}

// Close all files and quit editor.
// Key: :qa Enter
func (ed *Editor) PromptQuitAll() {
	for ed.alive {
		if !ed.Close(false) {
			return
		}
		ed.CheckQuit()
	}
}

// Force close all files and quit editor.
// Key: :qa! Enter
func (ed *Editor) PromptForceQuitAll() {
	ed.alive = false
}

// Set tab stop size.
// Key: :set ts=<num> Enter
func (ed *Editor) PromptTabStop(n int) {
	if n < 1 {
		ed.Ring("set: the ts option may never be set to 0.")
		return
	}
	ed.cfg.TabStop = n
	termi.TabWidth = n
}

// Set auto indent enabled.
// Key: :set ai Enter
func (ed *Editor) PromptAutoIndent() {
	ed.cfg.AutoIndent = true
}

// Set auto indent disabled.
// Key: :set noai Enter
func (ed *Editor) PromptNoAutoIndent() {
	ed.cfg.AutoIndent = false
}

// Open file in new buffer.
// Key: :open Enter
func (ed *Editor) PromptOpen(name string) {
	if !ed.Open(name) {
		return
	}
	ed.ShowFileInfo()
}

// Set newline type.
// Key: :newline Enter
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

// Set colorscheme.
// Key: :colors Enter
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
		ed.Message("%s", strings.Join(list.Names, " "))
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

// Show memory usage.
// Key: :mem Enter
func (ed *Editor) PromptMem() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	sb := strings.Builder{}

	if m.Alloc >= 1024*1024 || m.HeapAlloc >= 1024*1024 {
		sb.WriteString(fmt.Sprintf(
			"Alloc     = % 6d MiB\n", m.Alloc/1024/1024,
		))
		sb.WriteString(fmt.Sprintf(
			"HeapAlloc = % 6d MiB\n", m.HeapAlloc/1024/1024,
		))
	} else {
		sb.WriteString(fmt.Sprintf(
			"Alloc     = % 6d KiB\n", m.Alloc/1024,
		))
		sb.WriteString(fmt.Sprintf(
			"HeapAlloc = % 6d KiB\n", m.HeapAlloc/1024,
		))
	}
	sb.WriteRune('\n')
	if m.HeapSys >= 1024*1024 || m.Sys >= 1024*1024 {
		sb.WriteString(fmt.Sprintf(
			"HeapSys   = % 6d MiB\n", m.HeapSys/1024/1024,
		))
		sb.WriteString(fmt.Sprintf(
			"Sys       = % 6d MiB\n", m.Sys/1024/1024,
		))
	} else {
		sb.WriteString(fmt.Sprintf(
			"HeapSys   = % 6d KiB\n", m.HeapSys/1024,
		))
		sb.WriteString(fmt.Sprintf(
			"Sys       = % 6d KiB\n", m.Sys/1024,
		))
	}
	sb.WriteRune('\n')
	sb.WriteString(fmt.Sprintf("NumGC     = % 6d\n", m.NumGC))

	ed.Message("%s", sb.String())
}

// Used by debug.
// Key: :hello Enter
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
