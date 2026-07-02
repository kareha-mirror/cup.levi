package editor

import (
	"fmt"
	"strings"

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
	Map       map[rune]Reg
	cfgDir    string
}

func IsValidRegName(name rune) bool {
	if name == 0 { // default
		return true
	}
	if name >= 'a' && name <= 'z' { // normal
		return true
	}
	if name >= 'A' && name <= 'Z' { // append
		return true
	}
	if name == '+' { // clipboard
		return true
	}
	return false
}

func NormalizeRegName(name rune) rune {
	if name >= 'A' && name <= 'Z' {
		return name + 'a' - 'A'
	}
	return name
}

func (regs *Regs) SetReg(name rune, reg Reg) bool {
	if !IsValidRegName(name) {
		return false
	}
	if name == 0 {
		regs.DefMode = reg.Mode
		regs.DefKilled = reg.Killed
		return true
	}
	if name == '+' {
		// omitted
		return true
	}
	name = NormalizeRegName(name)
	if regs.Map == nil {
		regs.Map = make(map[rune]Reg)
	}
	regs.Map[name] = reg
	return true
}

func (regs *Regs) Mode(name rune) (KillMode, error) {
	if !IsValidRegName(name) {
		return KillNone, fmt.Errorf("Invalid reg name")
	}
	if name == 0 {
		return regs.DefMode, nil
	}
	if name == '+' {
		return KillRunes, nil
	}
	name = NormalizeRegName(name)
	reg, ok := regs.Map[name]
	if !ok {
		return KillNone, nil
	}
	if reg.Shared {
		err := regs.LoadMeta(name)
		if err != nil {
			return KillNone, err
		}
		reg, _ = regs.Map[name]
	}
	return reg.Mode, nil
}

var clipboardInitialized = false

func EnsureClipboard() error {
	if clipboardInitialized {
		return nil
	}
	if err := clipboard.Init(); err != nil {
		return err
	}
	clipboardInitialized = true
	return nil
}

func (regs *Regs) Killed(name rune) ([]string, error) {
	if !IsValidRegName(name) {
		return nil, fmt.Errorf("Invalid reg name")
	}
	if name == 0 {
		return regs.DefKilled, nil
	}
	if name == '+' {
		err := EnsureClipboard()
		if err != nil {
			return nil, err
		}
		text := string(clipboard.Read(clipboard.FmtText))
		text = strings.ReplaceAll(text, "\r\n", "\n")
		return strings.Split(text, "\n"), nil
	}
	name = NormalizeRegName(name)
	reg, ok := regs.Map[name]
	if !ok {
		return nil, nil
	}
	if reg.Shared {
		err := regs.LoadContent(name)
		if err != nil {
			return nil, err
		}
		reg, _ = regs.Map[name]
	}
	return reg.Killed, nil
}

func (regs *Regs) Shared(name rune) bool {
	if !IsValidRegName(name) {
		return false
	}
	if name == 0 {
		return false
	}
	if name == '+' {
		return false
	}
	name = NormalizeRegName(name)
	reg, ok := regs.Map[name]
	if !ok {
		return false
	}
	return reg.Shared
}

func (regs *Regs) SetShared(name rune, shared bool) bool {
	if !IsValidRegName(name) {
		return false
	}
	if name == 0 {
		return true
	}
	if name == '+' {
		return true
	}
	name = NormalizeRegName(name)
	reg, _ := regs.Map[name]
	reg.Shared = shared
	regs.SetReg(name, reg)
	return true
}

func (regs *Regs) SetLines(name rune, killed []string, crlf bool) error {
	if !IsValidRegName(name) {
		return fmt.Errorf("Invalid reg name")
	}
	lines := append([]string{}, killed...)
	regs.DefMode = KillLines
	regs.DefKilled = lines
	if name == 0 {
		return nil
	}
	if name == '+' {
		lines := append([]string{}, killed...)
		lines = append(lines, "")
		return regs.SetRunes(name, lines, crlf)
	}
	name = NormalizeRegName(name)
	reg, _ := regs.Map[name]
	reg.Mode = KillLines
	reg.Killed = lines
	regs.SetReg(name, reg)
	if reg.Shared {
		return regs.Save(name)
	}
	return nil
}

