package editor

func (ed *Editor) RunPrompt(c Pcmd) bool {
	ed.Commit()

	switch c.Kind {
	case PcmdMoveToLine:
		ed.PromptMoveToLine(c.Num)
		return true

	case PcmdSaveAndQuit:
		ed.PromptSaveAndQuit()
		return true
	case PcmdSave:
		ed.PromptSave(c.Name)
		return true
	case PcmdForceSave:
		ed.PromptForceSave(c.Name)
		return true
	case PcmdQuit:
		ed.PromptQuit()
		return true
	case PcmdForceQuit:
		ed.PromptForceQuit()
		return true
	case PcmdOpen:
		ed.PromptOpen(c.Name)
		return true
	case PcmdForceOpen:
		ed.PromptForceOpen(c.Name)
		return true
	case PcmdRead:
		ed.PromptRead()
		return true
	case PcmdNext:
		ed.PromptNext()
		return true
	case PcmdPrev:
		ed.PromptPrev()
		return true

	case PcmdShell:
		ed.PromptShell()
		return true

	case PcmdSaveAll:
		ed.PromptSaveAll()
		return true
	case PcmdQuitAll:
		ed.PromptQuitAll()
		return true
	case PcmdForceQuitAll:
		ed.PromptForceQuitAll()
		return true

	case PcmdTabStop:
		ed.PromptTabStop(c.Num)
		return true
	case PcmdAutoIndent:
		ed.PromptAutoIndent()
		return true
	case PcmdNoAutoIndent:
		ed.PromptNoAutoIndent()
		return true

	case PcmdColors:
		ed.PromptColors(c.Name)
		return true
	}

	return false
}
