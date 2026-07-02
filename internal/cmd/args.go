package cmd

import (
	"fmt"
	"strings"
)

// Parsed arguments of command.
type Args struct {
	Reg    rune
	Has    bool
	Num    int
	Op     rune
	HasSub bool
	SubNum int
	Mv     rune
	Rune   rune
}

// Uses sub number as main number.
func (a Args) sub() Args {
	sub := a
	sub.Has = a.HasSub
	sub.Num = a.SubNum
	sub.HasSub = false
	sub.SubNum = 0
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

	if a.Has && a.Num > 0 {
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

	if a.HasSub && a.SubNum > 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("%d", a.SubNum))
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
