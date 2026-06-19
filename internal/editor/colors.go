package editor

import (
	"embed"
	"io/fs"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	"tea.kareha.org/cup/termi"
)

//go:embed colors/*.yaml
var colorsFS embed.FS

type Colors struct {
	Text    termi.ColorPair
	Status  termi.ColorPair
	Current termi.ColorPair
	Border  termi.ColorPair
}

type ColorsConfig struct {
	Text    string `yaml:"text"`
	Status  string `yaml:"status"`
	Current string `yaml:"current"`
	Border  string `yaml:"border"`
}

func (cfg *ColorsConfig) Colors() (*Colors, error) {
	text, err := termi.ParseColorPair(cfg.Text)
	if err != nil {
		return nil, err
	}
	status, err := termi.ParseColorPair(cfg.Status)
	if err != nil {
		return nil, err
	}
	current, err := termi.ParseColorPair(cfg.Current)
	if err != nil {
		return nil, err
	}
	border, err := termi.ParseColorPair(cfg.Border)
	if err != nil {
		return nil, err
	}

	return &Colors{
		Text:    text,
		Status:  status,
		Current: current,
		Border:  border,
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

func LoadColorsConfigFromString(s string) (*ColorsConfig, error) {
	var cfg ColorsConfig
	if err := yaml.Unmarshal([]byte(s), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadEmbeddedColorsConfig(path string) (*ColorsConfig, error) {
	data, err := fs.ReadFile(colorsFS, path)
	if err != nil {
		return nil, err
	}
	return LoadColorsConfigFromString(string(data))
}

func LoadColorsConfig(path string) (*ColorsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadColorsConfigFromString(string(data))
}
