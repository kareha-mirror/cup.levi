package cmd

func (a Args) compileMove(sub bool) (Cmd, bool) {
	if a.Rune != 0 {
		switch a.Mv {

		case '`':
			if a.Rune == '`' {
				return Cmd{Kind: BackToMark}, true
			} else {
				return Cmd{Kind: MoveToMark, Rune: a.Rune}, true
			}
		case '\'':
			if a.Rune == '\'' {
				return Cmd{Kind: BackToMarkLine}, true
			} else {
				return Cmd{Kind: MoveToMarkLine, Rune: a.Rune}, true
			}

		case 'f':
			return Cmd{
				Kind: Find,
				Num:  a.Num,
				Rune: a.Rune,
			}, true
		case 'F':
			return Cmd{
				Kind: FindBackward,
				Num:  a.Num,
				Rune: a.Rune,
			}, true
		case 't':
			return Cmd{
				Kind: FindBefore,
				Num:  a.Num,
				Rune: a.Rune,
			}, true
		case 'T':
			return Cmd{
				Kind: FindBeforeBackward,
				Num:  a.Num,
				Rune: a.Rune,
			}, true

		}
	}

	switch a.Mv {

	case 'h':
		return Cmd{Kind: MoveLeft, Num: a.Num}, true
	case 'j':
		return Cmd{Kind: MoveDown, Num: a.Num}, true
	//case 'g': // XXX debug
	//	return Cmd{Kind: MoveHere, Num: a.Num}, true
	case 'k':
		return Cmd{Kind: MoveUp, Num: a.Num}, true
	case 'l':
		return Cmd{Kind: MoveRight, Num: a.Num}, true
	case '0': // special
		return Cmd{Kind: MoveToStart}, true
	case '$':
		return Cmd{Kind: MoveToEnd, Num: a.Num}, true
	case '^':
		return Cmd{Kind: MoveToAfterIndent}, true
	case '|':
		return Cmd{Kind: MoveToColumn, Num: a.Num}, true

	case 'w':
		return Cmd{Kind: MoveByWord, Num: a.Num}, true
	//case 'g': // XXX debug
	//	return Cmd{Kind: MoveByChangeWord, Num: a.Num}, true
	case 'g': // XXX debug
		return Cmd{Kind: MoveByDeleteWord, Num: a.Num}, true
	case 'b':
		return Cmd{Kind: MoveBackwardByWord, Num: a.Num}, true
	case 'e':
		return Cmd{Kind: MoveToEndOfWord, Num: a.Num}, true
	case 'W':
		return Cmd{Kind: MoveByLooseWord, Num: a.Num}, true
	case 'B':
		return Cmd{Kind: MoveBackwardByLooseWord, Num: a.Num}, true
	case 'E':
		return Cmd{Kind: MoveToEndOfLooseWord, Num: a.Num}, true

	case '\r', '+':
		return Cmd{Kind: MoveByLine, Num: a.Num}, true
	case '-':
		return Cmd{Kind: MoveBackwardByLine, Num: a.Num}, true
	case 'G':
		if a.Has {
			return Cmd{Kind: MoveToLine, Num: a.Num}, true
		} else {
			return Cmd{Kind: MoveToLastLine}, true
		}

	case ')':
		return Cmd{Kind: MoveBySentence, Num: a.Num}, true
	case '(':
		return Cmd{Kind: MoveBackwardBySentence, Num: a.Num}, true
	case '}':
		return Cmd{Kind: MoveByParagraph, Num: a.Num}, true
	case '{':
		return Cmd{Kind: MoveBackwardByParagraph, Num: a.Num}, true

	// MoveBySection and MoveBackwardBySection are compound

	case 'H':
		if a.Has {
			return Cmd{Kind: MoveToBelowTopOfView, Num: a.Num}, true
		} else {
			return Cmd{Kind: MoveToTopOfView}, true
		}
	case 'M':
		return Cmd{Kind: MoveToMiddleOfView}, true
	case 'L':
		if a.Has {
			return Cmd{Kind: MoveToAboveBottomOfView, Num: a.Num}, true
		} else {
			return Cmd{Kind: MoveToBottomOfView}, true
		}

	// XXX search pattern
	//case '/':
	//	return Cmd{Kind: RepeatSearch}, true
	//	return Cmd{Kind: Search, Pat: pat}, true
	//case '?':
	//	return Cmd{Kind: RepeatBackwardSearch}, true
	//	return Cmd{Kind: SearchBackward, Pat: pat}, true
	case 'n':
		return Cmd{Kind: SearchNext}, true
	case 'N':
		return Cmd{Kind: SearchPrev}, true

	case ';':
		return Cmd{Kind: FindNext, Num: a.Num}, true
	case ',':
		return Cmd{Kind: FindPrev, Num: a.Num}, true

	}

	// Sub Motion Commands
	if sub {
		switch a.Mv {
		case 'y', 'd', 'c', '>', '<':
			return Cmd{Kind: MoveHere, Num: a.Num}, true
		}
	}

	return Cmd{}, false
}
