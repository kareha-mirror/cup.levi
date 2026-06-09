package editor

import (
	"strings"
)

/////////////////////
// Prompt Commands //
/////////////////////

// : : Prompt mode.
func (ed *Editor) PromptMode() {
	ed.mode = ModePrompt
}

func (ed *Editor) ParsePrompt() (Cmd, bool) {
	if ed.prompt.Len() < 1 {
		return Cmd{Kind: CmdInvalid}, false
	}

	parts := strings.Split(ed.prompt.String(), " ")

	// TODO
	//return Cmd{Kind: CmdPromptMoveToLine}, true

	switch parts[0] {
	case "wq":
		return Cmd{Kind: CmdPromptSaveAndQuit}, true
	case "w":
		return Cmd{Kind: CmdPromptSave}, true
	case "w!":
		return Cmd{Kind: CmdPromptForceSave}, true
	case "q":
		return Cmd{Kind: CmdPromptQuit}, true
	case "q!":
		return Cmd{Kind: CmdPromptForceQuit}, true
	case "e":
		return Cmd{Kind: CmdPromptOpen}, true
	case "e!":
		return Cmd{Kind: CmdPromptForceOpen}, true
	case "r":
		return Cmd{Kind: CmdPromptRead}, true
	case "n":
		return Cmd{Kind: CmdPromptNext}, true
	case "prev":
		return Cmd{Kind: CmdPromptPrev}, true

	case "sh":
		return Cmd{Kind: CmdPromptShell}, true

	case "wa":
		return Cmd{Kind: CmdPromptSaveAll}, true
	case "qa":
		return Cmd{Kind: CmdPromptQuitAll}, true
	case "qa!":
		return Cmd{Kind: CmdPromptForceQuitAll}, true

	default:
		return Cmd{Kind: CmdInvalid}, false
	}
}

// :<num> Enter : Move cursor to line <num>.
func (ed *Editor) PromptMoveToLine(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("PromptMoveToLine")
}

// :wq Enter : Save current file and quit.
func (ed *Editor) PromptSaveAndQuit() {
	ed.EnsureCommand()
	ed.save = true
	ed.quit = true
}

// :w Enter : Save current file.
func (ed *Editor) PromptSave() {
	ed.EnsureCommand()
	ed.Unimplemented("PromptSave")
}

// :w! Enter : Force save current file.
func (ed *Editor) PromptForceSave() {
	ed.EnsureCommand()
	ed.Unimplemented("PromptForceSave")
}

// :q Enter : Quit editor.
func (ed *Editor) PromptQuit() {
	ed.EnsureCommand()
	ed.quit = true
}

// :q! Enter : Force quit editor.
func (ed *Editor) PromptForceQuit() {
	ed.EnsureCommand()
	ed.Unimplemented("PromptForceQuit")
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
	ed.Unimplemented("PromptForceQuitAll")
}
