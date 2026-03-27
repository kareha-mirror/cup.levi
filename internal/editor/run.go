package editor

func (ed *Editor) Run(c Cmd) bool {
	switch c.Kind {
	case CmdInvalid:
		ed.Ring("not (yet) a vi command [" + ed.parser.String() + "]")
		ed.parser.Clear()
		return true

	case CmdMoveLeft:
		ed.MoveLeft(c.Num)
		return true
	case CmdMoveDown:
		ed.MoveDown(c.Num)
		return true
	case CmdMoveUp:
		ed.MoveUp(c.Num)
		return true
	case CmdMoveRight:
		ed.MoveRight(c.Num)
		return true

	case CmdMoveToStart:
		ed.MoveToStart()
		return true
	case CmdMoveToEnd:
		ed.MoveToEnd()
		return true
	case CmdMoveToNonBlank:
		ed.MoveToNonBlank()
		return true
	case CmdMoveToColumn:
		ed.MoveToColumn(c.Num)
		return true

	// TODO

	case CmdInsertBefore:
		ed.InsertBefore(c.Num)
		return true
	case CmdInsertAfter:
		ed.InsertAfter(c.Num)
		return true

	// TODO

	case CmdOpDelete:
		ed.OpDelete(c.Num)
		return true

	// TODO

	case CmdMiscSaveAndQuit:
		ed.MiscSaveAndQuit()
		return true

		// TODO
	}

	return false
}
