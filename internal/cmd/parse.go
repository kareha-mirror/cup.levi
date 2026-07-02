package cmd

import (
	"strconv"
)

var isRuneOp = map[rune]struct{}{
	'm': {},
	'r': {},
}

var isRuneMove = map[rune]struct{}{
	'`':  {},
	'\'': {},
	'f':  {},
	'F':  {},
	't':  {},
	'T':  {},
}

var isMove = map[rune]struct{}{
	'n': {},
	'N': {},

	'h':  {},
	'j':  {},
	'k':  {},
	'l':  {},
	'0':  {},
	'$':  {},
	'^':  {},
	'|':  {},
	'w':  {},
	'b':  {},
	'e':  {},
	'W':  {},
	'B':  {},
	'E':  {},
	'\r': {},
	'+':  {},
	'-':  {},
	'G':  {},
	')':  {},
	'(':  {},
	'}':  {},
	'{':  {},
	'H':  {},
	'M':  {},
	'L':  {},
	';':  {},
	',':  {},

	'g': {}, // XXX debug
}

func Parse(buf []rune) Args {
	a := Args{}
	if len(buf) < 1 {
		return a
	}
	i := 0

	if buf[i] == '0' { // special
		a.Mv = buf[i]
		return a
	}

	// register
	if buf[i] == '"' {
		i++
		if i < len(buf) {
			a.Reg = buf[i]
			i++
		}
	}

	// main number
	iPrev := i
	for i < len(buf) {
		if buf[i] < '0' || buf[i] > '9' {
			break
		}
		i++
	}
	a.Has = i > iPrev
	a.Num = 1
	if i > iPrev {
		s := string(buf[iPrev:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		a.Num = n
	}

	// register again
	if a.Reg == 0 && i < len(buf) {
		if buf[i] == '"' {
			i++
			if i < len(buf) {
				a.Reg = buf[i]
				i++
			}
		}
	}

	// rune command
	if i < len(buf) {
		_, ok := isRuneOp[buf[i]]
		if ok {
			a.Op = buf[i]
			if i+1 >= len(buf) {
				return a
			}
			a.Rune = buf[i+1]
			return a
		}
		_, ok = isRuneMove[buf[i]]
		if ok {
			a.Mv = buf[i]
			if i+1 >= len(buf) {
				return a
			}
			a.Rune = buf[i+1]
			return a
		}
		if a.Rune != 0 {
			a.Mv = 0
			return a
		}
	}

	// detect non number
	iPrev = i
	for i < len(buf) {
		if i+1-iPrev == 2 {
			break
		}
		if buf[i] >= '0' && buf[i] <= '9' {
			break
		}
		i++
	}
	if i <= iPrev {
		return a
	}

	// detect motion command
	a.Mv = buf[iPrev]
	_, ok := isMove[a.Mv]
	if ok {
		return a
	}

	// cast motion to operation
	a.Op = a.Mv
	a.Mv = 0

	// sub number
	iPrev = i
	for i < len(buf) {
		if buf[i] < '0' || buf[i] > '9' {
			break
		}
		i++
	}
	a.HasSub = i > iPrev
	a.SubNum = 1
	if i > iPrev {
		s := string(buf[iPrev:i])
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		a.SubNum = n
	}

	// rune motion command
	if i < len(buf) {
		_, ok := isRuneMove[buf[i]]
		if ok {
			if i+1 < len(buf) {
				a.Mv = buf[i]
				a.Rune = buf[i+1]
				return a
			}
		}
	}

	// last motion command
	if a.Mv == 0 {
		if i < len(buf) {
			a.Mv = buf[i]
		}
	}
	return a
}
