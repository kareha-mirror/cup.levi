package editor

import (
	"embed"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"
	"tea.kareha.org/cup/termi"
)

//go:embed colors/*.yaml
var colorsFS embed.FS

type Colors struct {
	TextFg    termi.Color
	TextBg    termi.Color
	StatusFg  termi.Color
	StatusBg  termi.Color
	CurrentFg termi.Color
	CurrentBg termi.Color
	BorderFg  termi.Color
	BorderBg  termi.Color
}

type ColorsConfig struct {
	TextFg    string `yaml:"text-fg"`
	TextBg    string `yaml:"text-bg"`
	StatusFg  string `yaml:"status-fg"`
	StatusBg  string `yaml:"status-bg"`
	CurrentFg string `yaml:"current-fg"`
	CurrentBg string `yaml:"current-bg"`
	BorderFg  string `yaml:"border-fg"`
	BorderBg  string `yaml:"border-bg"`
}

func (cfg *ColorsConfig) Colors() (*Colors, error) {
	textFg, err := termi.ParseColor(cfg.TextFg)
	if err != nil {
		return nil, err
	}
	textBg, err := termi.ParseColor(cfg.TextBg)
	if err != nil {
		return nil, err
	}
	statusFg, err := termi.ParseColor(cfg.StatusFg)
	if err != nil {
		return nil, err
	}
	statusBg, err := termi.ParseColor(cfg.StatusBg)
	if err != nil {
		return nil, err
	}
	currentFg, err := termi.ParseColor(cfg.CurrentFg)
	if err != nil {
		return nil, err
	}
	currentBg, err := termi.ParseColor(cfg.CurrentBg)
	if err != nil {
		return nil, err
	}
	borderFg, err := termi.ParseColor(cfg.BorderFg)
	if err != nil {
		return nil, err
	}
	borderBg, err := termi.ParseColor(cfg.BorderBg)
	if err != nil {
		return nil, err
	}

	return &Colors{
		TextFg:    textFg,
		TextBg:    textBg,
		StatusFg:  statusFg,
		StatusBg:  statusBg,
		CurrentFg: currentFg,
		CurrentBg: currentBg,
		BorderFg:  borderFg,
		BorderBg:  borderBg,
	}, nil
}

func ListEmbeddedColors() ([]string, error) {
	entries, err := fs.ReadDir(colorsFS, "colors")
	if err != nil {
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".yaml") {
			name = name[:len(name)-5]
		}
		names = append(names, name)
	}
	return names, nil
}

func LoadEmbeddedColorsConfig(path string) (*ColorsConfig, error) {
	data, err := fs.ReadFile(colorsFS, path)
	if err != nil {
		return nil, err
	}

	var cfg ColorsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
