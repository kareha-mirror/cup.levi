package editor

import (
	"strings"
	"unicode/utf8"

	"golang.design/x/clipboard"

	"tea.kareha.org/cup/levi/internal/buf"
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
	Map       map[string]*Reg
}

func (regs *Regs) SetReg(name string, reg *Reg) {
	if regs.Map == nil {
		regs.Map = make(map[string]*Reg)
	}
	regs.Map[name] = reg
}

func (regs *Regs) Mode(name string) KillMode {
	if name == "" {
		return regs.DefMode
	}
	reg, ok := regs.Map[name]
	if !ok {
		return KillNone
	}
	return reg.Mode
}

func (regs *Regs) Killed(name string) []string {
	if name == "" {
		return regs.DefKilled
	}
	reg, ok := regs.Map[name]
	if !ok {
		return nil
	}
	return reg.Killed
}

func (regs *Regs) Shared(name string) bool {
	if name == "" {
		return false
	}
	reg, ok := regs.Map[name]
	if !ok {
		return false
	}
	return reg.Shared
}

func (regs *Regs) SetShared(name string, shared bool) {
	if name == "" {
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

func IsValidRegName(name string) bool {
	rc := utf8.RuneCountInString(name)
	if rc != 1 {
		return false
	}
	rs := []rune(name)
	r := rs[0]
	return r >= 'a' && r <= 'z'
}

func ToDestRegName(name string) string {
	if len(name) != 1 {
		return ""
	}
	rs := []rune(name)
	r := rs[0]
	if r < 'A' || r > 'Z' {
		return ""
	}
	return string(r + 'a' - 'A')
}

func (regs *Regs) SetLines(name string, killed []string) {
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

func (regs *Regs) SetRunes(name string, killed []string) {
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

func (regs *Regs) AddLines(name string, killed []string) {
	name = ToDestRegName(name)
	if name == "" {
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

func (regs *Regs) AddRunes(name string, killed []string) {
	name = ToDestRegName(name)
	if name == "" {
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

func (regs *Regs) ApplyLines(name string, killed []string) {
	addName := ToDestRegName(name)
	if addName != "" {
		regs.AddLines(name, killed)
		return
	}
	regs.SetLines(name, killed)
}

func (regs *Regs) ApplyRunes(name string, killed []string) {
	addName := ToDestRegName(name)
	if addName != "" {
		regs.AddRunes(name, killed)
		return
	}
	regs.SetRunes(name, killed)
}

func (regs *Regs) SyncWithConfig(cfg *Config) {
	for _, reg := range regs.Map {
		reg.Shared = false
	}
	for _, r := range cfg.Shared {
		name := string(r)
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

func (ed *Editor) RegMode(name string) KillMode {
	if name == "+" {
		if err := ed.EnsureClipboard(); err != nil {
			ed.Error("%v", err)
			return KillNone
		}
		return KillRunes
	}
	return ed.regs.Mode(name)
}

func (ed *Editor) RegKilled(name string) []string {
	if name == "+" {
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

func (ed *Editor) ApplyRegLines(name string, killed []string) {
	if name == "+" {
		lines := append([]string{}, killed...)
		lines = append(lines, "")
		ed.ApplyRegRunes(name, lines)
		return
	}
	ed.regs.ApplyLines(name, killed)
}

func (ed *Editor) ApplyRegRunes(name string, killed []string) {
	if name == "+" {
		if err := ed.EnsureClipboard(); err != nil {
			ed.Error("%v", err)
			return
		}
		text := strings.Join(killed, buf.LineSep(ed.Buf().CRLF))
		clipboard.Write(clipboard.FmtText, []byte(text))
		return
	}
	ed.regs.ApplyRunes(name, killed)
}
