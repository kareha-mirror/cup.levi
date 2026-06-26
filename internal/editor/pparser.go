package editor

import (
	"strconv"
	"strings"
)

func (ed *Editor) ParsePrompt() (Pcmd, bool) {
	if ed.prompt.Len() < 1 {
		return Pcmd{Kind: PcmdInvalid}, false
	}
	prompt := ed.prompt.String()

	if prompt[0] == '+' {
		if len(prompt) < 2 {
			return Pcmd{Kind: PcmdMoveByLine, Num: 1}, true
		} else {
			n, err := strconv.ParseInt(prompt, 10, 32)
			if err == nil {
				return Pcmd{Kind: PcmdMoveByLine, Num: int(n)}, true
			}
		}
	}

	if prompt[0] == '-' {
		if len(prompt) < 2 {
			return Pcmd{Kind: PcmdMoveBackwardByLine, Num: 1}, true
		} else {
			n, err := strconv.ParseInt(prompt, 10, 32)
			if err == nil {
				return Pcmd{Kind: PcmdMoveBackwardByLine, Num: int(n)}, true
			}
		}
	}

	n, err := strconv.ParseInt(prompt, 10, 32)
	if err == nil {
		return Pcmd{Kind: PcmdMoveToLine, Num: int(n)}, true
	}

	parts := strings.Split(prompt, " ")

	switch parts[0] {
	case "wq":
		return Pcmd{Kind: PcmdSaveAndQuit}, true
	case "w":
		if len(parts) > 1 {
			return Pcmd{Kind: PcmdSave, Name: parts[1]}, true
		} else {
			return Pcmd{Kind: PcmdSave}, true
		}
	case "w!":
		if len(parts) > 1 {
			return Pcmd{Kind: PcmdForceSave, Name: parts[1]}, true
		} else {
			return Pcmd{Kind: PcmdForceSave}, true
		}
	case "q":
		return Pcmd{Kind: PcmdQuit}, true
	case "q!":
		return Pcmd{Kind: PcmdForceQuit}, true
	case "e":
		if len(parts) > 1 {
			return Pcmd{Kind: PcmdOpen, Name: parts[1]}, true
		} else {
			return Pcmd{Kind: PcmdOpen}, true
		}
	case "e!":
		if len(parts) > 1 {
			return Pcmd{Kind: PcmdForceOpen, Name: parts[1]}, true
		} else {
			return Pcmd{Kind: PcmdForceOpen}, true
		}
	case "r":
		return Pcmd{Kind: PcmdRead}, true
	case "n", "next":
		return Pcmd{Kind: PcmdNext}, true
	case "prev", "previous":
		return Pcmd{Kind: PcmdPrev}, true

	case "sh", "shell":
		return Pcmd{Kind: PcmdShell}, true

	case "wa":
		return Pcmd{Kind: PcmdSaveAll}, true
	case "qa":
		return Pcmd{Kind: PcmdQuitAll}, true
	case "qa!":
		return Pcmd{Kind: PcmdForceQuitAll}, true

	case "set":
		if len(parts) < 2 {
			// TODO show variables
			return Pcmd{Kind: PcmdInvalid}, false
		}
		if strings.HasPrefix(parts[1], "ts=") {
			ns := parts[1][3:]
			n, err := strconv.ParseUint(ns, 10, 16)
			if err != nil {
				// XXX
				//ed.Ring("set: %s option: %s is an illegal number.", ns, ns)
				return Pcmd{Kind: PcmdInvalid}, false
			}
			return Pcmd{Kind: PcmdTabStop, Num: int(n)}, true
		}
		switch parts[1] {
		case "ai", "autoindent":
			return Pcmd{Kind: PcmdAutoIndent}, true
		case "noai", "noautoindent":
			return Pcmd{Kind: PcmdNoAutoIndent}, true
		}
		// TODO set all
		ed.Ring(
			"set no %s option: 'set all' gives all option values.",
			parts[1],
		)
		return Pcmd{Kind: PcmdInvalid}, false

	case "col", "colors", "colorscheme":
		if len(parts) > 1 {
			return Pcmd{Kind: PcmdColors, Name: parts[1]}, true
		} else {
			return Pcmd{Kind: PcmdColors}, true
		}

	default:
		ed.Ring("The %s command is unknown.", parts[0])
		return Pcmd{Kind: PcmdInvalid}, false
	}
}