func (regs *Regs) SetRunes(name rune, killed []string, crlf bool) error {
	if !IsValidRegName(name) {
		return fmt.Errorf("Invalid reg name")
	}
	lines := append([]string{}, killed...)
	regs.DefMode = KillRunes
	regs.DefKilled = lines
	if name == 0 {
		return nil
	}
	if name == '+' {
		err := EnsureClipboard()
		if err != nil {
			return err
		}
		text := strings.Join(killed, buf.LineSep(crlf))
		clipboard.Write(clipboard.FmtText, []byte(text))
		return nil
	}
	name = NormalizeRegName(name)
	reg, _ := regs.Map[name]
	reg.Mode = KillRunes
	reg.Killed = lines
	regs.SetReg(name, reg)
	if reg.Shared {
		return regs.Save(name)
	}
	return nil
}

func (regs *Regs) AddLines(name rune, killed []string) error {
	if !IsValidRegName(name) {
		return fmt.Errorf("Invalid reg name")
	}
	if name == 0 {
		return nil
	}
	if name == '+' {
		return fmt.Errorf("Not supported")
	}
	name = NormalizeRegName(name)
	reg, ok := regs.Map[name]
	if !ok {
		reg = Reg{
			Mode:   KillLines,
			Killed: append([]string{}, killed...),
			Shared: false,
		}
		regs.SetReg(name, reg)
		return nil
	}
	if reg.Mode == KillLines {
		reg.Killed = append(reg.Killed, killed...)
	} else {
		reg.Mode = KillLines
		reg.Killed = append([]string{}, killed...)
	}
	regs.SetReg(name, reg)
	if reg.Shared {
		return regs.Save(name)
	}
	return nil
}

func (regs *Regs) AddRunes(name rune, killed []string) error {
	if !IsValidRegName(name) {
		return fmt.Errorf("Invalid reg name")
	}
	if name == 0 {
		return nil
	}
	if name == '+' {
		return fmt.Errorf("Not supported")
	}
	name = NormalizeRegName(name)
	reg, ok := regs.Map[name]
	if !ok {
		reg = Reg{
			Mode:   KillRunes,
			Killed: append([]string{}, killed...),
			Shared: false,
		}
		regs.SetReg(name, reg)
		return nil
	}
	if reg.Mode == KillRunes {
		lines := append([]string{}, reg.Killed[:len(reg.Killed)-1]...)
		line := reg.Killed[len(reg.Killed)-1] + killed[0]
		lines = append(lines, line)
		if 1 < len(killed) {
			lines = append(lines, killed[1:]...)
		}
		reg.Killed = lines
	} else {
		reg.Mode = KillRunes
		reg.Killed = append([]string{}, killed...)
	}
	regs.SetReg(name, reg)
	if reg.Shared {
		return regs.Save(name)
	}
	return nil
}

func (regs *Regs) ApplyLines(name rune, killed []string, crlf bool) error {
	if !IsValidRegName(name) {
		return fmt.Errorf("Invalid reg name")
	}
	normalized := NormalizeRegName(name)
	if normalized != name {
		return regs.AddLines(name, killed)
	}
	return regs.SetLines(name, killed, crlf)
}

func (regs *Regs) ApplyRunes(name rune, killed []string, crlf bool) error {
	if !IsValidRegName(name) {
		return fmt.Errorf("Invalid reg name")
	}
	normalized := NormalizeRegName(name)
	if normalized != name {
		return regs.AddRunes(name, killed)
	}
	return regs.SetRunes(name, killed, crlf)
}
