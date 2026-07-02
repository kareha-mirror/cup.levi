package editor

import (
	"unicode/utf8"
)

func (ed *Editor) RegMode(name rune) KillMode {
	mode, err := ed.regs.Mode(name)
	if err != nil {
		ed.Error("%v", err)
		return KillNone
	}
	return mode
}

func (ed *Editor) RegKilled(name rune) []string {
	killed, err := ed.regs.Killed(name)
	if err != nil {
		ed.Error("%v", err)
		return []string{""}
	}
	return killed
}

func (ed *Editor) ApplyRegLines(name rune, killed []string) bool {
	err := ed.regs.ApplyLines(name, killed, ed.Buf().CRLF)
	if err != nil {
		ed.Error("%v", err)
		return false
	}

	numLines := len(killed)
	if numLines >= 5 {
		ed.Message("%d lines yanked", numLines)
	}
	return true
}

func (ed *Editor) ApplyRegRunes(name rune, killed []string) bool {
	err := ed.regs.ApplyRunes(name, killed, ed.Buf().CRLF)
	if err != nil {
		ed.Error("%v", err)
		return false
	}

	numLines := len(killed)
	if numLines >= 5 {
		ed.Message("%d lines yanked", numLines)
	} else if numLines == 1 {
		rc := utf8.RuneCountInString(killed[0])
		if rc >= 25 {
			ed.Message("%d bytes, %d runes yanked", len(killed[0]), rc)
		}
	}
	return true
}
