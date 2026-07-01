package editor

import (
	"strconv"
)

var isRuneOp = map[rune]struct{}{
	'm': {},
	'r': {},
}

var isRuneMove = map[rune]struct{}{
	'\'': {},
	'`':  {},
	'f':  {},
	'F':  {},
	't':  {},
	'T':  {},
}

var isCompound = map[rune]struct{}{
	'y': {},
	'd': {},
	'c': {},
	'>': {},
	'<': {},

	']': {},
	'[': {},

	'`':  {},
	'\'': {},

	'z': {},

	'Z': {},
}

type Parser struct {
	buf   []rune
	Cache string
	Args  Args
	Ok    bool
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
	p.Args = Args{}
	p.Ok = false
}

func (p *Parser) Parse() (CmdPair, bool) {
	p.Args = Args{}
	p.Ok = false
	a := &p.Args

	if len(p.buf) < 1 {
		return CmdPair{}, false
	}

	if p.buf[0] == '0' { // special
		a.Mv = p.buf[0]
		p.Ok = true
		return CmdPair{
			Mv: Cmd{Kind: MoveToStart},
		}, true
	}

	i := 0
	if p.buf[i] == '"' {
		i++
		if i < len(p.buf) {
			a.Reg = p.buf[i]
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
	a.NoNum = i <= iPrev
	a.Num = 1
	if i > iPrev {
		s := string(p.buf[iPrev:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		a.Num = n
	}

	if a.Reg == 0 && i < len(p.buf) {
		if p.buf[i] == '"' {
			i++
			if i < len(p.buf) {
				a.Reg = p.buf[i]
				i++
			}
		}
	}

	if i < len(p.buf) {
		_, ok := isRuneOp[p.buf[i]]
		if ok {
			a.Op = p.buf[i]
			if i+1 >= len(p.buf) {
				return CmdPair{}, false
			}
			a.Rune = p.buf[i+1]
			cmd, ok := a.ParseRune()
			if ok {
				p.Ok = true
				return CmdPair{Op: cmd}, true
			}
		}
		_, ok = isRuneMove[p.buf[i]]
		if ok {
			a.Mv = p.buf[i]
			if i+1 >= len(p.buf) {
				return CmdPair{}, false
			}
			a.Rune = p.buf[i+1]
			cmd, ok := a.ParseMoveRune()
			if ok {
				p.Ok = true
				return CmdPair{Mv: cmd}, true
			}
		}
		if a.Rune != 0 {
			p.Ok = true
			return CmdPair{
				Op: Cmd{Kind: InvalidCmd},
			}, true
		}
	}

	iPrev = i
	for i < len(p.buf) {
		if i+1-iPrev == 2 {
			break
		}
		if p.buf[i] >= '0' && p.buf[i] <= '9' {
			break
		}
		i++
	}
	if i <= iPrev {
		return CmdPair{}, false
	}

	a.Mv = p.buf[iPrev]
	cmd, ok := a.ParseMove(false)
	if ok {
		p.Ok = true
		return CmdPair{Mv: cmd}, true
	}
	if a.Mv == '/' || a.Mv == '?' {
		// XXX input pat
	}
	a.Op = a.Mv
	a.Mv = 0
	opFirst := p.buf[iPrev]

	cmd, ok = a.ParseView()
	if ok {
		p.Ok = true
		return CmdPair{Op: cmd}, true
	}
	cmd, ok = a.ParseInsert()
	if ok {
		p.Ok = true
		return CmdPair{Op: cmd}, true
	}
	cmd, ok = a.ParseMisc()
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
	a.NoSubnum = i <= iPrev
	a.Subnum = 1
	if i > iPrev {
		s := string(p.buf[iPrev:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		a.Subnum = n
	}

	if i < len(p.buf) {
		_, ok := isRuneMove[p.buf[i]]
		if ok {
			if i+1 < len(p.buf) {
				a.Mv = p.buf[i]
				a.Rune = p.buf[i+1]
			}
		}
	}

	if a.Mv == 0 {
		if i < len(p.buf) {
			a.Mv = p.buf[i]
		}
	}

	c, ok := a.ParseOp()
	if ok {
		p.Ok = true
		return c, true
	}
	c, ok = a.ParseEdit()
	if ok {
		p.Ok = true
		return c, true
	}
	c, ok = a.ParseCompound()
	if ok {
		p.Ok = true
		return c, true
	}

	_, ok = isCompound[opFirst]
	if ok {
		return CmdPair{}, false
	}
	p.Ok = true
	return CmdPair{
		Op: Cmd{Kind: InvalidCmd},
	}, true
}
