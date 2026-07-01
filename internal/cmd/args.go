package cmd

import (
	"fmt"
	"strings"
)

// Parsed arguments of command.
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
func (a Args) sub() Args {
	sub := a
	sub.NoNum = a.NoSubnum
	sub.Num = a.Subnum
	sub.NoSubnum = true
	sub.Subnum = 1
	return sub
}

// Compiles into command.
func (a Args) Compile() (Pair, bool) {
	op, ok := a.compileOp()
	if ok {
		return op, true
	}
	mv, ok := a.compileMove(false)
	if ok {
		return Pair{Mv: mv}, true
	}
	return Pair{}, false
}

// Returns mnemonic code for parsed arguments of command.
func (a Args) Code() string {
	b := strings.Builder{}
	first := true

	if a.Reg != 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("Rg(%c)", a.Reg))
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
		b.WriteString(fmt.Sprintf("Ch(%c)", a.Rune))
		first = false
	}

	return b.String()
}
