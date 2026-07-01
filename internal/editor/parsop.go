package editor

func (ed *Editor) ParseRune(num int, op rune, r rune) (Cmd, bool) {
	if r == 0 {
		return Cmd{}, false
	}

	switch op {

	case 'm':
		return Cmd{Kind: Mark, Rune: r}, true

	case 'r':
		return Cmd{
			Kind: Replace,
			Num:  num,
			Rune: r,
		}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseView(num int, op rune) (Cmd, bool) {
	switch op {

	case 0x06: // Ctrl-F
		return Cmd{Kind: ViewDown, Num: num}, true
	case 0x02: // Ctrl-B
		return Cmd{Kind: ViewUp, Num: num}, true
	case 0x04: // Ctrl-D
		return Cmd{Kind: ViewDownHalf, Num: num}, true
	case 0x15: // Ctrl-U
		return Cmd{Kind: ViewUpHalf, Num: num}, true
	case 0x19: // Ctrl-Y
		return Cmd{Kind: ViewDownLine, Num: num}, true
	case 0x05: // Ctrl-E
		return Cmd{Kind: ViewUpLine, Num: num}, true

	case 0x0c: // Ctrl-L
		return Cmd{Kind: Redraw}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseInsert(num int, op rune) (Cmd, bool) {
	switch op {

	case 'i':
		return Cmd{Kind: Insert, Num: num}, true
	case 'a':
		return Cmd{Kind: InsertAfter, Num: num}, true
	case 'I':
		return Cmd{Kind: InsertAfterIndent, Num: num}, true
	case 'A':
		return Cmd{Kind: InsertAfterEnd, Num: num}, true
	case 'R':
		return Cmd{Kind: Overwrite, Num: num}, true

	case 'o':
		return Cmd{Kind: OpenBelow, Num: num}, true
	case 'O':
		return Cmd{Kind: OpenAbove, Num: num}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseMisc(noNum bool, num int, op rune) (Cmd, bool) {
	switch op {

	case 0x07: // Ctrl-G
		return Cmd{Kind: ShowInfo}, true
	case '.':
		return Cmd{Kind: Repeat, Num: num}, true
	case 'u':
		return Cmd{Kind: Undo, Num: num}, true
	case 'U':
		return Cmd{Kind: Restore}, true
	case 0x1a: // Ctrl-Z
		return Cmd{Kind: Suspend}, true

	case 0x1e, 0x1f: // Ctrl-^, Ctrl-_
		if noNum {
			return Cmd{Kind: LastBuf}, true
		} else {
			return Cmd{Kind: GoToBuf, Num: num}, true
		}
	}

	return Cmd{}, false
}

func (ed *Editor) ParseOp(
	reg rune, num int, op rune, noSubnum bool, subnum int, mv rune, r rune,
) (CmdPair, bool) {
	if mv != 0 {
		switch op {

		case 'y':
			cmd, ok := ed.ParseMove(noSubnum, subnum, mv, r, true)
			if ok {
				return CmdPair{
					Reg: reg,
					Op:  Cmd{Kind: CopyRegion, Num: num},
					Mv:  cmd,
				}, true
			}
			return CmdPair{}, false
		case 'd':
			cmd, ok := ed.ParseMove(noSubnum, subnum, mv, r, true)
			if ok {
				return CmdPair{
					Reg: reg,
					Op:  Cmd{Kind: DeleteRegion, Num: num},
					Mv:  cmd,
				}, true
			}
			return CmdPair{}, false
		case 'c':
			cmd, ok := ed.ParseMove(noSubnum, subnum, mv, r, true)
			if ok {
				return CmdPair{
					Reg: reg,
					Op:  Cmd{Kind: ChangeRegion, Num: num},
					Mv:  cmd,
				}, true
			}
			return CmdPair{}, false

		}
	}

	switch op {

	case 'p':
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: Paste, Num: num},
		}, true
	case 'P':
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: PasteBefore, Num: num},
		}, true

	case 'x':
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: Delete, Num: num},
		}, true
	case 'X':
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: DeleteBefore, Num: num},
		}, true
	case 'D':
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: DeleteRegion, Num: num},
			Mv:  Cmd{Kind: MoveToEnd, Num: 1},
		}, true

	case 'C':
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: ChangeRegion, Num: num},
			Mv:  Cmd{Kind: MoveToEnd, Num: 1},
		}, true
	case 's':
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: Subst, Num: num},
		}, true
	case 'S':
		return CmdPair{
			Reg: reg,
			Op:  Cmd{Kind: ChangeRegion, Num: num},
			Mv:  Cmd{Kind: MoveHere, Num: 1},
		}, true

	}

	return CmdPair{}, false
}

