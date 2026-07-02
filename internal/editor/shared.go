package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"tea.kareha.org/cup/levi/internal/config"
)

type RegMeta struct {
	Lines bool `yaml:"lines"`
}

func LoadRegMeta(path string) (RegMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return RegMeta{}, err
	}
	var meta RegMeta
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return RegMeta{}, err
	}
	return meta, nil
}

func (meta *RegMeta) Save(path string) error {
	data, err := yaml.Marshal(meta)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (regs *Regs) SyncWithConfig(cfgDir string, cfg *config.Config) {
	regs.cfgDir = cfgDir
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

func (regs *Regs) SharedMetaPath(name rune) string {
	return filepath.Join(regs.cfgDir, "regs", string(name)+".yaml")
}

func (regs *Regs) SharedTextPath(name rune) string {
	return filepath.Join(regs.cfgDir, "regs", string(name)+".txt")
}

func (regs *Regs) Load(name rune) error {
	reg, ok := regs.Map[name]
	if !ok {
		return fmt.Errorf("Reg not found")
	}

	err := lock(regs.cfgDir)
	if err != nil {
		return err
	}
	defer Unlock(regs.cfgDir)

	metaPath := regs.SharedMetaPath(name)
	meta, err := LoadRegMeta(metaPath)
	if err != nil {
		return err
	}

	textPath := regs.SharedTextPath(name)
	data, err := os.ReadFile(textPath)
	if err != nil {
		return err
	}

	if meta.Lines {
		reg.Mode = KillLines
	} else {
		reg.Mode = KillRunes
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	reg.Killed = strings.Split(text, "\n")

	regs.Map[name] = reg

	return nil
}

func (regs *Regs) Save(name rune) error {
	reg, ok := regs.Map[name]
	if !ok {
		return fmt.Errorf("Reg not found")
	}

	err := lock(regs.cfgDir)
	if err != nil {
		return err
	}
	defer Unlock(regs.cfgDir)

	meta := RegMeta{Lines: reg.Mode == KillLines}
	metaPath := regs.SharedMetaPath(name)
	err = meta.Save(metaPath)
	if err != nil {
		return err
	}

	textPath := regs.SharedTextPath(name)
	err = os.MkdirAll(filepath.Dir(textPath), 0777)
	if err != nil {
		return err
	}
	text := strings.Join(reg.Killed, "\n")
	err = os.WriteFile(textPath, []byte(text), 0666)
	if err != nil {
		return err
	}

	return nil
}
