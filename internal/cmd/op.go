package cmd

// Simple Commands
func (a Args) compileOp() (Pair, bool) {
	//
	// Rune Commands
	//

	if a.Rune != 0 {
		switch a.Op {

		// Mark Commands
		case 'm':
			return Pair{Op: Cmd{Kind: Mark, Rune: a.Rune}}, true

		// Edit Commands
		case 'r':
			return Pair{Op: Cmd{Kind: Replace, Num: a.Num, Rune: a.Rune}}, true

		}
	}

	switch a.Op {

	//
	// Insert Commands
	//

	case 'i':
		return Pair{Op: Cmd{Kind: Insert, Num: a.Num}}, true
	case 'a':
		return Pair{Op: Cmd{Kind: InsertAfter, Num: a.Num}}, true
	case 'I':
		return Pair{Op: Cmd{Kind: InsertAfterIndent, Num: a.Num}}, true
	case 'A':
		return Pair{Op: Cmd{Kind: InsertAfterEnd, Num: a.Num}}, true

	case 'o':
		return Pair{Op: Cmd{Kind: InsertLine, Num: a.Num}}, true
	case 'O':
		return Pair{Op: Cmd{Kind: InsertLineAbove, Num: a.Num}}, true

	case 'c':
		mv, ok := a.sub().compileMove(true)
		if ok {
			return Pair{
				Reg: a.Reg,
				Op:  Cmd{Kind: ChangeRegion, Num: a.Num},
				Mv:  mv,
			}, true
		}
		return Pair{}, false
	case 'C':
		return Pair{
			Reg: a.Reg,
			Op:  Cmd{Kind: ChangeRegion, Num: a.Num},
			Mv:  Cmd{Kind: MoveToEnd, Num: 1},
		}, true
	case 's':
		return Pair{
			Reg: a.Reg,
			Op:  Cmd{Kind: Subst, Num: a.Num},
		}, true
	case 'S':
		return Pair{
			Reg: a.Reg,
			Op:  Cmd{Kind: ChangeRegion, Num: a.Num},
			Mv:  Cmd{Kind: MoveHere, Num: 1},
		}, true

	case 'R': // unsupported
		return Pair{Op: Cmd{Kind: Overwrite}}, true

	//
	// Edit Commands
	//

	case 'p':
		return Pair{Reg: a.Reg, Op: Cmd{Kind: Paste, Num: a.Num}}, true
	case 'P':
		return Pair{
			Reg: a.Reg, Op: Cmd{Kind: PasteBefore, Num: a.Num},
		}, true

	case 'x':
		return Pair{Reg: a.Reg, Op: Cmd{Kind: Delete, Num: a.Num}}, true
	case 'X':
		return Pair{
			Reg: a.Reg,
			Op:  Cmd{Kind: DeleteBefore, Num: a.Num},
		}, true
	case 'd':
		mv, ok := a.sub().compileMove(true)
		if ok {
			return Pair{
				Reg: a.Reg,
				Op:  Cmd{Kind: DeleteRegion, Num: a.Num},
				Mv:  mv,
			}, true
		}
		return Pair{}, false
	case 'D':
		return Pair{
			Reg: a.Reg,
			Op:  Cmd{Kind: DeleteRegion, Num: a.Num},
			Mv:  Cmd{Kind: MoveToEnd, Num: 1},
		}, true

	case 'J':
		return Pair{Op: Cmd{Kind: Join, Num: a.Num}}, true
	case '>':
		mv, ok := a.sub().compileMove(true)
		if ok {
			attr, ok := AttrOf[mv.Kind]
			if !ok {
				return Pair{}, false
			}
			if attr.Linewise {
				return Pair{
					Op: Cmd{Kind: IndentRegion, Num: a.Num},
					Mv: mv,
				}, true
			} else {
				return Pair{
					Op: Cmd{Kind: IndentRegion, Num: a.Num},
					Mv: mv,
				}, true
			}
		}
		return Pair{}, false
	case '<':
		mv, ok := a.sub().compileMove(true)
		if ok {
			attr, ok := AttrOf[mv.Kind]
			if !ok {
				return Pair{}, false
			}
			if attr.Linewise {
				return Pair{
					Op: Cmd{Kind: OutdentRegion, Num: a.Num},
					Mv: mv,
				}, true
			} else {
				return Pair{
					Op: Cmd{Kind: OutdentRegion, Num: a.Num},
					Mv: mv,
				}, true
			}
		}
		return Pair{}, false

	case 'U':
		return Pair{Op: Cmd{Kind: Restore}}, true

	//
	// Copy Commands
	//

	case 'y':
		mv, ok := a.sub().compileMove(true)
		if ok {
			return Pair{
				Reg: a.Reg,
				Op:  Cmd{Kind: CopyRegion, Num: a.Num},
				Mv:  mv,
			}, true
		}
		return Pair{}, false

	//
	// View Commands
	//

	case 0x06: // Ctrl-F
		return Pair{Op: Cmd{Kind: ViewDown, Num: a.Num}}, true
	case 0x02: // Ctrl-B
		return Pair{Op: Cmd{Kind: ViewUp, Num: a.Num}}, true
	case 0x04: // Ctrl-D
		return Pair{Op: Cmd{Kind: ViewDownHalf, Num: a.Num}}, true
	case 0x15: // Ctrl-U
		return Pair{Op: Cmd{Kind: ViewUpHalf, Num: a.Num}}, true
	case 0x19: // Ctrl-Y
		return Pair{Op: Cmd{Kind: ViewDownLine, Num: a.Num}}, true
	case 0x05: // Ctrl-E
		return Pair{Op: Cmd{Kind: ViewUpLine, Num: a.Num}}, true

	//
	// Miscellaneous Commands
	//

	case 0x07: // Ctrl-G
		return Pair{Op: Cmd{Kind: ShowInfo}}, true
	case 0x0c: // Ctrl-L
		return Pair{Op: Cmd{Kind: Redraw}}, true
	case '.':
		return Pair{Op: Cmd{Kind: Repeat, Num: a.Num}}, true
	case 'u':
		return Pair{Op: Cmd{Kind: Undo, Num: a.Num}}, true
	case 'Z':
		if a.Mv == 0 {
			return Pair{}, false
		}
		if a.Mv != 'Z' {
			return Pair{Op: Cmd{Kind: Ring, Pat: "Usage: ZZ"}}, true
		}
		return Pair{Op: Cmd{Kind: SaveAndClose}}, true
	case 0x1a: // Ctrl-Z
		return Pair{Op: Cmd{Kind: Suspend}}, true

	//
	// Select Current Buffer
	//

	case 0x1e, 0x1f: // Ctrl-^, Ctrl-_
		if a.NoNum {
			return Pair{Op: Cmd{Kind: LastBuf}}, true
		} else {
			return Pair{Op: Cmd{Kind: GoToBuf, Num: a.Num}}, true
		}

	//
	// z Commands
	//

	case 'z':
		if a.Mv == 0 {
			return Pair{}, false
		}
		switch a.Mv {
		case '\r':
			return Pair{Op: Cmd{Kind: ViewToTop}}, true
		case '.':
			return Pair{Op: Cmd{Kind: ViewToMiddle}}, true
		case '-':
			return Pair{Op: Cmd{Kind: ViewToBottom}}, true

		case 'j':
			return Pair{Op: Cmd{Kind: NextBuf}}, true
		case 'k':
			return Pair{Op: Cmd{Kind: PrevBuf}}, true

		default:
			return Pair{Op: Cmd{
				Kind: Ring,
				Pat:  "Usage: [line]z[window_size][-|.|+|^|<CR>]",
			}}, true
		}

	//
	// Motion Commands // XXX
	//

	case ']':
		if a.Mv == 0 {
			return Pair{}, false
		}
		if a.Mv != ']' {
			return Pair{Op: Cmd{
				Kind: Ring,
				Pat:  "Usage: ]]",
			}}, true
		}
		return Pair{Mv: Cmd{Kind: MoveBySection, Num: a.Num}}, true
	case '[':
		if a.Mv == 0 {
			return Pair{}, false
		}
		if a.Mv != '[' {
			return Pair{Op: Cmd{
				Kind: Ring,
				Pat:  "Usage: [[",
			}}, true
		}
		return Pair{Mv: Cmd{Kind: MoveBackwardBySection, Num: a.Num}}, true

	}
	return Pair{}, false
}
