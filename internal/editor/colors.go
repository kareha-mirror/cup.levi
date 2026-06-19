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
	Text    string `yaml:"text"`
	Status  string `yaml:"status"`
	Current string `yaml:"current"`
	Border  string `yaml:"border"`
}

func splitColorPairString(s string) (string, string) {
	parts := strings.Split(s, ",")
	if len(parts[0]) < 1 {
		parts[0] = "default"
	}
	if len(parts) < 2 {
		parts = append(parts, "default")
	}
	return parts[0], parts[1]
}

func parseColorPair(s string) (termi.Color, termi.Color, error) {
	fgStr, bgStr := splitColorPairString(s)
	fg, err := termi.ParseColor(fgStr)
	if err != nil {
		return termi.Color{}, termi.Color{}, err
	}
	bg, err := termi.ParseColor(bgStr)
	if err != nil {
		return termi.Color{}, termi.Color{}, err
	}
	return fg, bg, nil
}

func (cfg *ColorsConfig) Colors() (*Colors, error) {
	textFg, textBg, err := parseColorPair(cfg.Text)
	if err != nil {
		return nil, err
	}
	statusFg, statusBg, err := parseColorPair(cfg.Status)
	if err != nil {
		return nil, err
	}
	currentFg, currentBg, err := parseColorPair(cfg.Current)
	if err != nil {
		return nil, err
	}
	borderFg, borderBg, err := parseColorPair(cfg.Border)
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
