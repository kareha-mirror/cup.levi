package editor

func (ed *Editor) ParseLetter(num int, op string, letter rune) (Cmd, bool) {
	if letter == 0 {
		return Cmd{}, false
	}

	switch op {

	case "m":
		return Cmd{Kind: CmdMarkSet, Letter: letter}, true

	case "r":
		return Cmd{
			Kind:   CmdEditReplace,
			Num:    num,
			Letter: letter,
		}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseView(num int, op string) (Cmd, bool) {
	switch op {

	case "\x06": // Ctrl-F
		return Cmd{Kind: CmdViewDown, Num: num}, true
	case "\x02": // Ctrl-B
		return Cmd{Kind: CmdViewUp, Num: num}, true
	case "\x04": // Ctrl-D
		return Cmd{Kind: CmdViewDownHalf, Num: num}, true
	case "\x15": // Ctrl-U
		return Cmd{Kind: CmdViewUpHalf, Num: num}, true
	case "\x19": // Ctrl-Y
		return Cmd{Kind: CmdViewDownLine, Num: num}, true
	case "\x05": // Ctrl-E
		return Cmd{Kind: CmdViewUpLine, Num: num}, true

	case "z\r":
		return Cmd{Kind: CmdViewToTop}, true
	case "z.":
		return Cmd{Kind: CmdViewToMiddle}, true
	case "z-":
		return Cmd{Kind: CmdViewToBottom}, true

	case "\x0c": // Ctrl-L
		return Cmd{Kind: CmdViewRedraw}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseInsert(num int, op string) (Cmd, bool) {
	switch op {

	case "i":
		return Cmd{Kind: CmdInsertBefore, Num: num}, true
	case "a":
		return Cmd{Kind: CmdInsertAfter, Num: num}, true
	case "I":
		return Cmd{Kind: CmdInsertBeforeNonBlank, Num: num}, true
	case "A":
		return Cmd{Kind: CmdInsertAfterEnd, Num: num}, true
	case "R":
		return Cmd{Kind: CmdInsertOverwrite, Num: num}, true

	case "o":
		return Cmd{Kind: CmdInsertOpenBelow, Num: num}, true
	case "O":
		return Cmd{Kind: CmdInsertOpenAbove, Num: num}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseMisc(num int, op string) (Cmd, bool) {
	switch op {

	case "\x07": // Ctrl-G
		return Cmd{Kind: CmdMiscShowInfo}, true
	case ".":
		return Cmd{Kind: CmdMiscRepeat, Num: num}, true
	case "u":
		return Cmd{Kind: CmdMiscUndo, Num: num}, true
	case "U":
		return Cmd{Kind: CmdMiscRestore}, true
	case "ZZ":
		return Cmd{Kind: CmdMiscSaveAndQuit}, true
	case "\x1a": // Ctrl-Z
		return Cmd{Kind: CmdMiscSuspend}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseOp(
	reg string, num int, op string, noSubnum bool, subnum int,
	mv string, letter rune,
) (CmdPair, bool) {
	if mv != "" {
		switch op {

		case "y":
			cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
			if ok {
				return CmdPair{
					Reg:  reg,
					Main: Cmd{Kind: CmdOpCopyRegion},
					Sub:  cmd,
				}, true
			}
			return CmdPair{}, false
		case "d":
			cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
			if ok {
				return CmdPair{
					Reg:  reg,
					Main: Cmd{Kind: CmdOpDeleteRegion},
					Sub:  cmd,
				}, true
			}
			return CmdPair{}, false
		case "c":
			cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
			if cmd.Kind == CmdMoveByWord {
				cmd.Kind = CmdMoveByWordEx
			}
			if ok {
				return CmdPair{
					Reg:  reg,
					Main: Cmd{Kind: CmdOpChangeRegion},
					Sub:  cmd,
				}, true
			}
			return CmdPair{}, false

		}
	}

	switch op {

	case "yy", "Y":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpCopyLine, Num: num},
		}, true
	case "yw":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpCopyWord, Num: num},
		}, true
	case "y$":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpCopyToEnd, Num: num},
		}, true

	case "p":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpPaste, Num: num},
		}, true
	case "P":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpPasteBefore, Num: num},
		}, true

	case "x":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpDelete, Num: num},
		}, true
	case "X":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpDeleteBefore, Num: num},
		}, true
	case "dd":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpDeleteLine, Num: num},
		}, true
	case "dw":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpDeleteWord, Num: num},
		}, true
	case "d$", "D":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpDeleteToEnd, Num: num},
		}, true

	case "cc":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpChangeLine, Num: num},
		}, true
	case "cw":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpChangeWord, Num: num},
		}, true
	case "C":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpChangeToEnd, Num: num},
		}, true
	case "s":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpSubst, Num: num},
		}, true
	case "S":
		return CmdPair{
			Reg:  reg,
			Main: Cmd{Kind: CmdOpSubstLine, Num: num},
		}, true

	}

	return CmdPair{}, false
}

func (ed *Editor) ParseEdit(
	num int, op string, noSubnum bool, subnum int, mv string, letter rune,
) (CmdPair, bool) {
	switch op {

	case "J":
		return CmdPair{
			Main: Cmd{Kind: CmdEditJoin, Num: num},
		}, true

	case ">>":
		return CmdPair{
			Main: Cmd{Kind: CmdEditIndent, Num: num},
		}, true
	case "<<":
		return CmdPair{
			Main: Cmd{Kind: CmdEditOutdent, Num: num},
		}, true

	case ">":
		cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
		if ok {
			meta, ok := MoveMetas[cmd.Kind]
			if !ok {
				return CmdPair{}, false
			}
			if meta.Linewise {
				return CmdPair{
					Main: Cmd{Kind: CmdEditIndentRegion},
					Sub:  cmd,
				}, true
			} else {
				return CmdPair{
					Main: Cmd{Kind: CmdEditIndentRegion},
					Sub:  cmd,
				}, true
			}
		}
		return CmdPair{}, false
	case "<":
		cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
		if ok {
			meta, ok := MoveMetas[cmd.Kind]
			if !ok {
				return CmdPair{}, false
			}
			if meta.Linewise {
				return CmdPair{
					Main: Cmd{Kind: CmdEditOutdentRegion},
					Sub:  cmd,
				}, true
			} else {
				return CmdPair{
					Main: Cmd{Kind: CmdEditOutdentRegion},
					Sub:  cmd,
				}, true
			}
		}
		return CmdPair{}, false

	}

	return CmdPair{}, false
}
