package editor

import (
	"strconv"
	"strings"
)

type Parser struct {
	buf   []rune
	cache string
}

const maxParserLen = 256

func NewParser() *Parser {
	return &Parser{
		buf:   make([]rune, 0),
		cache: "",
	}
}

func (p *Parser) String() string {
	b := strings.Builder{}
	for _, r := range p.buf {
		switch r {
		case '\f':
			b.WriteString("(Ctrl-L)")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (p *Parser) InsertRune(r rune) {
	p.buf = append(p.buf, r)
	p.cache = p.String()
}

func (p *Parser) Clear() {
	if len(p.buf) > maxParserLen {
		p.buf = make([]rune, 0)
	} else {
		p.buf = p.buf[:0]
	}
}

func (p *Parser) Cache() string {
	return p.cache
}

func (p *Parser) ClearAll() {
	p.Clear()
	p.cache = ""
}

var letterOpSet = map[rune]struct{}{
	'm': {},
	'r': {},
}

var letterMoveSet = map[rune]struct{}{
	'\'': {},
	'`':  {},
	'f':  {},
	'F':  {},
	't':  {},
	'T':  {},
}

func (ed *Editor) ParseMoveLetterLines(num int, op string, letter rune) (Cmd, bool) {
	switch op {
	case "'":
		if letter == 0 {
			return Cmd{}, false
		}
		if letter == '\'' {
			return Cmd{Kind: CmdMoveBackToMarkLine}, true
		} else {
			return Cmd{
				Kind:   CmdMoveToMarkLine,
				Letter: letter,
			}, true
		}
	}

	return Cmd{}, false
}

func (ed *Editor) ParseMoveLetterRunes(num int, op string, letter rune) (Cmd, bool) {
	switch op {
	case "f":
		if letter == 0 {
			return Cmd{}, false
		}
		return Cmd{
			Kind:   CmdMoveFindForward,
			Num:    num,
			Letter: letter,
		}, true
	case "F":
		if letter == 0 {
			return Cmd{}, false
		}
		return Cmd{
			Kind:   CmdMoveFindBackward,
			Num:    num,
			Letter: letter,
		}, true
	case "t":
		if letter == 0 {
			return Cmd{}, false
		}
		return Cmd{
			Kind:   CmdMoveFindBeforeForward,
			Num:    num,
			Letter: letter,
		}, true
	case "T":
		if letter == 0 {
			return Cmd{}, false
		}
		return Cmd{
			Kind:   CmdMoveFindBeforeBackward,
			Num:    num,
			Letter: letter,
		}, true
	case ";":
		return Cmd{
			Kind: CmdMoveFindNextMatch,
			Num:  num,
		}, true
	case ",":
		return Cmd{
			Kind: CmdMoveFindPrevMatch,
			Num:  num,
		}, true
	}

	switch op {
	case "`":
		if letter == 0 {
			return Cmd{}, false
		}
		if letter == '`' {
			return Cmd{Kind: CmdMoveBackToMark}, true
		} else {
			return Cmd{
				Kind:   CmdMoveToMark,
				Letter: letter,
			}, true
		}
	}

	return Cmd{}, false
}

func (ed *Editor) ParseMoveLines(noNum bool, num int, mv string, letter rune) (Cmd, bool) {
	switch mv {
	case "j":
		return Cmd{
			Kind: CmdMoveDown,
			Num:  num,
		}, true
	case "k":
		return Cmd{
			Kind: CmdMoveUp,
			Num:  num,
		}, true

	case "\r", "+":
		return Cmd{
			Kind: CmdMoveByLine,
			Num:  num,
		}, true
	case "-":
		return Cmd{
			Kind: CmdMoveBackwardByLine,
			Num:  num,
		}, true
	case "G":
		if noNum {
			return Cmd{Kind: CmdMoveToLastLine}, true
		} else {
			return Cmd{
				Kind: CmdMoveToLine,
				Num:  num,
			}, true
		}

	case ")":
		return Cmd{
			Kind: CmdMoveBySentence,
			Num:  num,
		}, true
	case "(":
		return Cmd{
			Kind: CmdMoveBackwardBySentence,
			Num:  num,
		}, true
	case "}":
		return Cmd{
			Kind: CmdMoveByParagraph,
			Num:  num,
		}, true
	case "{":
		return Cmd{
			Kind: CmdMoveBackwardByParagraph,
			Num:  num,
		}, true
	case "]]":
		return Cmd{
			Kind: CmdMoveBySection,
			Num:  num,
		}, true
	case "[[":
		return Cmd{
			Kind: CmdMoveBackwardBySection,
			Num:  num,
		}, true

	case "H":
		if noNum {
			return Cmd{Kind: CmdMoveToTopOfView}, true
		} else {
			return Cmd{
				Kind: CmdMoveToBelowTopOfView,
				Num:  num,
			}, true
		}
	case "M":
		return Cmd{Kind: CmdMoveToMiddleOfView}, true
	case "L":
		if noNum {
			return Cmd{Kind: CmdMoveToBottomOfView}, true
		} else {
			return Cmd{
				Kind: CmdMoveToAboveBottomOfView,
				Num:  num,
			}, true
		}
	}

	return ed.ParseMoveLetterLines(num, mv, letter)
}

func (ed *Editor) ParseMoveRunes(noNum bool, num int, mv string, letter rune) (Cmd, bool) {
	switch mv {
	case "h":
		return Cmd{
			Kind: CmdMoveLeft,
			Num:  num,
		}, true
	case "l":
		return Cmd{
			Kind: CmdMoveRight,
			Num:  num,
		}, true

	case "0": // special
		return Cmd{Kind: CmdMoveToStart}, true
	case "$":
		return Cmd{Kind: CmdMoveToEnd}, true
	case "^":
		return Cmd{Kind: CmdMoveToNonBlank}, true
	case "|":
		return Cmd{
			Kind: CmdMoveToColumn,
			Num:  num,
		}, true

	case "w":
		return Cmd{
			Kind: CmdMoveByWord,
			Num:  num,
		}, true
	case "g": // XXX debug
		return Cmd{
			Kind: CmdMoveByWordEx,
			Num:  num,
		}, true
	case "b":
		return Cmd{
			Kind: CmdMoveBackwardByWord,
			Num:  num,
		}, true
	case "e":
		return Cmd{
			Kind: CmdMoveToEndOfWord,
			Num:  num,
		}, true
	case "W":
		return Cmd{
			Kind: CmdMoveByLooseWord,
			Num:  num,
		}, true
	case "B":
		return Cmd{
			Kind: CmdMoveBackwardByLooseWord,
			Num:  num,
		}, true
	case "E":
		return Cmd{
			Kind: CmdMoveToEndOfLooseWord,
			Num:  num,
		}, true
	}

	cmd, ok := ed.ParseMoveLetterRunes(num, mv, letter)
	if ok {
		return cmd, true
	}
	return ed.ParseSearch(mv, "") // XXX pat
}

func (ed *Editor) ParseLetter(num int, op string, letter rune) (Cmd, bool) {
	if letter == 0 {
		return Cmd{}, false
	}

	switch op {
	case "m":
		return Cmd{
			Kind:   CmdMarkSet,
			Letter: letter,
		}, true
	}

	switch op {
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
		return Cmd{
			Kind: CmdViewDown,
			Num:  num,
		}, true
	case "\x02": // Ctrl-B
		return Cmd{
			Kind: CmdViewUp,
			Num:  num,
		}, true
	case "\x04": // Ctrl-D
		return Cmd{
			Kind: CmdViewDownHalf,
			Num:  num,
		}, true
	case "\x15": // Ctrl-U
		return Cmd{
			Kind: CmdViewUpHalf,
			Num:  num,
		}, true
	case "\x19": // Ctrl-Y
		return Cmd{
			Kind: CmdViewDownLine,
			Num:  num,
		}, true
	case "\x05": // Ctrl-E
		return Cmd{
			Kind: CmdViewUpLine,
			Num:  num,
		}, true

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

func (ed *Editor) ParseSearch(op string, pat string) (Cmd, bool) {
	switch op {
	case "/":
		if pat == "" {
			return Cmd{Kind: CmdMoveSearchRepeatForward}, true
		} else {
			return Cmd{
				Kind: CmdMoveSearchForward,
				Pat:  pat,
			}, true
		}
	case "?":
		if pat == "" {
			return Cmd{Kind: CmdMoveSearchRepeatBackward}, true
		} else {
			return Cmd{
				Kind: CmdMoveSearchBackward,
				Pat:  pat,
			}, true
		}
	case "n":
		return Cmd{Kind: CmdMoveSearchNextMatch}, true
	case "N":
		return Cmd{Kind: CmdMoveSearchPrevMatch}, true
	}

	return Cmd{}, false
}

func (ed *Editor) ParseInsert(num int, op string) (Cmd, bool) {
	switch op {
	case "i":
		return Cmd{
			Kind: CmdInsertBefore,
			Num:  num,
		}, true
	case "a":
		return Cmd{
			Kind: CmdInsertAfter,
			Num:  num,
		}, true
	case "I":
		return Cmd{
			Kind: CmdInsertBeforeNonBlank,
			Num:  num,
		}, true
	case "A":
		return Cmd{
			Kind: CmdInsertAfterEnd,
			Num:  num,
		}, true
	case "R":
		return Cmd{
			Kind: CmdInsertOverwrite,
			Num:  num,
		}, true

	case "o":
		return Cmd{
			Kind: CmdInsertOpenBelow,
			Num:  num,
		}, true
	case "O":
		return Cmd{
			Kind: CmdInsertOpenAbove,
			Num:  num,
		}, true
	}

	return Cmd{}, false
}

func (ed *Editor) ParseMisc(num int, op string) (Cmd, bool) {
	switch op {
	case "\x07": // Ctrl-G
		return Cmd{Kind: CmdMiscShowInfo}, true
	case ".":
		return Cmd{
			Kind: CmdMiscRepeat,
			Num:  num,
		}, true
	case "u":
		return Cmd{
			Kind: CmdMiscUndo,
			Num:  num,
		}, true
	case "U":
		return Cmd{Kind: CmdMiscRestore}, true
	case "ZZ":
		return Cmd{Kind: CmdMiscSaveAndQuit}, true
	case "\x1a": // Ctrl-Z
		return Cmd{Kind: CmdMiscSuspend}, true
	}

	return Cmd{}, false
}

func (ed *Editor) ParseOp(reg rune, num int, op string, noSubnum bool, subnum int, mv string, letter rune) (Cmd, bool) {
	if mv != "" {
		switch op {
		case "y":
			b := ed.Buf()
			start := b.Loc
			cmd, ok := ed.ParseMoveLines(noSubnum, subnum, mv, letter)
			if ok {
				end, ok := ed.RunMove(cmd)
				if !ok {
					return Cmd{}, false
				}
				return Cmd{
					Kind:     CmdOpCopyLineRegion,
					StartRow: start.Row,
					EndRow:   end.Row,
				}, true
			}
			cmd, ok = ed.ParseMoveRunes(noSubnum, subnum, mv, letter)
			if ok {
				meta, ok := MoveMetas[cmd.Kind]
				if !ok {
					return Cmd{}, false
				}
				end, ok := ed.RunMove(cmd)
				if !ok {
					return Cmd{}, false
				}
				return Cmd{
					Kind:      CmdOpCopyRegion,
					Start:     start,
					End:       end,
					Inclusive: meta.Inclusive,
				}, true
			}
			return Cmd{}, false
		case "d":
			b := ed.Buf()
			start := b.Loc
			cmd, ok := ed.ParseMoveLines(noSubnum, subnum, mv, letter)
			if ok {
				end, ok := ed.RunMove(cmd)
				if !ok {
					return Cmd{}, false
				}
				return Cmd{
					Kind:     CmdOpDeleteLineRegion,
					StartRow: start.Row,
					EndRow:   end.Row,
				}, true
			}
			cmd, ok = ed.ParseMoveRunes(noSubnum, subnum, mv, letter)
			if ok {
				meta, ok := MoveMetas[cmd.Kind]
				if !ok {
					return Cmd{}, false
				}
				end, ok := ed.RunMove(cmd)
				if !ok {
					return Cmd{}, false
				}
				return Cmd{
					Kind:      CmdOpDeleteRegion,
					Start:     start,
					End:       end,
					Inclusive: meta.Inclusive,
				}, true
			}
			return Cmd{}, false
		case "c":
			b := ed.Buf()
			start := b.Loc
			cmd, ok := ed.ParseMoveLines(noSubnum, subnum, mv, letter)
			if ok {
				end, ok := ed.RunMove(cmd)
				if !ok {
					return Cmd{}, false
				}
				return Cmd{
					Kind:     CmdOpChangeLineRegion,
					StartRow: start.Row,
					EndRow:   end.Row,
				}, true
			}
			cmd, ok = ed.ParseMoveRunes(noSubnum, subnum, mv, letter)
			if ok {
				meta, ok := MoveMetas[cmd.Kind]
				if !ok {
					return Cmd{}, false
				}
				end, ok := ed.RunMove(cmd)
				if !ok {
					return Cmd{}, false
				}
				return Cmd{
					Kind:      CmdOpChangeRegion,
					Start:     start,
					End:       end,
					Inclusive: meta.Inclusive,
				}, true
			}
			return Cmd{}, false
		}
	}

	switch op {
	case "yy", "Y":
		if reg == 0 {
			return Cmd{
				Kind: CmdOpCopyLine,
				Num:  num,
			}, true
		} else {
			return Cmd{
				Kind: CmdOpCopyLineIntoReg,
				Num:  num,
				Reg:  reg,
			}, true
		}
	case "yw":
		return Cmd{
			Kind: CmdOpCopyWord,
			Num:  num,
		}, true
	case "y$":
		return Cmd{
			Kind: CmdOpCopyToEnd,
			Num:  num,
		}, true

	case "p":
		if reg == 0 {
			return Cmd{
				Kind: CmdOpPaste,
				Num:  num,
			}, true
		} else {
			return Cmd{
				Kind: CmdOpPasteFromReg,
				Num:  num,
				Reg:  reg,
			}, true
		}
	case "P":
		return Cmd{
			Kind: CmdOpPasteBefore,
			Num:  num,
		}, true

	case "x":
		return Cmd{
			Kind: CmdOpDelete,
			Num:  num,
		}, true
	case "X":
		return Cmd{
			Kind: CmdOpDeleteBefore,
			Num:  num,
		}, true
	case "dd":
		return Cmd{
			Kind: CmdOpDeleteLine,
			Num:  num,
		}, true
	case "dw":
		return Cmd{
			Kind: CmdOpDeleteWord,
			Num:  num,
		}, true
	case "d$", "D":
		return Cmd{
			Kind: CmdOpDeleteToEnd,
			Num:  num,
		}, true

	case "cc":
		return Cmd{
			Kind: CmdOpChangeLine,
			Num:  num,
		}, true
	case "cw":
		return Cmd{
			Kind: CmdOpChangeWord,
			Num:  num,
		}, true
	case "C":
		return Cmd{
			Kind: CmdOpChangeToEnd,
			Num:  num,
		}, true
	case "s":
		return Cmd{
			Kind: CmdOpSubst,
			Num:  num,
		}, true
	case "S":
		return Cmd{
			Kind: CmdOpSubstLine,
			Num:  num,
		}, true
	}

	return Cmd{}, false
}

func (ed *Editor) ParseEdit(num int, op string, noSubnum bool, subnum int, mv string) (Cmd, bool) {
	switch op {
	case "J":
		return Cmd{
			Kind: CmdEditJoin,
			Num:  num,
		}, true
	case ">>":
		return Cmd{
			Kind: CmdEditIndent,
			Num:  num,
		}, true
	case "<<":
		return Cmd{
			Kind: CmdEditOutdent,
			Num:  num,
		}, true
	case ">":
		/*
			return Cmd{
				Kind: CmdEditIndentRegion,
				Start: Loc{},
				End: Loc{},
			}, true
		*/
	case "<":
		/*
			return Cmd{
				Kind: CmdEditOutdentRegion,
				Start: Loc{},
				End: Loc{},
			}, true
		*/
	}

	return Cmd{}, false
}

var compoundSet = map[string]struct{}{
	"]]": {},
	"[[": {},

	"``": {},
	"''": {},

	"z\r": {},
	"z.":  {},
	"z-":  {},

	"yy": {},
	"dd": {},
	"cc": {},
	">>": {},
	"<<": {},
	"ZZ": {},

	"dw": {},
	"cw": {},
}

var compoundHeadSet = map[rune]struct{}{
	']': {},
	'[': {},

	'`':  {},
	'\'': {},

	'z': {},

	'y': {},
	'd': {},
	'c': {},
	'>': {},
	'<': {},
	'Z': {},
}

func (ed *Editor) Parse() (Cmd, bool) {
	p := ed.parser

	if len(p.buf) < 1 {
		return Cmd{}, false
	}

	if p.buf[0] == '0' { // special
		return Cmd{Kind: CmdMoveToStart}, true
	}

	i := 0
	var reg rune = 0
	if p.buf[0] == '"' {
		if len(p.buf) > 1 {
			reg = p.buf[1]
			i += 2
		}
	}

	iPrev := i
	for i < len(p.buf) {
		if p.buf[i] < '0' || p.buf[i] > '9' {
			break
		}
		i++
	}
	noNum := i <= iPrev
	num := 1
	if i > iPrev {
		s := string(p.buf[iPrev:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		num = n
	}

	if i < len(p.buf) {
		var letter rune = 0
		_, ok := letterOpSet[p.buf[i]]
		if ok {
			if i+1 >= len(p.buf) {
				return Cmd{}, false
			}
			op := string(p.buf[i : i+1])
			letter = p.buf[i+1]
			cmd, ok := ed.ParseLetter(num, op, letter)
			if ok {
				return cmd, true
			}
		}
		_, ok = letterMoveSet[p.buf[i]]
		if ok {
			if i+1 >= len(p.buf) {
				return Cmd{}, false
			}
			mv := string(p.buf[i : i+1])
			letter = p.buf[i+1]
			cmd, ok := ed.ParseMoveLetterLines(num, mv, letter)
			if ok {
				return cmd, true
			}
			cmd, ok = ed.ParseMoveLetterRunes(num, mv, letter)
			if ok {
				return cmd, true
			}
		}
		if letter != 0 {
			return Cmd{Kind: CmdInvalid}, true
		}
	}

	iPrev = i
	for i < len(p.buf) {
		if i+1-iPrev == 2 {
			_, ok := compoundSet[string(p.buf[iPrev:i+1])]
			if !ok {
				break
			}
		}
		if p.buf[i] >= '0' && p.buf[i] <= '9' {
			break
		}
		i++
	}
	if i <= iPrev {
		return Cmd{}, false
	}

	mv := string(p.buf[iPrev:i])

	cmd, ok := ed.ParseMoveLines(noNum, num, mv, 0)
	if ok {
		return cmd, true
	}
	cmd, ok = ed.ParseMoveRunes(noNum, num, mv, 0)
	if ok {
		return cmd, true
	}
	if mv == "/" || mv == "?" {
		// XXX input pat
	}
	op := mv
	opFirst := p.buf[iPrev]

	cmd, ok = ed.ParseView(num, op)
	if ok {
		return cmd, true
	}
	cmd, ok = ed.ParseInsert(num, op)
	if ok {
		return cmd, true
	}
	cmd, ok = ed.ParseMisc(num, op)
	if ok {
		return cmd, true
	}

	iPrev = i
	for i < len(p.buf) {
		if p.buf[i] < '0' || p.buf[i] > '9' {
			break
		}
		i++
	}
	noSubnum := i <= iPrev
	subnum := 1
	if i > iPrev {
		s := string(p.buf[iPrev:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		subnum = n
	}

	mv = ""
	var letter rune = 0
	if i < len(p.buf) {
		_, ok := letterMoveSet[p.buf[i]]
		if ok {
			if i+1 < len(p.buf) {
				mv = string(p.buf[i : i+1])
				letter = p.buf[i+1]
			}
		}
	}

	if mv == "" {
		if i < len(p.buf) {
			mv = string(p.buf[i:])
		}
	}

	cmd, ok = ed.ParseOp(reg, num, op, noSubnum, subnum, mv, letter)
	if ok {
		return cmd, true
	}
	cmd, ok = ed.ParseEdit(num, op, noSubnum, subnum, mv)
	if ok {
		return cmd, true
	}

	if len(op) < 2 {
		_, ok := compoundHeadSet[opFirst]
		if ok {
			return Cmd{}, false
		}
	}
	return Cmd{Kind: CmdInvalid}, true
}
