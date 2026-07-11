package kill

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"tea.kareha.org/cup/termi/lock"
)

var SharedDirName = "kills"

var ReadFile func(string) ([]byte, error) = os.ReadFile
var WriteFile func(string, []byte, os.FileMode) error = os.WriteFile

type meta struct {
	lines bool `yaml:"lines"`
}

func loadMeta(path string) (meta, error) {
	data, err := ReadFile(path)
	if err != nil {
		return meta{}, err
	}
	var m meta
	if err := yaml.Unmarshal(data, &m); err != nil {
		return meta{}, err
	}
	return m, nil
}

func (m *meta) save(path string) error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}
	err = WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Init(cfgDir string, shared string) {
	s.cfgDir = cfgDir
	for _, sl := range s.set {
		sl.shared = false
	}
	for _, name := range shared {
		if !IsValidName(name) {
			continue
		}
		s.SetShared(name, true)
	}
}

func (s *Store) sharedMetaPath(name rune) string {
	return filepath.Join(s.cfgDir, SharedDirName, string(name)+".yaml")
}

func (s *Store) sharedTextPath(name rune) string {
	return filepath.Join(s.cfgDir, SharedDirName, string(name)+".txt")
}

func (s *Store) loadMeta(name rune) error {
	sl, ok := s.set[name]
	if !ok {
		return fmt.Errorf("slot not found")
	}

	err := lock.Lock(s.cfgDir)
	if err != nil {
		return err
	}
	defer lock.Unlock(s.cfgDir)

	metaPath := s.sharedMetaPath(name)
	m, err := loadMeta(metaPath)
	if err != nil {
		return err
	}

	if m.lines {
		sl.mode = Lines
	} else {
		sl.mode = Runes
	}

	sl.content = nil

	s.setSlot(name, sl)
	return nil
}

func (s *Store) loadContent(name rune) error {
	sl, ok := s.set[name]
	if !ok {
		return fmt.Errorf("slot not found")
	}

	err := lock.Lock(s.cfgDir)
	if err != nil {
		return err
	}
	defer lock.Unlock(s.cfgDir)

	textPath := s.sharedTextPath(name)
	data, err := ReadFile(textPath)
	if err != nil {
		return err
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if text != "" && text[len(text)-1] == '\n' {
		text = text[:len(text)-1]
	}
	sl.content = strings.Split(text, "\n")

	s.setSlot(name, sl)
	return nil
}

func (s *Store) save(name rune) error {
	sl, ok := s.set[name]
	if !ok {
		return fmt.Errorf("slot not found")
	}

	err := lock.Lock(s.cfgDir)
	if err != nil {
		return err
	}
	defer lock.Unlock(s.cfgDir)

	m := meta{lines: sl.mode == Lines}
	metaPath := s.sharedMetaPath(name)
	err = m.save(metaPath)
	if err != nil {
		return err
	}

	textPath := s.sharedTextPath(name)
	err = os.MkdirAll(filepath.Dir(textPath), 0777)
	if err != nil {
		return err
	}
	text := strings.Join(sl.content, "\n") + "\n"
	err = WriteFile(textPath, []byte(text), 0666)
	if err != nil {
		return err
	}

	return nil
}
