package editor

import (
	"fmt"
	"strings"
)

type Args struct {
	Reg      string
	NoNum    bool
	Num      int
	Op       string
	NoSubnum bool
	Subnum   int
	Mv       string
	Rune     rune
}

func (args *Args) Code() string {
	code := strings.Builder{}
	found := false

	if args.Reg != "" {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("Rg(%s)", args.Reg))
		found = true
	}

	if !args.NoNum && args.Num > 0 {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("%d", args.Num))
		found = true
	}

	if args.Op != "" {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("Op(%s)", args.Op))
		found = true
	}

	if !args.NoSubnum && args.Subnum > 0 {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("%d", args.Subnum))
		found = true
	}

	if args.Mv != "" {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("Mv(%s)", args.Mv))
		found = true
	}

	if args.Rune != 0 {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("Ch(%s)", args.Rune))
		found = true
	}

	return code.String()
}
