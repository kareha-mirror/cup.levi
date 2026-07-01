package editor

import (
	"fmt"
	"strings"
)

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

func (args *Args) Code() string {
	code := strings.Builder{}
	found := false

	if args.Reg != 0 {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("Rg(%c)", args.Reg))
		found = true
	}

	if !args.NoNum && args.Num > 0 {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("%d", args.Num))
		found = true
	}

	if args.Op != 0 {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("Op(%c)", args.Op))
		found = true
	}

	if !args.NoSubnum && args.Subnum > 0 {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("%d", args.Subnum))
		found = true
	}

	if args.Mv != 0 {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("Mv(%c)", args.Mv))
		found = true
	}

	if args.Rune != 0 {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("Ch(%c)", args.Rune))
		found = true
	}

	return code.String()
}
