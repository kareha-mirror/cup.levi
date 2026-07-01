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

func (ed *Editor) Parse() (CmdPair, bool) {
	p := &ed.parser
	p.Args = Args{}
	p.Ok = false
	args := &p.Args

	if len(p.buf) < 1 {
		return CmdPair{}, false
	}

	if p.buf[0] == '0' { // special
		args.Mv = p.buf[0]
		p.Ok = true
		return CmdPair{
			Mv: Cmd{Kind: MoveToStart},
		}, true
	}

	i := 0
	if p.buf[i] == '"' {
		i++
		if i < len(p.buf) {
			args.Reg = p.buf[i]
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
	args.NoNum = i <= iPrev
	args.Num = 1
	if i > iPrev {
		s := string(p.buf[iPrev:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		args.Num = n
	}

	if args.Reg == 0 && i < len(p.buf) {
		if p.buf[i] == '"' {
			i++
			if i < len(p.buf) {
				args.Reg = p.buf[i]
				i++
			}
		}
	}

	if i < len(p.buf) {
		_, ok := isRuneOp[p.buf[i]]
		if ok {
			args.Op = p.buf[i]
			if i+1 >= len(p.buf) {
				return CmdPair{}, false
			}
			args.Rune = p.buf[i+1]
			cmd, ok := ed.ParseRune(args.Num, args.Op, args.Rune)
			if ok {
				p.Ok = true
				return CmdPair{Op: cmd}, true
			}
		}
		_, ok = isRuneMove[p.buf[i]]
		if ok {
			args.Mv = p.buf[i]
			if i+1 >= len(p.buf) {
				return CmdPair{}, false
			}
			args.Rune = p.buf[i+1]
			cmd, ok := ed.ParseMoveRune(args.Num, args.Mv, args.Rune)
			if ok {
				p.Ok = true
				return CmdPair{Mv: cmd}, true
			}
		}
		if args.Rune != 0 {
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

	args.Mv = p.buf[iPrev]
	cmd, ok := ed.ParseMove(args.NoNum, args.Num, args.Mv, 0, false)
	if ok {
		p.Ok = true
		return CmdPair{Mv: cmd}, true
	}
	if args.Mv == '/' || args.Mv == '?' {
		// XXX input pat
	}
	args.Op = args.Mv
	args.Mv = 0
	opFirst := p.buf[iPrev]

	cmd, ok = ed.ParseView(args.Num, args.Op)
	if ok {
		p.Ok = true
		return CmdPair{Op: cmd}, true
	}
	cmd, ok = ed.ParseInsert(args.Num, args.Op)
	if ok {
		p.Ok = true
		return CmdPair{Op: cmd}, true
	}
	cmd, ok = ed.ParseMisc(args.NoNum, args.Num, args.Op)
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
	args.NoSubnum = i <= iPrev
	args.Subnum = 1
	if i > iPrev {
		s := string(p.buf[iPrev:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		args.Subnum = n
	}

	if i < len(p.buf) {
		_, ok := isRuneMove[p.buf[i]]
		if ok {
			if i+1 < len(p.buf) {
				args.Mv = p.buf[i]
				args.Rune = p.buf[i+1]
			}
		}
	}

	if args.Mv == 0 {
		if i < len(p.buf) {
			args.Mv = p.buf[i]
		}
	}

	c, ok := ed.ParseOp(
		args.Reg, args.Num, args.Op,
		args.NoSubnum, args.Subnum, args.Mv, args.Rune,
	)
	if ok {
		p.Ok = true
		return c, true
	}
	c, ok = ed.ParseEdit(
		args.Num, args.Op, args.NoSubnum, args.Subnum, args.Mv, args.Rune,
	)
	if ok {
		p.Ok = true
		return c, true
	}
	c, ok = ed.ParseCompound(
		args.Num, args.Op, args.NoSubnum, args.Subnum, args.Mv, args.Rune,
	)
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
