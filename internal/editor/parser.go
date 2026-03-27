package editor

import (
	"strconv"
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
	return string(p.buf)
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
	'm':  {},
	'\'': {},
	'`':  {},
	'r':  {},
}

var letterMoveSet = map[rune]struct{}{
	'f': {},
	'F': {},
	't': {},
	'T': {},
	';': {},
	',': {},
}

func (p *Parser) ParseMove(noNum bool, num int, op string, letter rune) (Cmd, bool) {
	switch op {
	case "h":
		return Cmd{
			Kind: CmdMoveLeft,
			Num:  num,
		}, true
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

		// TODO
	}

	return Cmd{}, false
}

func (p *Parser) ParseLetter(num int, op string, letter rune) (Cmd, bool) {
	switch op {
	case "m":
		return Cmd{}, false // TODO
	case "'":
		return Cmd{}, false // TODO
	case "`":
		return Cmd{}, false // TODO
	case "r":
		return Cmd{}, false // TODO
	}

	return Cmd{}, false
}

func (p *Parser) ParseInsert(num int, op string) (Cmd, bool) {
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

		// TODO
	}

	return Cmd{}, false
}

func (p *Parser) ParseMisc(op string) (Cmd, bool) {
	switch op {
	case "ZZ":
		return Cmd{Kind: CmdMiscSaveAndQuit}, true

		// TODO
	}

	return Cmd{}, false
}

func (p *Parser) ParseOp(reg rune, num int, op string, noSubnum bool, subnum int, mv string) (Cmd, bool) {
	switch op {
	case "x":
		return Cmd{
			Kind: CmdOpDelete,
			Num:  num,
		}, true

		// TODO
	}

	return Cmd{}, false
}

var compoundSet = map[rune]struct{}{
	'y': {},
	'd': {},
	'c': {},
	'>': {},
	'<': {},
	'Z': {},
}

func (p *Parser) Parse() (Cmd, bool) {
	if len(p.buf) < 1 {
		return Cmd{}, false
	}

	if len(p.buf) == 1 {
		if p.buf[0] == '0' { // special
			return Cmd{Kind: CmdMoveToStart}, true
		}
	}

	i := 0
	var reg rune = 0
	if p.buf[i] == '"' {
		if len(p.buf) > i+1 {
			reg = p.buf[i+1]
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
	noNum := i == iPrev
	num := 1
	if i > 0 {
		s := string(p.buf[iPrev:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		num = n
	}

	iPrev = i
	for i < len(p.buf) {
		if p.buf[i] >= '0' && p.buf[i] <= '9' {
			break
		}
		i++
	}
	if i <= iPrev {
		return Cmd{}, false
	}

	if i+1 < len(p.buf) {
		var letter rune = 0
		_, ok := letterOpSet[p.buf[i]]
		if ok {
			op := string(p.buf[i : i+1])
			letter = p.buf[i+1]
			cmd, ok := p.ParseLetter(num, op, letter)
			if ok {
				return cmd, true
			}
		}
		_, ok = letterMoveSet[p.buf[i]]
		if ok {
			mv := string(p.buf[i : i+1])
			letter = p.buf[i+1]
			cmd, ok := p.ParseMove(noNum, num, mv, letter)
			if ok {
				return cmd, true
			}
		}
		if letter != 0 {
			return Cmd{Kind: CmdInvalid}, true
		}
	}

	mv := string(p.buf[iPrev:i])

	cmd, ok := p.ParseMove(noNum, num, mv, 0)
	if ok {
		return cmd, true
	}
	op := mv
	opFirst := p.buf[iPrev]

	cmd, ok = p.ParseInsert(num, op)
	if ok {
		return cmd, true
	}
	cmd, ok = p.ParseMisc(op)
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
	noSubnum := i == iPrev
	subnum := 1
	if i > iPrev {
		s := string(p.buf[iPrev:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		subnum = n
	}

	iPrev = i
	for i < len(p.buf) {
		if p.buf[i] >= '0' && p.buf[i] <= '9' {
			break
		}
		i++
	}

	var letter rune = 0
	if i+1 < len(p.buf) {
		_, ok := letterMoveSet[p.buf[i]]
		if ok {
			mv = string(p.buf[i : i+1])
			letter = p.buf[i+1]
		}
	}

	if letter == 0 {
		mv = ""
		if i > iPrev {
			mv = string(p.buf[iPrev:i])
		}
	}

	cmd, ok = p.ParseOp(reg, num, op, noSubnum, subnum, mv)
	if ok {
		return cmd, true
	}

	if len(op) < 2 {
		_, ok := compoundSet[opFirst]
		if ok {
			return Cmd{}, false
		}
	}
	return Cmd{Kind: CmdInvalid}, true
}