func (ed *Editor) ParseEdit(
	num int, op rune, noSubnum bool, subnum int, mv rune, r rune,
) (CmdPair, bool) {
	switch op {

	case 'J':
		return CmdPair{
			Op: Cmd{Kind: Join, Num: num},
		}, true

	case '>':
		cmd, ok := ed.ParseMove(noSubnum, subnum, mv, r, true)
		if ok {
			attr, ok := MoveAttrs[cmd.Kind]
			if !ok {
				return CmdPair{}, false
			}
			if attr.Linewise {
				return CmdPair{
					Op: Cmd{Kind: IndentRegion, Num: num},
					Mv: cmd,
				}, true
			} else {
				return CmdPair{
					Op: Cmd{Kind: IndentRegion, Num: num},
					Mv: cmd,
				}, true
			}
		}
		return CmdPair{}, false
	case '<':
		cmd, ok := ed.ParseMove(noSubnum, subnum, mv, r, true)
		if ok {
			attr, ok := MoveAttrs[cmd.Kind]
			if !ok {
				return CmdPair{}, false
			}
			if attr.Linewise {
				return CmdPair{
					Op: Cmd{Kind: OutdentRegion, Num: num},
					Mv: cmd,
				}, true
			} else {
				return CmdPair{
					Op: Cmd{Kind: OutdentRegion, Num: num},
					Mv: cmd,
				}, true
			}
		}
		return CmdPair{}, false

	}

	return CmdPair{}, false
}

func (ed *Editor) ParseCompound(
	num int, op rune, noSubnum bool, subnum int, mv rune, r rune,
) (CmdPair, bool) {
	switch op {

	case ']':
		if mv == 0 {
			return CmdPair{}, false
		}
		if mv != ']' {
			ed.Ring("Usage: ]]")
			return CmdPair{}, true
		}
		return CmdPair{Mv: Cmd{Kind: MoveBySection, Num: num}}, true
	case '[':
		if mv == 0 {
			return CmdPair{}, false
		}
		if mv != '[' {
			ed.Ring("Usage: [[")
			return CmdPair{}, true
		}
		return CmdPair{Mv: Cmd{Kind: MoveBackwardBySection, Num: num}}, true

	case 'z':
		if mv == 0 {
			return CmdPair{}, false
		}
		switch mv {
		case '\r':
			return CmdPair{Op: Cmd{Kind: ViewToTop}}, true
		case '.':
			return CmdPair{Op: Cmd{Kind: ViewToMiddle}}, true
		case '-':
			return CmdPair{Op: Cmd{Kind: ViewToBottom}}, true

		case 'j':
			return CmdPair{Op: Cmd{Kind: NextBuf}}, true
		case 'k':
			return CmdPair{Op: Cmd{Kind: PrevBuf}}, true

		default:
			ed.Ring("Usage: [line]z[window_size][-|.|+|^|<CR>]")
			return CmdPair{}, true
		}

	case 'Z':
		if mv == 0 {
			return CmdPair{}, false
		}
		if mv != 'Z' {
			ed.Ring("Usage: ZZ")
			return CmdPair{}, true
		}
		return CmdPair{Op: Cmd{Kind: SaveAndClose}}, true

	}

	return CmdPair{}, false
}
