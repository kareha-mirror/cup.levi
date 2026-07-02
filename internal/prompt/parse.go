package prompt

import (
	"fmt"
	"strconv"
	"strings"
)

func Parse(prompt string) (Cmd, bool) {
	if prompt == "" {
		return Cmd{Kind: Invalid}, false
	}

	if prompt[0] == '+' {
		if len(prompt) < 2 {
			return Cmd{Kind: MoveByLine, Num: 1}, true
		} else {
			n, err := strconv.ParseInt(prompt, 10, 32)
			if err == nil {
				return Cmd{Kind: MoveByLine, Num: int(n)}, true
			}
		}
	}

	if prompt[0] == '-' {
		if len(prompt) < 2 {
			return Cmd{Kind: MoveBackwardByLine, Num: 1}, true
		} else {
			n, err := strconv.ParseInt(prompt, 10, 32)
			if err == nil {
				return Cmd{Kind: MoveBackwardByLine, Num: int(n)}, true
			}
		}
	}

	n, err := strconv.ParseInt(prompt, 10, 32)
	if err == nil {
		return Cmd{Kind: MoveToLine, Num: int(n)}, true
	}

	parts := strings.Split(prompt, " ")

	switch parts[0] {
	case "wq":
		return Cmd{Kind: SaveAndQuit}, true
	case "w":
		if len(parts) > 1 {
			return Cmd{Kind: Save, Name: parts[1]}, true
		} else {
			return Cmd{Kind: Save}, true
		}
	case "w!":
		if len(parts) > 1 {
			return Cmd{Kind: ForceSave, Name: parts[1]}, true
		} else {
			return Cmd{Kind: ForceSave}, true
		}
	case "q":
		return Cmd{Kind: Quit}, true
	case "q!":
		return Cmd{Kind: ForceQuit}, true
	case "e":
		if len(parts) > 1 {
			return Cmd{Kind: Load, Name: parts[1]}, true
		} else {
			return Cmd{Kind: Load}, true
		}
	case "e!":
		if len(parts) > 1 {
			return Cmd{Kind: ForceLoad, Name: parts[1]}, true
		} else {
			return Cmd{Kind: ForceLoad}, true
		}
	case "r":
		return Cmd{Kind: Read}, true
	case "n", "next":
		return Cmd{Kind: Next}, true
	case "prev", "previous":
		return Cmd{Kind: Prev}, true

	case "sh", "shell":
		return Cmd{Kind: Shell}, true

	case "wa":
		return Cmd{Kind: SaveAll}, true
	case "qa":
		return Cmd{Kind: QuitAll}, true
	case "qa!":
		return Cmd{Kind: ForceQuitAll}, true

	case "set":
		if len(parts) < 2 {
			// TODO show variables
			return Cmd{Kind: Invalid}, false
		}
		if strings.HasPrefix(parts[1], "ts=") {
			ns := parts[1][3:]
			n, err := strconv.ParseUint(ns, 10, 16)
			if err != nil {
				name := fmt.Sprintf(
					"set: %s option: %s is an illegal number.", ns, ns,
				)
				return Cmd{Kind: Ring, Name: name}, false
			}
			return Cmd{Kind: TabStop, Num: int(n)}, true
		}
		switch parts[1] {
		case "ai", "autoindent":
			return Cmd{Kind: AutoIndent}, true
		case "noai", "noautoindent":
			return Cmd{Kind: NoAutoIndent}, true
		}
		// TODO set all
		name := fmt.Sprintf(
			"set no %s option: 'set all' gives all option values.",
			parts[1],
		)
		return Cmd{Kind: Ring, Name: name}, false

	case "open":
		if len(parts) > 1 {
			return Cmd{Kind: Open, Name: parts[1]}, true
		} else {
			return Cmd{Kind: Open}, true
		}
	case "nl", "newline":
		if len(parts) > 1 {
			return Cmd{Kind: Newline, Name: parts[1]}, true
		} else {
			return Cmd{Kind: Newline}, true
		}
	case "col", "colors", "colorscheme":
		if len(parts) > 1 {
			return Cmd{Kind: Colors, Name: parts[1]}, true
		} else {
			return Cmd{Kind: Colors}, true
		}

	case "mem":
		return Cmd{Kind: Mem}, true
	case "hello":
		if len(parts) > 1 {
			n, err := strconv.ParseUint(parts[1], 10, 16)
			if err != nil {
				name := "Goodbye, World!"
				return Cmd{Kind: Error, Name: name}, false
			}
			return Cmd{Kind: Hello, Num: int(n)}, true
		} else {
			return Cmd{Kind: Hello}, true
		}

	default:
		name := fmt.Sprintf("The %s command is unknown.", parts[0])
		return Cmd{Kind: Error, Name: name}, false
	}
}
