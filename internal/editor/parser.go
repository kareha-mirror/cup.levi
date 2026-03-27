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

func (p *Parser) ParseMove(noNum bool, num int, op string) (Cmd, bool) {
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

func (p *Parser) ParseOp(num int, op string) (Cmd, bool) {
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

func (p *Parser) ParseMisc(op string) (Cmd, bool) {
	switch op {
	case "ZZ":
		return Cmd{Kind: CmdMiscSaveAndQuit}, true

		// TODO
	}

	return Cmd{}, false
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
	for i < len(p.buf) {
		if p.buf[i] < '0' || p.buf[i] > '9' {
			break
		}
		i++
	}
	noNum := i == 0
	num := 1
	if i > 0 {
		s := string(p.buf[:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		num = n
	}

	iPrev := i
	for i < len(p.buf) {
		if p.buf[i] >= '0' && p.buf[i] <= '9' {
			break
		}
		i++
	}
	if i <= iPrev {
		return Cmd{}, false
	}
	op := string(p.buf[iPrev:i])
	cmd, ok := p.ParseMove(noNum, num, op)
	if ok {
		return cmd, true
	}
	cmd, ok = p.ParseInsert(num, op)
	if ok {
		return cmd, true
	}
	cmd, ok = p.ParseOp(num, op)
	if ok {
		return cmd, true
	}
	cmd, ok = p.ParseMisc(op)
	if ok {
		return cmd, true
	}

	// TODO

	return Cmd{}, false
}
