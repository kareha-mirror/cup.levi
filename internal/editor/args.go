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

// Returns mnemonic code for parsed arguments of vi command.
func (args *Args) Code() string {
	b := strings.Builder{}
	first := true

	if args.Reg != 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("Reg(%c)", args.Reg))
		first = false
	}

	if !args.NoNum && args.Num > 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("%d", args.Num))
		first = false
	}

	if args.Op != 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("Op(%c)", args.Op))
		first = false
	}

	if !args.NoSubnum && args.Subnum > 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("%d", args.Subnum))
		first = false
	}

	if args.Mv != 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("Mv(%c)", args.Mv))
		first = false
	}

	if args.Rune != 0 {
		if !first {
			b.WriteRune('-')
		}
		b.WriteString(fmt.Sprintf("Rune(%c)", args.Rune))
		first = false
	}

	return b.String()
}
