package editor

import (
	"strings"
	"unicode/utf8"

	"golang.design/x/clipboard"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/config"
)

type KillMode int

const (
	KillNone = iota
	KillLines
	KillRunes
)

type Reg struct {
	Mode   KillMode
	Killed []string
	Shared bool
}

type Regs struct {
	DefMode   KillMode
	DefKilled []string
	Map       map[rune]*Reg
}

func (regs *Regs) SetReg(name rune, reg *Reg) {
	if regs.Map == nil {
		regs.Map = make(map[rune]*Reg)
	}
	regs.Map[name] = reg
}

func (regs *Regs) Mode(name rune) KillMode {
	if name == 0 {
		return regs.DefMode
	}
	reg, ok := regs.Map[name]
	if !ok {
		return KillNone
	}
	return reg.Mode
}

func (regs *Regs) Killed(name rune) []string {
	if name == 0 {
		return regs.DefKilled
	}
	reg, ok := regs.Map[name]
	if !ok {
		return nil
	}
	return reg.Killed
}

func (regs *Regs) Shared(name rune) bool {
	if name == 0 {
		return false
	}
	reg, ok := regs.Map[name]
	if !ok {
		return false
	}
	return reg.Shared
}

func IsValidRegName(name rune) bool {
	return name >= 'a' && name <= 'z'
}

func (regs *Regs) SetShared(name rune, shared bool) {
	if !IsValidRegName(name) {
		return
	}
	reg, ok := regs.Map[name]
	if !ok {
		reg = &Reg{
			Mode:   KillNone,
			Killed: nil,
			Shared: shared,
		}
		regs.SetReg(name, reg)
		return
	}
	reg.Shared = shared
}

func ToDestRegName(name rune) rune {
	if name < 'A' || name > 'Z' {
		return 0
	}
	return name + 'a' - 'A'
}

func (regs *Regs) SetLines(name rune, killed []string) {
	lines := append([]string{}, killed...)
	regs.DefMode = KillLines
	regs.DefKilled = lines

	if !IsValidRegName(name) {
		return
	}

	reg, ok := regs.Map[name]
	if !ok {
		reg = &Reg{
			Mode:   KillLines,
			Killed: lines,
			Shared: false,
		}
		regs.SetReg(name, reg)
		return
	}
	reg.Mode = KillLines
	reg.Killed = lines
}

func (regs *Regs) SetRunes(name rune, killed []string) {
	lines := append([]string{}, killed...)
	regs.DefMode = KillRunes
	regs.DefKilled = lines

	if !IsValidRegName(name) {
		return
	}

	reg, ok := regs.Map[name]
	if !ok {
		reg = &Reg{
			Mode:   KillRunes,
			Killed: lines,
			Shared: false,
		}
		regs.SetReg(name, reg)
		return
	}
	reg.Mode = KillRunes
	reg.Killed = lines
}

func (regs *Regs) AddLines(name rune, killed []string) {
	name = ToDestRegName(name)
	if name == 0 {
		return
	}
	reg, ok := regs.Map[name]
	if !ok {
		reg = &Reg{
			Mode:   KillLines,
			Killed: append([]string{}, killed...),
			Shared: false,
		}
		regs.SetReg(name, reg)
		return
	}
	if reg.Mode == KillLines {
		reg.Killed = append(reg.Killed, killed...)
		return
	}
	reg.Mode = KillLines
	reg.Killed = append([]string{}, killed...)
}

func (regs *Regs) AddRunes(name rune, killed []string) {
	name = ToDestRegName(name)
	if name == 0 {
		return
	}
	reg, ok := regs.Map[name]
	if !ok {
		reg = &Reg{
			Mode:   KillRunes,
			Killed: append([]string{}, killed...),
			Shared: false,
		}
		regs.SetReg(name, reg)
		return
	}
	if reg.Mode == KillRunes {
		lines := append([]string{}, reg.Killed[:len(reg.Killed)-1]...)
		line := reg.Killed[len(reg.Killed)-1] + killed[0]
		lines = append(lines, line)
		if 1 < len(killed) {
			lines = append(lines, killed[1:]...)
		}
		reg.Killed = lines
		return
	}
	reg.Mode = KillRunes
	reg.Killed = append([]string{}, killed...)
}

func (regs *Regs) ApplyLines(name rune, killed []string) {
	addName := ToDestRegName(name)
	if addName != 0 {
		regs.AddLines(name, killed)
		return
	}
	regs.SetLines(name, killed)
}

func (regs *Regs) ApplyRunes(name rune, killed []string) {
	addName := ToDestRegName(name)
	if addName != 0 {
		regs.AddRunes(name, killed)
		return
	}
	regs.SetRunes(name, killed)
}

func (regs *Regs) SyncWithConfig(cfg *config.Config) {
	for _, reg := range regs.Map {
		reg.Shared = false
	}
	for _, name := range cfg.Shared {
		if !IsValidRegName(name) {
			continue
		}
		regs.SetShared(name, true)
	}
}

func (ed *Editor) EnsureClipboard() error {
	if ed.clipUsed {
		return nil
	}
	if err := clipboard.Init(); err != nil {
		return err
	}
	ed.clipUsed = true
	return nil
}

func (ed *Editor) RegMode(name rune) KillMode {
	if name == '+' {
		if err := ed.EnsureClipboard(); err != nil {
			ed.Error("%v", err)
			return KillNone
		}
		return KillRunes
	}
	return ed.regs.Mode(name)
}

func (ed *Editor) RegKilled(name rune) []string {
	if name == '+' {
		if err := ed.EnsureClipboard(); err != nil {
			ed.Error("%v", err)
			return []string{""}
		}
		text := string(clipboard.Read(clipboard.FmtText))
		text = strings.ReplaceAll(text, "\r\n", "\n")
		return strings.Split(text, "\n")
	}
	return ed.regs.Killed(name)
}

func (ed *Editor) ApplyRegLines(name rune, killed []string) bool {
	if name == '+' {
		lines := append([]string{}, killed...)
		lines = append(lines, "")
		if !ed.ApplyRegRunes(name, lines) {
			return false
		}
	} else {
		ed.regs.ApplyLines(name, killed)
	}
	numLines := len(killed)
	if numLines >= 5 {
		ed.Message("%d lines yanked", numLines)
	}
	return true
}

func (ed *Editor) ApplyRegRunes(name rune, killed []string) bool {
	if name == '+' {
		if err := ed.EnsureClipboard(); err != nil {
			ed.Error("%v", err)
			return false
		}
		text := strings.Join(killed, buf.LineSep(ed.Buf().CRLF))
		clipboard.Write(clipboard.FmtText, []byte(text))
	} else {
		ed.regs.ApplyRunes(name, killed)
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
