package editor

import (
	"fmt"
	"strings"
)

// Parsed arguments of vi command.
type Args struct {
	Reg      rune
	NoNum    bool
	Num      int
	Op       rune
	NoSubnum bool
	Subnum   int
	Mv       rune
	Rune     rune
}

// Uses sub number as main number.
func (a *Args) Sub() *Args {
	sub := *a
	sub.NoNum = a.NoSubnum
	sub.Num = a.Subnum
	sub.NoSubnum = true
	sub.Subnum = 1
	return &sub
}

// Parses into command.
func (a *Args) Parse() (CmdPair, bool) {
	cp, ok := a.ParseOp()
	if ok {
		return cp, true
	}

	cp, ok = a.ParseEdit()
	if ok {
		return cp, true
	}

	cp, ok = a.ParseCompound()
	if ok {
		return cp, true
	}

	op, ok := a.ParseInsert()
	if ok {
		return CmdPair{
			Op: op,
		}, true
	}

	op, ok = a.ParseRune()
	if ok {
		return CmdPair{
			Op: op,
		}, true
	}

	op, ok = a.ParseView()
	if ok {
		return CmdPair{
			Op: op,
		}, true
	}

	op, ok = a.ParseMisc()
	if ok {
		return CmdPair{
			Op: op,
		}, true
	}

	mv, ok := a.ParseMove(false)
	if ok {
		return CmdPair{
			Mv: mv,
		}, true
	}

	return CmdPair{}, false
}

// Returns mnemonic code for parsed arguments of vi command.
func (a *Args) Code() string {
	b := strings.Builder{}
	first := true

	if a.Reg != 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("Reg(%c)", a.Reg))
		first = false
	}

	if !a.NoNum && a.Num > 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("%d", a.Num))
		first = false
	}

	if a.Op != 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("Op(%c)", a.Op))
		first = false
	}

	if !a.NoSubnum && a.Subnum > 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("%d", a.Subnum))
		first = false
	}

	if a.Mv != 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("Mv(%c)", a.Mv))
		first = false
	}

	if a.Rune != 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("Rune(%c)", a.Rune))
		first = false
	}

	return b.String()
}
