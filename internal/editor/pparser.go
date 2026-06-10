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
		return Pcmd{Kind: PcmdSave}, true
	case "w!":
		return Pcmd{Kind: PcmdForceSave}, true
	case "q":
		return Pcmd{Kind: PcmdQuit}, true
	case "q!":
		return Pcmd{Kind: PcmdForceQuit}, true
	case "e":
		return Pcmd{Kind: PcmdOpen}, true
	case "e!":
		return Pcmd{Kind: PcmdForceOpen}, true
	case "r":
		return Pcmd{Kind: PcmdRead}, true
	case "n":
		return Pcmd{Kind: PcmdNext}, true
	case "prev":
		return Pcmd{Kind: PcmdPrev}, true

	case "sh":
		return Pcmd{Kind: PcmdShell}, true

	case "wa":
		return Pcmd{Kind: PcmdSaveAll}, true
	case "qa":
		return Pcmd{Kind: PcmdQuitAll}, true
	case "qa!":
		return Pcmd{Kind: PcmdForceQuitAll}, true

	default:
		return Pcmd{Kind: PcmdInvalid}, false
	}
}
