package editor

func (ed *Editor) ParseMoveLetter(
	num int, op string, letter rune,
) (Cmd, bool) {
	if letter == 0 {
		return Cmd{}, false
	}
	switch op {

	//
	// Linewise
	//

	case "'":
		if letter == '\'' {
			return Cmd{Kind: CmdMoveBackToMarkLine}, true
		} else {
			return Cmd{Kind: CmdMoveToMarkLine, Letter: letter}, true
		}

	//
	// Runewise
	//

	case "`":
		if letter == '`' {
			return Cmd{Kind: CmdMoveBackToMark}, true
		} else {
			return Cmd{Kind: CmdMoveToMark, Letter: letter}, true
		}

	case "f":
		return Cmd{
			Kind:   CmdMoveFindForward,
			Num:    num,
			Letter: letter,
		}, true
	case "F":
		return Cmd{
			Kind:   CmdMoveFindBackward,
			Num:    num,
			Letter: letter,
		}, true
	case "t":
		return Cmd{
			Kind:   CmdMoveFindBeforeForward,
			Num:    num,
			Letter: letter,
		}, true
	case "T":
		return Cmd{
			Kind:   CmdMoveFindBeforeBackward,
			Num:    num,
			Letter: letter,
		}, true

	}

	return Cmd{}, false
}

func (ed *Editor) ParseMove(
	noNum bool, num int, mv string, letter rune,
) (Cmd, bool) {
	switch mv {

	//
	// Linewise
	//

	case "j":
		return Cmd{Kind: CmdMoveDown, Num: num}, true
	case "k":
		return Cmd{Kind: CmdMoveUp, Num: num}, true

	case "\r", "+":
		return Cmd{Kind: CmdMoveByLine, Num: num}, true
	case "-":
		return Cmd{Kind: CmdMoveBackwardByLine, Num: num}, true
	case "G":
		if noNum {
			return Cmd{Kind: CmdMoveToLastLine}, true
		} else {
			return Cmd{Kind: CmdMoveToLine, Num: num}, true
		}

	case ")":
		return Cmd{Kind: CmdMoveBySentence, Num: num}, true
	case "(":
		return Cmd{Kind: CmdMoveBackwardBySentence, Num: num}, true
	case "}":
		return Cmd{Kind: CmdMoveByParagraph, Num: num}, true
	case "{":
		return Cmd{Kind: CmdMoveBackwardByParagraph, Num: num}, true
	case "]]":
		return Cmd{Kind: CmdMoveBySection, Num: num}, true
	case "[[":
		return Cmd{Kind: CmdMoveBackwardBySection, Num: num}, true

	case "H":
		if noNum {
			return Cmd{Kind: CmdMoveToTopOfView}, true
		} else {
			return Cmd{Kind: CmdMoveToBelowTopOfView, Num: num}, true
		}
	case "M":
		return Cmd{Kind: CmdMoveToMiddleOfView}, true
	case "L":
		if noNum {
			return Cmd{Kind: CmdMoveToBottomOfView}, true
		} else {
			return Cmd{Kind: CmdMoveToAboveBottomOfView, Num: num}, true
		}

	//
	// Runewise
	//

	case "h":
		return Cmd{Kind: CmdMoveLeft, Num: num}, true
	case "l":
		return Cmd{Kind: CmdMoveRight, Num: num}, true

	case "0": // special
		return Cmd{Kind: CmdMoveToStart}, true
	case "$":
		return Cmd{Kind: CmdMoveToEnd}, true
	case "^":
		return Cmd{Kind: CmdMoveToNonBlank}, true
	case "|":
		return Cmd{Kind: CmdMoveToColumn, Num: num}, true

	case "w":
		return Cmd{Kind: CmdMoveByWord, Num: num}, true
	case "g": // XXX debug
		return Cmd{Kind: CmdMoveByWordEx, Num: num}, true
	case "b":
		return Cmd{Kind: CmdMoveBackwardByWord, Num: num}, true
	case "e":
		return Cmd{Kind: CmdMoveToEndOfWord, Num: num}, true
	case "W":
		return Cmd{Kind: CmdMoveByLooseWord, Num: num}, true
	case "B":
		return Cmd{Kind: CmdMoveBackwardByLooseWord, Num: num}, true
	case "E":
		return Cmd{Kind: CmdMoveToEndOfLooseWord, Num: num}, true

	case ";":
		return Cmd{Kind: CmdMoveFindNextMatch, Num: num}, true
	case ",":
		return Cmd{Kind: CmdMoveFindPrevMatch, Num: num}, true

	}

	//
	// Additions
	//

	cmd, ok := ed.ParseMoveLetter(num, mv, letter)
	if ok {
		return cmd, true
	}
	return ed.ParseSearch(mv, "") // XXX pat
}

func (ed *Editor) ParseSearch(op string, pat string) (Cmd, bool) {
	switch op {

	case "/":
		if pat == "" {
			return Cmd{Kind: CmdMoveSearchRepeatForward}, true
		} else {
			return Cmd{Kind: CmdMoveSearchForward, Pat: pat}, true
		}
	case "?":
		if pat == "" {
			return Cmd{Kind: CmdMoveSearchRepeatBackward}, true
		} else {
			return Cmd{Kind: CmdMoveSearchBackward, Pat: pat}, true
		}

	case "n":
		return Cmd{Kind: CmdMoveSearchNextMatch}, true
	case "N":
		return Cmd{Kind: CmdMoveSearchPrevMatch}, true

	}

	return Cmd{}, false
}
