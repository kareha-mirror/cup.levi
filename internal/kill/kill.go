package kill

import (
	"fmt"
	"strings"

	"golang.design/x/clipboard"

	"tea.kareha.org/cup/levi/internal/buf"
)

var ClipboardInitialized = false

type Mode int

const (
	None Mode = iota
	Lines
	Runes
)

type slot struct {
	mode    Mode
	content []string
	shared  bool
}

type Store struct {
	defMode    Mode
	defContent []string
	set        map[rune]slot
	cfgDir     string
}

func IsValidName(name rune) bool {
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

func NormalizeName(name rune) rune {
	if name >= 'A' && name <= 'Z' {
		return name + 'a' - 'A'
	}
	return name
}

func (s *Store) setSlot(name rune, sl slot) bool {
	if !IsValidName(name) {
		return false
	}
	if name == 0 {
		s.defMode = sl.mode
		s.defContent = sl.content
		return true
	}
	if name == '+' {
		// not supported
		return true
	}
	name = NormalizeName(name)
	if s.set == nil {
		s.set = make(map[rune]slot)
	}
	s.set[name] = sl
	return true
}

func (s *Store) Mode(name rune) (Mode, error) {
	if !IsValidName(name) {
		return None, fmt.Errorf("Invalid slot name")
	}
	if name == 0 {
		return s.defMode, nil
	}
	if name == '+' {
		return Runes, nil
	}
	name = NormalizeName(name)
	sl, ok := s.set[name]
	if !ok {
		return None, nil
	}
	if sl.shared {
		err := s.loadMeta(name)
		if err != nil {
			return None, err
		}
		sl, _ = s.set[name]
	}
	return sl.mode, nil
}

func ensureClipboard() error {
	if ClipboardInitialized {
		return nil
	}
	if err := clipboard.Init(); err != nil {
		return err
	}
	ClipboardInitialized = true
	return nil
}

func (s *Store) Content(name rune) ([]string, error) {
	if !IsValidName(name) {
		return nil, fmt.Errorf("Invalid slot name")
	}
	if name == 0 {
		return s.defContent, nil
	}
	if name == '+' {
		err := ensureClipboard()
		if err != nil {
			return nil, err
		}
		text := string(clipboard.Read(clipboard.FmtText))
		text = strings.ReplaceAll(text, "\r\n", "\n")
		return strings.Split(text, "\n"), nil
	}
	name = NormalizeName(name)
	sl, ok := s.set[name]
	if !ok {
		return nil, nil
	}
	if sl.shared {
		err := s.loadContent(name)
		if err != nil {
			return nil, err
		}
		sl, _ = s.set[name]
	}
	return sl.content, nil
}

func (s *Store) Shared(name rune) bool {
	if !IsValidName(name) {
		return false
	}
	if name == 0 {
		return false
	}
	if name == '+' {
		return false
	}
	name = NormalizeName(name)
	sl, ok := s.set[name]
	if !ok {
		return false
	}
	return sl.shared
}

func (s *Store) SetShared(name rune, shared bool) bool {
	if !IsValidName(name) {
		return false
	}
	if name == 0 {
		return true
	}
	if name == '+' {
		return true
	}
	name = NormalizeName(name)
	sl, _ := s.set[name]
	sl.shared = shared
	s.setSlot(name, sl)
	return true
}

func (s *Store) SetLines(name rune, killed []string) error {
	if !IsValidName(name) {
		return fmt.Errorf("Invalid slot name")
	}
	lines := append([]string{}, killed...)
	s.defMode = Lines
	s.defContent = lines
	if name == 0 {
		return nil
	}
	if name == '+' {
		lines := append([]string{}, killed...)
		lines = append(lines, "")
		return s.SetRunes(name, lines)
	}
	name = NormalizeName(name)
	sl, _ := s.set[name]
	sl.mode = Lines
	sl.content = append([]string{}, lines...)
	s.setSlot(name, sl)
	if sl.shared {
		return s.save(name)
	}
	return nil
}

func (s *Store) SetRunes(name rune, killed []string) error {
	if !IsValidName(name) {
		return fmt.Errorf("Invalid slot name")
	}
	lines := append([]string{}, killed...)
	s.defMode = Runes
	s.defContent = lines
	if name == 0 {
		return nil
	}
	if name == '+' {
		err := ensureClipboard()
		if err != nil {
			return err
		}
		text := strings.Join(killed, buf.LineSep(ClipboardCRLF))
		clipboard.Write(clipboard.FmtText, []byte(text))
		return nil
	}
	name = NormalizeName(name)
	sl, _ := s.set[name]
	sl.mode = Runes
	sl.content = append([]string{}, lines...)
	s.setSlot(name, sl)
	if sl.shared {
		return s.save(name)
	}
	return nil
}

func (s *Store) AddLines(name rune, killed []string) error {
	if !IsValidName(name) {
		return fmt.Errorf("Invalid slot name")
	}
	if name == 0 {
		return nil
	}
	if name == '+' {
		return fmt.Errorf("Not supported")
	}
	name = NormalizeName(name)
	sl, ok := s.set[name]
	if !ok {
		sl = slot{
			mode:    Lines,
			content: append([]string{}, killed...),
			shared:  false,
		}
		s.setSlot(name, sl)
		return nil
	}
	if sl.mode == Lines {
		sl.content = append(sl.content, killed...)
	} else {
		sl.mode = Lines
		sl.content = append([]string{}, killed...)
	}
	s.setSlot(name, sl)
	if sl.shared {
		return s.save(name)
	}
	return nil
}

func (s *Store) AddRunes(name rune, killed []string) error {
	if !IsValidName(name) {
		return fmt.Errorf("Invalid slot name")
	}
	if name == 0 {
		return nil
	}
	if name == '+' {
		return fmt.Errorf("Not supported")
	}
	name = NormalizeName(name)
	sl, ok := s.set[name]
	if !ok {
		sl = slot{
			mode:    Runes,
			content: append([]string{}, killed...),
			shared:  false,
		}
		s.setSlot(name, sl)
		return nil
	}
	if sl.mode == Runes {
		if len(killed) > 0 {
			lines := []string{}
			if len(sl.content) > 0 {
				lines = append(lines, sl.content[:len(sl.content)-1]...)
				line := sl.content[len(sl.content)-1] + killed[0]
				lines = append(lines, line)
			} else {
				lines = append(lines, killed[0])
			}
			if 1 < len(killed) {
				lines = append(lines, killed[1:]...)
			}
			sl.content = lines
		}
	} else {
		sl.mode = Runes
		sl.content = append([]string{}, killed...)
	}
	s.setSlot(name, sl)
	if sl.shared {
		return s.save(name)
	}
	return nil
}

func (s *Store) ApplyLines(name rune, killed []string) error {
	if !IsValidName(name) {
		return fmt.Errorf("Invalid slot name")
	}
	normalized := NormalizeName(name)
	if normalized != name {
		return s.AddLines(name, killed)
	}
	return s.SetLines(name, killed)
}

func (s *Store) ApplyRunes(name rune, killed []string) error {
	if !IsValidName(name) {
		return fmt.Errorf("Invalid slot name")
	}
	normalized := NormalizeName(name)
	if normalized != name {
		return s.AddRunes(name, killed)
	}
	return s.SetRunes(name, killed)
}
