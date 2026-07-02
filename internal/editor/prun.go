package editor

import (
	"tea.kareha.org/cup/levi/internal/prompt"
)

func (ed *Editor) RunPrompt(c prompt.Cmd) bool {
	ed.Commit()
	switch c.Kind {

	case prompt.MoveToLine:
		ed.PromptMoveToLine(c.Num)
		return true

	case prompt.SaveAndQuit:
		ed.PromptSaveAndQuit()
		return true
	case prompt.Save:
		ed.PromptSave(c.Name)
		return true
	case prompt.ForceSave:
		ed.PromptForceSave(c.Name)
		return true
	case prompt.Quit:
		ed.PromptQuit()
		return true
	case prompt.ForceQuit:
		ed.PromptForceQuit()
		return true
	case prompt.Load:
		ed.PromptLoad(c.Name)
		return true
	case prompt.ForceLoad:
		ed.PromptForceLoad(c.Name)
		return true
	case prompt.Read:
		ed.PromptRead()
		return true
	case prompt.Next:
		ed.NextBuf()
		return true
	case prompt.Prev:
		ed.PrevBuf()
		return true

	case prompt.Shell:
		ed.PromptShell()
		return true

	case prompt.SaveAll:
		ed.PromptSaveAll()
		return true
	case prompt.QuitAll:
		ed.PromptQuitAll()
		return true
	case prompt.ForceQuitAll:
		ed.PromptForceQuitAll()
		return true

	case prompt.TabStop:
		ed.PromptTabStop(c.Num)
		return true
	case prompt.AutoIndent:
		ed.PromptAutoIndent()
		return true
	case prompt.NoAutoIndent:
		ed.PromptNoAutoIndent()
		return true

	case prompt.Open:
		ed.PromptOpen(c.Name)
		return true
	case prompt.Newline:
		ed.PromptNewline(c.Name)
		return true
	case prompt.Colors:
		ed.PromptColors(c.Name)
		return true

	case prompt.Mem:
		ed.PromptMem()
		return true
	case prompt.Hello:
		ed.PromptHello(c.Num)
		return true

	case prompt.Ring:
		ed.Ring("%s", c.Name)
	case prompt.Error:
		ed.Error("%s", c.Name)

	}
	return false
}
