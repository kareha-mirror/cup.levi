package editor

import (
	"fmt"
	"strings"
)

func (ed *Editor) SeqCode() string {
	t := ed.parser.Tokens
	code := strings.Builder{}
	found := false

	if t.Reg != "" {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("Rg(%s)", t.Reg))
		found = true
	}

	if !t.NoNum && t.Num > 0 {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("%d", t.Num))
		found = true
	}

	if t.Op != "" {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("Op(%s)", t.Op))
		found = true
	}

	if !t.NoSubnum && t.Subnum > 0 {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("%d", t.Subnum))
		found = true
	}

	if t.Mv != "" {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("Mv(%s)", t.Mv))
		found = true
	}

	if t.Letter != 0 {
		if found {
			code.WriteRune('-')
		}
		code.WriteString(fmt.Sprintf("Lt(%s)", t.Letter))
		found = true
	}

	if ed.parser.Ok {
		code.WriteRune('.')
	}

	return code.String()
}
