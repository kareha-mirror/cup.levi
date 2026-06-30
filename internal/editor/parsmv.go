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
			return Cmd{Kind: BackToMarkLine}, true
		} else {
			return Cmd{Kind: MoveToMarkLine, Ltr: letter}, true
		}

	//
	// Runewise
	//

	case "`":
		if letter == '`' {
			return Cmd{Kind: BackToMark}, true
		} else {
			return Cmd{Kind: MoveToMark, Ltr: letter}, true
		}

	case "f":
		return Cmd{
			Kind: FindForward,
			Num:  num,
			Ltr:  letter,
		}, true
	case "F":
		return Cmd{
			Kind: FindBackward,
			Num:  num,
			Ltr:  letter,
		}, true
	case "t":
		return Cmd{
			Kind: FindBeforeForward,
			Num:  num,
			Ltr:  letter,
		}, true
	case "T":
		return Cmd{
			Kind: FindBeforeBackward,
			Num:  num,
			Ltr:  letter,
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
		return Cmd{Kind: MoveDown, Num: num}, true
	case "k":
		return Cmd{Kind: MoveUp, Num: num}, true

	case "\r", "+":
		return Cmd{Kind: MoveByLine, Num: num}, true
	case "-":
		return Cmd{Kind: MoveBackwardByLine, Num: num}, true
	case "G":
		if noNum {
			return Cmd{Kind: MoveToLastLine}, true
		} else {
			return Cmd{Kind: MoveToLine, Num: num}, true
		}

	case ")":
		return Cmd{Kind: MoveBySentence, Num: num}, true
	case "(":
		return Cmd{Kind: MoveBackwardBySentence, Num: num}, true
	case "}":
		return Cmd{Kind: MoveByParagraph, Num: num}, true
	case "{":
		return Cmd{Kind: MoveBackwardByParagraph, Num: num}, true
	case "]]":
		return Cmd{Kind: MoveBySection, Num: num}, true
	case "[[":
		return Cmd{Kind: MoveBackwardBySection, Num: num}, true

	case "H":
		if noNum {
			return Cmd{Kind: MoveToTopOfView}, true
		} else {
			return Cmd{Kind: MoveToBelowTopOfView, Num: num}, true
		}
	case "M":
		return Cmd{Kind: MoveToMiddleOfView}, true
	case "L":
		if noNum {
			return Cmd{Kind: MoveToBottomOfView}, true
		} else {
			return Cmd{Kind: MoveToAboveBottomOfView, Num: num}, true
		}

	//
	// Runewise
	//

	case "h":
		return Cmd{Kind: MoveLeft, Num: num}, true
	case "l":
		return Cmd{Kind: MoveRight, Num: num}, true

	case "0": // special
		return Cmd{Kind: MoveToStart}, true
	case "$":
		return Cmd{Kind: MoveToEnd}, true
	case "^":
		return Cmd{Kind: MoveToNonBlank}, true
	case "|":
		return Cmd{Kind: MoveToColumn, Num: num}, true

	case "w":
		return Cmd{Kind: MoveByWord, Num: num}, true
	case "g": // XXX debug
		return Cmd{Kind: MoveByChangeWord, Num: num}, true
	case "b":
		return Cmd{Kind: MoveBackwardByWord, Num: num}, true
	case "e":
		return Cmd{Kind: MoveToEndOfWord, Num: num}, true
	case "W":
		return Cmd{Kind: MoveByLooseWord, Num: num}, true
	case "B":
		return Cmd{Kind: MoveBackwardByLooseWord, Num: num}, true
	case "E":
		return Cmd{Kind: MoveToEndOfLooseWord, Num: num}, true

	case ";":
		return Cmd{Kind: FindNext, Num: num}, true
	case ",":
		return Cmd{Kind: FindPrev, Num: num}, true

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
			return Cmd{Kind: RepeatSearchForward}, true
		} else {
			return Cmd{Kind: SearchForward, Pat: pat}, true
		}
	case "?":
		if pat == "" {
			return Cmd{Kind: RepeatSearchBackward}, true
		} else {
			return Cmd{Kind: SearchBackward, Pat: pat}, true
		}

	case "n":
		return Cmd{Kind: SearchNext}, true
	case "N":
		return Cmd{Kind: SearchPrev}, true

	}

	return Cmd{}, false
}
