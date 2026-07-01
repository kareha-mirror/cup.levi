package editor

func (a *Args) ParseRune() (Cmd, bool) {
	if a.Rune == 0 {
		return Cmd{}, false
	}

	switch a.Op {

	case 'm':
		return Cmd{Kind: Mark, Rune: a.Rune}, true

	case 'r':
		return Cmd{
			Kind: Replace,
			Num:  a.Num,
			Rune: a.Rune,
		}, true

	}

	return Cmd{}, false
}

func (a *Args) ParseView() (Cmd, bool) {
	switch a.Op {

	case 0x06: // Ctrl-F
		return Cmd{Kind: ViewDown, Num: a.Num}, true
	case 0x02: // Ctrl-B
		return Cmd{Kind: ViewUp, Num: a.Num}, true
	case 0x04: // Ctrl-D
		return Cmd{Kind: ViewDownHalf, Num: a.Num}, true
	case 0x15: // Ctrl-U
		return Cmd{Kind: ViewUpHalf, Num: a.Num}, true
	case 0x19: // Ctrl-Y
		return Cmd{Kind: ViewDownLine, Num: a.Num}, true
	case 0x05: // Ctrl-E
		return Cmd{Kind: ViewUpLine, Num: a.Num}, true

	}

	return Cmd{}, false
}

func (a *Args) ParseInsert() (Cmd, bool) {
	switch a.Op {

	case 'i':
		return Cmd{Kind: Insert, Num: a.Num}, true
	case 'a':
		return Cmd{Kind: InsertAfter, Num: a.Num}, true
	case 'I':
		return Cmd{Kind: InsertAfterIndent, Num: a.Num}, true
	case 'A':
		return Cmd{Kind: InsertAfterEnd, Num: a.Num}, true

	case 'o':
		return Cmd{Kind: InsertLine, Num: a.Num}, true
	case 'O':
		return Cmd{Kind: InsertLineAbove, Num: a.Num}, true

	case 'R': // unsupported
		return Cmd{Kind: Overwrite}, true
	}

	return Cmd{}, false
}

func (a *Args) ParseMisc() (Cmd, bool) {
	switch a.Op {

	case 0x07: // Ctrl-G
		return Cmd{Kind: ShowInfo}, true
	case 0x0c: // Ctrl-L
		return Cmd{Kind: Redraw}, true
	case '.':
		return Cmd{Kind: Repeat, Num: a.Num}, true
	case 'u':
		return Cmd{Kind: Undo, Num: a.Num}, true
	case 0x1a: // Ctrl-Z
		return Cmd{Kind: Suspend}, true

	case 0x1e, 0x1f: // Ctrl-^, Ctrl-_
		if a.NoNum {
			return Cmd{Kind: LastBuf}, true
		} else {
			return Cmd{Kind: GoToBuf, Num: a.Num}, true
		}
	}

	return Cmd{}, false
}

func (a *Args) ParseOp() (CmdPair, bool) {
	if a.Mv != 0 {
		switch a.Op {

		case 'y':

			cmd, ok := a.Sub().ParseMove(true)
			if ok {
				return CmdPair{
					Reg: a.Reg,
					Op:  Cmd{Kind: CopyRegion, Num: a.Num},
					Mv:  cmd,
				}, true
			}
			return CmdPair{}, false
		case 'd':
			cmd, ok := a.Sub().ParseMove(true)
			if ok {
				return CmdPair{
					Reg: a.Reg,
					Op:  Cmd{Kind: DeleteRegion, Num: a.Num},
					Mv:  cmd,
				}, true
			}
			return CmdPair{}, false
		case 'c':
			cmd, ok := a.Sub().ParseMove(true)
			if ok {
				return CmdPair{
					Reg: a.Reg,
					Op:  Cmd{Kind: ChangeRegion, Num: a.Num},
					Mv:  cmd,
				}, true
			}
			return CmdPair{}, false

		}
	}

	switch a.Op {

	case 'p':
		return CmdPair{Reg: a.Reg, Op: Cmd{Kind: Paste, Num: a.Num}}, true
	case 'P':
		return CmdPair{
			Reg: a.Reg, Op: Cmd{Kind: PasteBefore, Num: a.Num},
		}, true

	case 'x':
		return CmdPair{Reg: a.Reg, Op: Cmd{Kind: Delete, Num: a.Num}}, true
	case 'X':
		return CmdPair{
			Reg: a.Reg,
			Op:  Cmd{Kind: DeleteBefore, Num: a.Num},
		}, true
	case 'D':
		return CmdPair{
			Reg: a.Reg,
			Op:  Cmd{Kind: DeleteRegion, Num: a.Num},
			Mv:  Cmd{Kind: MoveToEnd, Num: 1},
		}, true

	case 'C':
		return CmdPair{
			Reg: a.Reg,
			Op:  Cmd{Kind: ChangeRegion, Num: a.Num},
			Mv:  Cmd{Kind: MoveToEnd, Num: 1},
		}, true
	case 's':
		return CmdPair{
			Reg: a.Reg,
			Op:  Cmd{Kind: Subst, Num: a.Num},
		}, true
	case 'S':
		return CmdPair{
			Reg: a.Reg,
			Op:  Cmd{Kind: ChangeRegion, Num: a.Num},
			Mv:  Cmd{Kind: MoveHere, Num: 1},
		}, true

	}

	return CmdPair{}, false
}

