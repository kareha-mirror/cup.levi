package editor

import "tea.kareha.org/cup/levi/internal/prompt"

func (ed *Editor) RunPrompt(c prompt.Cmd) bool {
	ed.Commit()
	switch c.Kind {

	case prompt.MoveToLine:
		ed.PromptMoveToLine(c.Num)
		return true

	case prompt.SaveAndClose:
		ed.PromptSaveAndClose()
		return true
	case prompt.WriteAndClose:
		ed.PromptWriteAndClose()
		return true
	case prompt.Write:
		ed.PromptWrite(c.Name)
		return true
	case prompt.ForceWrite:
		ed.PromptForceWrite(c.Name)
		return true
	case prompt.Close:
		ed.PromptClose()
		return true
	case prompt.ForceClose:
		ed.PromptForceClose()
		return true
	case prompt.Load:
		ed.PromptLoad(c.Name)
		return true
	case prompt.ForceLoad:
		ed.PromptForceLoad(c.Name)
		return true
	case prompt.Read:
		ed.PromptRead(c.Name)
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
	case prompt.ForceSaveAll:
		ed.PromptForceSaveAll()
		return true
	case prompt.CloseAll:
		ed.PromptCloseAll()
		return true
	case prompt.ForceCloseAll:
		ed.PromptForceCloseAll()
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
