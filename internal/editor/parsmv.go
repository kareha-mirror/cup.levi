package editor

func (ed *Editor) ParseMoveRune(num int, op rune, r rune) (Cmd, bool) {
	if r == 0 {
		return Cmd{}, false
	}
	switch op {

	//
	// Linewise
	//

	case '\'':
		if r == '\'' {
			return Cmd{Kind: BackToMarkLine}, true
		} else {
			return Cmd{Kind: MoveToMarkLine, Rune: r}, true
		}

	//
	// Runewise
	//

	case '`':
		if r == '`' {
			return Cmd{Kind: BackToMark}, true
		} else {
			return Cmd{Kind: MoveToMark, Rune: r}, true
		}

	case 'f':
		return Cmd{
			Kind: Find,
			Num:  num,
			Rune: r,
		}, true
	case 'F':
		return Cmd{
			Kind: FindBackward,
			Num:  num,
			Rune: r,
		}, true
	case 't':
		return Cmd{
			Kind: FindBefore,
			Num:  num,
			Rune: r,
		}, true
	case 'T':
		return Cmd{
			Kind: FindBeforeBackward,
			Num:  num,
			Rune: r,
		}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseMove(
	noNum bool, num int, mv rune, r rune, meta bool,
) (Cmd, bool) {
	switch mv {

	//
	// Linewise
	//

	case 'j':
		return Cmd{Kind: MoveDown, Num: num}, true
	//case 'g': // XXX debug
	//	return Cmd{Kind: MoveHere, Num: num}, true
	case 'k':
		return Cmd{Kind: MoveUp, Num: num}, true

	case '\r', '+':
		return Cmd{Kind: MoveByLine, Num: num}, true
	case '-':
		return Cmd{Kind: MoveBackwardByLine, Num: num}, true
	case 'G':
		if noNum {
			return Cmd{Kind: MoveToLastLine}, true
		} else {
			return Cmd{Kind: MoveToLine, Num: num}, true
		}

	case ')':
		return Cmd{Kind: MoveBySentence, Num: num}, true
	case '(':
		return Cmd{Kind: MoveBackwardBySentence, Num: num}, true
	case '}':
		return Cmd{Kind: MoveByParagraph, Num: num}, true
	case '{':
		return Cmd{Kind: MoveBackwardByParagraph, Num: num}, true

	case 'H':
		if noNum {
			return Cmd{Kind: MoveToTopOfView}, true
		} else {
			return Cmd{Kind: MoveToBelowTopOfView, Num: num}, true
		}
	case 'M':
		return Cmd{Kind: MoveToMiddleOfView}, true
	case 'L':
		if noNum {
			return Cmd{Kind: MoveToBottomOfView}, true
		} else {
			return Cmd{Kind: MoveToAboveBottomOfView, Num: num}, true
		}

	//
	// Runewise
	//

	case 'h':
		return Cmd{Kind: MoveLeft, Num: num}, true
	case 'l':
		return Cmd{Kind: MoveRight, Num: num}, true

	case '0': // special
		return Cmd{Kind: MoveToStart}, true
	case '$':
		return Cmd{Kind: MoveToEnd, Num: num}, true
	case '^':
		return Cmd{Kind: MoveToAfterIndent}, true
	case '|':
		return Cmd{Kind: MoveToColumn, Num: num}, true

	case 'w':
		return Cmd{Kind: MoveByWord, Num: num}, true
	//case 'g': // XXX debug
	//	return Cmd{Kind: MoveByChangeWord, Num: num}, true
	case 'g': // XXX debug
		return Cmd{Kind: MoveByDeleteWord, Num: num}, true
	case 'b':
		return Cmd{Kind: MoveBackwardByWord, Num: num}, true
	case 'e':
		return Cmd{Kind: MoveToEndOfWord, Num: num}, true
	case 'W':
		return Cmd{Kind: MoveByLooseWord, Num: num}, true
	case 'B':
		return Cmd{Kind: MoveBackwardByLooseWord, Num: num}, true
	case 'E':
		return Cmd{Kind: MoveToEndOfLooseWord, Num: num}, true

	case ';':
		return Cmd{Kind: FindNext, Num: num}, true
	case ',':
		return Cmd{Kind: FindPrev, Num: num}, true

	}

	//
	// Additions
	//

	cmd, ok := ed.ParseMoveRune(num, mv, r)
	if ok {
		return cmd, true
	}
	cmd, ok = ed.ParseSearch(mv, "") // XXX pat
	if ok {
		return cmd, true
	}

	//
	// Meta Motion Commands
	//

	if meta {
		switch mv {
		case 'y', 'd', 'c', '>', '<':
			return Cmd{Kind: MoveHere, Num: num}, true
		}
	}

	return Cmd{}, false
}

func (ed *Editor) ParseSearch(op rune, pat string) (Cmd, bool) {
	switch op {

	case '/':
		if pat == "" {
			return Cmd{Kind: RepeatSearch}, true
		} else {
			return Cmd{Kind: Search, Pat: pat}, true
		}
	case '?':
		if pat == "" {
			return Cmd{Kind: RepeatBackwardSearch}, true
		} else {
			return Cmd{Kind: SearchBackward, Pat: pat}, true
		}

	case 'n':
		return Cmd{Kind: SearchNext}, true
	case 'N':
		return Cmd{Kind: SearchPrev}, true

	}

	return Cmd{}, false
}
