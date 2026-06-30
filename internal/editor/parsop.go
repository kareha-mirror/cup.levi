package editor

func (ed *Editor) ParseLetter(num int, op string, letter rune) (Cmd, bool) {
	if letter == 0 {
		return Cmd{}, false
	}

	switch op {

	case "m":
		return Cmd{Kind: Mark, Ltr: letter}, true

	case "r":
		return Cmd{
			Kind: Replace,
			Num:  num,
			Ltr:  letter,
		}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseView(num int, op string) (Cmd, bool) {
	switch op {

	case "\x06": // Ctrl-F
		return Cmd{Kind: ViewDown, Num: num}, true
	case "\x02": // Ctrl-B
		return Cmd{Kind: ViewUp, Num: num}, true
	case "\x04": // Ctrl-D
		return Cmd{Kind: ViewDownHalf, Num: num}, true
	case "\x15": // Ctrl-U
		return Cmd{Kind: ViewUpHalf, Num: num}, true
	case "\x19": // Ctrl-Y
		return Cmd{Kind: ViewDownLine, Num: num}, true
	case "\x05": // Ctrl-E
		return Cmd{Kind: ViewUpLine, Num: num}, true

	case "z\r":
		return Cmd{Kind: ViewToTop}, true
	case "z.":
		return Cmd{Kind: ViewToMiddle}, true
	case "z-":
		return Cmd{Kind: ViewToBottom}, true

	case "\x0c": // Ctrl-L
		return Cmd{Kind: Redraw}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseInsert(num int, op string) (Cmd, bool) {
	switch op {

	case "i":
		return Cmd{Kind: InsertBefore, Num: num}, true
	case "a":
		return Cmd{Kind: InsertAfter, Num: num}, true
	case "I":
		return Cmd{Kind: InsertBeforeNonBlank, Num: num}, true
	case "A":
		return Cmd{Kind: InsertAfterEnd, Num: num}, true
	case "R":
		return Cmd{Kind: Overwrite, Num: num}, true

	case "o":
		return Cmd{Kind: OpenBelow, Num: num}, true
	case "O":
		return Cmd{Kind: OpenAbove, Num: num}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseMisc(num int, op string) (Cmd, bool) {
	switch op {

	case "\x07": // Ctrl-G
		return Cmd{Kind: ShowInfo}, true
	case ".":
		return Cmd{Kind: Repeat, Num: num}, true
	case "u":
		return Cmd{Kind: Undo, Num: num}, true
	case "U":
		return Cmd{Kind: Restore}, true
	case "ZZ":
		return Cmd{Kind: SaveAndClose}, true
	case "\x1a": // Ctrl-Z
		return Cmd{Kind: Suspend}, true

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
					Reg: reg,
					Op:  Cmd{Kind: CopyRegion},
					Mv:  cmd,
				}, true
			}
			return CmdPair{}, false
		case "d":
			cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
			if ok {
				return CmdPair{
					Reg: reg,
					Op:  Cmd{Kind: DeleteRegion},
					Mv:  cmd,
				}, true
			}
			return CmdPair{}, false
		case "c":
			cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
			if cmd.Kind == MoveByWord {
				cmd.Kind = MoveByChangeWord
			}
			if ok {
				return CmdPair{
					Reg: reg,
					Op:  Cmd{Kind: ChangeRegion},
					Mv:  cmd,
				}, true
			}
			return CmdPair{}, false

		}
	}

	switch op {

	case "yy", "Y":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: CopyLine, Num: num},
		}, true
	case "yw":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: CopyWord, Num: num},
		}, true
	case "y$":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: CopyToEnd, Num: num},
		}, true

	case "p":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: Paste, Num: num},
		}, true
	case "P":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: PasteBefore, Num: num},
		}, true

	case "x":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: Delete, Num: num},
		}, true
	case "X":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: DeleteBefore, Num: num},
		}, true
	case "dd":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: DeleteLine, Num: num},
		}, true
	case "dw":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: DeleteWord, Num: num},
		}, true
	case "d$", "D":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: DeleteToEnd, Num: num},
		}, true

	case "cc":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: ChangeLine, Num: num},
		}, true
	case "cw":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: ChangeWord, Num: num},
		}, true
	case "C":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: ChangeToEnd, Num: num},
		}, true
	case "s":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: Subst, Num: num},
		}, true
	case "S":
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: SubstLine, Num: num},
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
			Op: Cmd{Kind: Join, Num: num},
		}, true

	case ">>":
		return CmdPair{
			Op: Cmd{Kind: Indent, Num: num},
		}, true
	case "<<":
		return CmdPair{
			Op: Cmd{Kind: Outdent, Num: num},
		}, true

	case ">":
		cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
		if ok {
			attr, ok := MoveAttrs[cmd.Kind]
			if !ok {
				return CmdPair{}, false
			}
			if attr.Linewise {
				return CmdPair{
					Op: Cmd{Kind: IndentRegion},
					Mv: cmd,
				}, true
			} else {
				return CmdPair{
					Op: Cmd{Kind: IndentRegion},
					Mv: cmd,
				}, true
			}
		}
		return CmdPair{}, false
	case "<":
		cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
		if ok {
			attr, ok := MoveAttrs[cmd.Kind]
			if !ok {
				return CmdPair{}, false
			}
			if attr.Linewise {
				return CmdPair{
					Op: Cmd{Kind: OutdentRegion},
					Mv: cmd,
				}, true
			} else {
				return CmdPair{
					Op: Cmd{Kind: OutdentRegion},
					Mv: cmd,
				}, true
			}
		}
		return CmdPair{}, false

	}

	return CmdPair{}, false
}
