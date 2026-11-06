package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/regs"
)

func (ed *Editor) CopyMode(name rune) regs.Mode {
	mode, err := ed.regs.Mode(name)
	if err != nil {
		ed.Error("%v", err)
		return regs.None
	}
	return mode
}

func (ed *Editor) CopiedContent(name rune) []string {
	copied, err := ed.regs.Content(name)
	if err != nil {
		ed.Error("%v", err)
		return []string{""}
	}
	if len(copied) < 1 {
		return []string{""}
	}
	return copied
}

func (ed *Editor) StoreLines(name rune, copied []string) bool {
	err := ed.regs.ApplyLines(name, copied)
	if err != nil {
		ed.Error("%v", err)
		return false
	}

	numLines := len(copied)
	if numLines >= 5 {
		ed.Message("%d lines yanked", numLines)
	}
	return true
}

func (ed *Editor) StoreRunes(name rune, copied []string) bool {
	err := ed.regs.ApplyRunes(name, copied)
	if err != nil {
		ed.Error("%v", err)
		return false
	}

	numLines := len(copied)
	if numLines >= 5 {
		ed.Message("%d lines yanked", numLines)
	} else if numLines == 1 {
		rc := utf8.RuneCountInString(copied[0])
		if rc >= 25 {
			ed.Message("%d bytes, %d runes yanked", len(copied[0]), rc)
		}
	}
	return true
}
