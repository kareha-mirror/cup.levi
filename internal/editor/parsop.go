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
) (Cmd, bool) {
	if mv != "" {
		switch op {

		case "y":
			b := ed.Buf()
			start := b.Loc
			cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
			if ok {
				meta, ok := MoveMetas[cmd.Kind]
				if !ok {
					return Cmd{}, false
				}
				end, ok := ed.RunMove(cmd)
				if !ok {
					return Cmd{}, false
				}
				if meta.Linewise {
					return Cmd{
						Kind:  CmdOpCopyLineRegion,
						Reg:   reg,
						Start: start,
						End:   end,
					}, true
				} else {
					return Cmd{
						Kind:      CmdOpCopyRegion,
						Reg:       reg,
						Start:     start,
						End:       end,
						Inclusive: meta.Inclusive,
					}, true
				}
			}
			return Cmd{}, false
		case "d":
			b := ed.Buf()
			start := b.Loc
			cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
			if ok {
				meta, ok := MoveMetas[cmd.Kind]
				if !ok {
					return Cmd{}, false
				}
				end, ok := ed.RunMove(cmd)
				if !ok {
					return Cmd{}, false
				}
				if meta.Linewise {
					return Cmd{
						Kind:  CmdOpDeleteLineRegion,
						Reg:   reg,
						Start: start,
						End:   end,
					}, true
				} else {
					return Cmd{
						Kind:      CmdOpDeleteRegion,
						Reg:       reg,
						Start:     start,
						End:       end,
						Inclusive: meta.Inclusive,
					}, true
				}
			}
			return Cmd{}, false
		case "c":
			b := ed.Buf()
			start := b.Loc
			cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
			if ok {
				meta, ok := MoveMetas[cmd.Kind]
				if !ok {
					return Cmd{}, false
				}
				end, ok := ed.RunMove(cmd)
				if !ok {
					return Cmd{}, false
				}
				if meta.Linewise {
					return Cmd{
						Kind:  CmdOpChangeLineRegion,
						Reg:   reg,
						Start: start,
						End:   end,
					}, true
				} else {
					return Cmd{
						Kind:      CmdOpChangeRegion,
						Reg:       reg,
						Start:     start,
						End:       end,
						Inclusive: meta.Inclusive,
					}, true
				}
			}
			return Cmd{}, false

		}
	}

	switch op {

	case "yy", "Y":
		return Cmd{Kind: CmdOpCopyLine, Reg: reg, Num: num}, true
	case "yw":
		return Cmd{Kind: CmdOpCopyWord, Reg: reg, Num: num}, true
	case "y$":
		return Cmd{Kind: CmdOpCopyToEnd, Reg: reg, Num: num}, true

	case "p":
		return Cmd{Kind: CmdOpPaste, Reg: reg, Num: num}, true
	case "P":
		return Cmd{Kind: CmdOpPasteBefore, Reg: reg, Num: num}, true

	case "x":
		return Cmd{Kind: CmdOpDelete, Reg: reg, Num: num}, true
	case "X":
		return Cmd{Kind: CmdOpDeleteBefore, Reg: reg, Num: num}, true
	case "dd":
		return Cmd{Kind: CmdOpDeleteLine, Reg: reg, Num: num}, true
	case "dw":
		return Cmd{Kind: CmdOpDeleteWord, Reg: reg, Num: num}, true
	case "d$", "D":
		return Cmd{Kind: CmdOpDeleteToEnd, Reg: reg, Num: num}, true

	case "cc":
		return Cmd{Kind: CmdOpChangeLine, Reg: reg, Num: num}, true
	case "cw":
		return Cmd{Kind: CmdOpChangeWord, Reg: reg, Num: num}, true
	case "C":
		return Cmd{Kind: CmdOpChangeToEnd, Reg: reg, Num: num}, true
	case "s":
		return Cmd{Kind: CmdOpSubst, Reg: reg, Num: num}, true
	case "S":
		return Cmd{Kind: CmdOpSubstLine, Reg: reg, Num: num}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseEdit(
	num int, op string, noSubnum bool, subnum int, mv string, letter rune,
) (Cmd, bool) {
	switch op {

	case "J":
		return Cmd{Kind: CmdEditJoin, Num: num}, true

	case ">>":
		return Cmd{Kind: CmdEditIndent, Num: num}, true
	case "<<":
		return Cmd{Kind: CmdEditOutdent, Num: num}, true

	case ">":
		b := ed.Buf()
		start := b.Loc
		cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
		if ok {
			meta, ok := MoveMetas[cmd.Kind]
			if !ok {
				return Cmd{}, false
			}
			end, ok := ed.RunMove(cmd)
			if !ok {
				return Cmd{}, false
			}
			if meta.Linewise {
				return Cmd{
					Kind:  CmdEditIndentRegion,
					Start: start,
					End:   end,
				}, true
			} else {
				return Cmd{
					Kind:      CmdEditIndentRegion,
					Start:     start,
					End:       end,
					Inclusive: meta.Inclusive,
				}, true
			}
		}
		return Cmd{}, false
	case "<":
		b := ed.Buf()
		start := b.Loc
		cmd, ok := ed.ParseMove(noSubnum, subnum, mv, letter)
		if ok {
			meta, ok := MoveMetas[cmd.Kind]
			if !ok {
				return Cmd{}, false
			}
			end, ok := ed.RunMove(cmd)
			if !ok {
				return Cmd{}, false
			}
			if meta.Linewise {
				return Cmd{
					Kind:  CmdEditOutdentRegion,
					Start: start,
					End:   end,
				}, true
			} else {
				return Cmd{
					Kind:      CmdEditOutdentRegion,
					Start:     start,
					End:       end,
					Inclusive: meta.Inclusive,
				}, true
			}
		}
		return Cmd{}, false

	}

	return Cmd{}, false
}
