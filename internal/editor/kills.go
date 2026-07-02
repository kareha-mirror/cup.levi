package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/kill"
)

func (ed *Editor) KillMode(name rune) kill.Mode {
	mode, err := ed.kills.Mode(name)
	if err != nil {
		ed.Error("%v", err)
		return kill.None
	}
	return mode
}

func (ed *Editor) KilledContent(name rune) []string {
	killed, err := ed.kills.Content(name)
	if err != nil {
		ed.Error("%v", err)
		return []string{""}
	}
	if len(killed) < 1 {
		return []string{""}
	}
	return killed
}

func (ed *Editor) StoreLines(name rune, killed []string) bool {
	err := ed.kills.ApplyLines(name, killed)
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

func (ed *Editor) StoreRunes(name rune, killed []string) bool {
	err := ed.kills.ApplyRunes(name, killed)
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