func (a *Args) ParseEdit() (CmdPair, bool) {
	switch a.Op {

	case 'J':
		return CmdPair{Op: Cmd{Kind: Join, Num: a.Num}}, true

	case '>':
		cmd, ok := a.Sub().ParseMove(true)
		if ok {
			attr, ok := MoveAttrs[cmd.Kind]
			if !ok {
				return CmdPair{}, false
			}
			if attr.Linewise {
				return CmdPair{
					Op: Cmd{Kind: IndentRegion, Num: a.Num},
					Mv: cmd,
				}, true
			} else {
				return CmdPair{
					Op: Cmd{Kind: IndentRegion, Num: a.Num},
					Mv: cmd,
				}, true
			}
		}
		return CmdPair{}, false
	case '<':
		cmd, ok := a.Sub().ParseMove(true)
		if ok {
			attr, ok := MoveAttrs[cmd.Kind]
			if !ok {
				return CmdPair{}, false
			}
			if attr.Linewise {
				return CmdPair{
					Op: Cmd{Kind: OutdentRegion, Num: a.Num},
					Mv: cmd,
				}, true
			} else {
				return CmdPair{
					Op: Cmd{Kind: OutdentRegion, Num: a.Num},
					Mv: cmd,
				}, true
			}
		}
		return CmdPair{}, false

	case 'U':
		return CmdPair{Op: Cmd{Kind: Restore}}, true

	}

	return CmdPair{}, false
}

func (a *Args) ParseCompound() (CmdPair, bool) {
	switch a.Op {

	case ']':
		if a.Mv == 0 {
			return CmdPair{}, false
		}
		if a.Mv != ']' {
			return CmdPair{Op: Cmd{
				Kind: Ring,
				Pat:  "Usage: ]]",
			}}, true
		}
		return CmdPair{Mv: Cmd{Kind: MoveBySection, Num: a.Num}}, true
	case '[':
		if a.Mv == 0 {
			return CmdPair{}, false
		}
		if a.Mv != '[' {
			return CmdPair{Op: Cmd{
				Kind: Ring,
				Pat:  "Usage: [[",
			}}, true
		}
		return CmdPair{Mv: Cmd{Kind: MoveBackwardBySection, Num: a.Num}}, true

	case 'z':
		if a.Mv == 0 {
			return CmdPair{}, false
		}
		switch a.Mv {
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
			return CmdPair{Op: Cmd{
				Kind: Ring,
				Pat:  "Usage: [line]z[window_size][-|.|+|^|<CR>]",
			}}, true
			return CmdPair{}, true
		}

	case 'Z':
		if a.Mv == 0 {
			return CmdPair{}, false
		}
		if a.Mv != 'Z' {
			return CmdPair{Op: Cmd{
				Kind: Ring,
				Pat:  "Usage: ZZ",
			}}, true
			return CmdPair{}, true
		}
		return CmdPair{Op: Cmd{Kind: SaveAndClose}}, true

	}

	return CmdPair{}, false
}
