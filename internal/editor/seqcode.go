package editor

import (
	"fmt"
	"strings"

	"tea.kareha.org/cup/levi/internal/rkind"
)

func (ed *Editor) SeqCode() string {
	t := ed.parser.Tokens
	code := strings.Builder{}
	found := false

	if t.Reg != "" {
		if found {
			code.WriteRune('-')
		}
		escaped := rkind.Escape(t.Reg)
		code.WriteString(fmt.Sprintf("Rg(%s)", escaped))
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
		escaped := rkind.Escape(t.Op)
		code.WriteString(fmt.Sprintf("Op(%s)", escaped))
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
		escaped := rkind.Escape(t.Mv)
		code.WriteString(fmt.Sprintf("Mv(%s)", escaped))
		found = true
	}

	if t.Letter != 0 {
		if found {
			code.WriteRune('-')
		}
		escaped := rkind.Escape(string(t.Letter))
		code.WriteString(fmt.Sprintf("Lt(%s)", escaped))
		found = true
	}

	if ed.parser.Ok {
		code.WriteRune('.')
	}

	return code.String()
}
