package editor

import (
	"strconv"
)

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

type Tokens struct {
	Reg      string
	NoNum    bool
	Num      int
	Op       string
	NoSubnum bool
	Subnum   int
	Mv       string
	Letter   rune
}

type Parser struct {
	buf    []rune
	Cache  string
	Tokens Tokens
	Ok     bool
}

func (p *Parser) String() string {
	return string(p.buf)
}

func (p *Parser) WriteRune(r rune) {
	p.buf = append(p.buf, r)
	p.Cache = p.String()
}

func (p *Parser) Backspace() bool {
	if len(p.buf) < 1 {
		return false
	}
	p.buf = p.buf[:len(p.buf)-1]
	p.Cache = p.String()
	return true
}

func (p *Parser) Reset() {
	p.buf = p.buf[:0]
}

func (p *Parser) ResetAll() {
	p.Reset()
	p.Cache = ""
	p.Tokens = Tokens{}
	p.Ok = false
}

func (ed *Editor) Parse() (CmdPair, bool) {
	p := &ed.parser
	p.Tokens = Tokens{}
	p.Ok = false
	t := &p.Tokens

	if len(p.buf) < 1 {
		return CmdPair{}, false
	}

	if p.buf[0] == '0' { // special
		t.Mv = string(p.buf[0])
		p.Ok = true
		return CmdPair{
			Mv: Cmd{Kind: MoveToStart},
		}, true
	}

	i := 0
	if p.buf[0] == '"' {
		i++
		if len(p.buf) > i {
			t.Reg = string(p.buf[i])
			i++
		}
	}

	iPrev := i
	for i < len(p.buf) {
		if p.buf[i] < '0' || p.buf[i] > '9' {
			break
		}
		i++
	}
	t.NoNum = i <= iPrev
	t.Num = 1
	if i > iPrev {
		s := string(p.buf[iPrev:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		t.Num = n
	}

	if i < len(p.buf) {
		_, ok := letterOpSet[p.buf[i]]
		if ok {
			t.Op = string(p.buf[i])
			if i+1 >= len(p.buf) {
				return CmdPair{}, false
			}
			t.Letter = p.buf[i+1]
			cmd, ok := ed.ParseLetter(t.Num, t.Op, t.Letter)
			if ok {
				p.Ok = true
				return CmdPair{Op: cmd}, true
			}
		}
		_, ok = letterMoveSet[p.buf[i]]
		if ok {
			t.Mv = string(p.buf[i])
			if i+1 >= len(p.buf) {
				return CmdPair{}, false
			}
			t.Letter = p.buf[i+1]
			cmd, ok := ed.ParseMoveLetter(t.Num, t.Mv, t.Letter)
			if ok {
				p.Ok = true
				return CmdPair{Mv: cmd}, true
			}
		}
		if t.Letter != 0 {
			p.Ok = true
			return CmdPair{
				Op: Cmd{Kind: InvalidCmd},
			}, true
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
		return CmdPair{}, false
	}

	t.Mv = string(p.buf[iPrev:i])

	cmd, ok := ed.ParseMove(t.NoNum, t.Num, t.Mv, 0)
	if ok {
		p.Ok = true
		return CmdPair{Mv: cmd}, true
	}
	if t.Mv == "/" || t.Mv == "?" {
		// XXX input pat
	}
	t.Op = t.Mv
	t.Mv = ""
	opFirst := p.buf[iPrev]

	cmd, ok = ed.ParseView(t.Num, t.Op)
	if ok {
		p.Ok = true
		return CmdPair{Op: cmd}, true
	}
	cmd, ok = ed.ParseInsert(t.Num, t.Op)
	if ok {
		p.Ok = true
		return CmdPair{Op: cmd}, true
	}
	cmd, ok = ed.ParseMisc(t.Num, t.Op)
	if ok {
		p.Ok = true
		return CmdPair{Op: cmd}, true
	}

	iPrev = i
	for i < len(p.buf) {
		if p.buf[i] < '0' || p.buf[i] > '9' {
			break
		}
		i++
	}
	t.NoSubnum = i <= iPrev
	t.Subnum = 1
	if i > iPrev {
		s := string(p.buf[iPrev:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		t.Subnum = n
	}

	if i < len(p.buf) {
		_, ok := letterMoveSet[p.buf[i]]
		if ok {
			if i+1 < len(p.buf) {
				t.Mv = string(p.buf[i])
				t.Letter = p.buf[i+1]
			}
		}
	}

	if t.Mv == "" {
		if i < len(p.buf) {
			t.Mv = string(p.buf[i:])
		}
	}

	c, ok := ed.ParseOp(
		t.Reg, t.Num, t.Op, t.NoSubnum, t.Subnum, t.Mv, t.Letter,
	)
	if ok {
		p.Ok = true
		return c, true
	}
	c, ok = ed.ParseEdit(t.Num, t.Op, t.NoSubnum, t.Subnum, t.Mv, t.Letter)
	if ok {
		p.Ok = true
		return c, true
	}

	if len(t.Op) < 2 {
		_, ok := compoundHeadSet[opFirst]
		if ok {
			return CmdPair{}, false
		}
	}
	p.Ok = true
	return CmdPair{
		Op: Cmd{Kind: InvalidCmd},
	}, true
}
