package editor

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"tea.kareha.org/cup/termi"
)

//go:embed colors/*.yaml
var colorsFS embed.FS

const userColors = `text: "252,235"
status: "248,233"
current: "254,238"
border: "250,234"`

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

type ColorsList struct {
	dir      string
	User     []string
	Embedded []string
}

func LoadColorsList(dir string) *ColorsList {
	colorsDir := filepath.Join(dir, "colors")
	os.Mkdir(colorsDir, 0777)
	colorsPath := filepath.Join(colorsDir, "user.yaml")
	_, err := os.Stat(colorsPath)
	if err != nil {
		os.WriteFile(colorsPath, []byte(userColors), 0666)
	}

	user, err := ListUserColors(dir)
	if err != nil {
		user = []string{}
	}
	embedded, err := ListEmbeddedColors()
	if err != nil {
		embedded = []string{}
	}
	return &ColorsList{dir, user, embedded}
}

func (cl *ColorsList) Load(name string) (*Colors, error) {
	for _, n := range cl.User {
		if n == name {
			path := filepath.Join(cl.dir, "colors", name+".yaml")
			cfg, err := LoadColorsConfig(path)
			if err != nil {
				return nil, err
			}
			return cfg.Colors()
		}
	}
	for _, n := range cl.Embedded {
		if n == name {
			path := filepath.Join("colors", name+".yaml")
			cfg, err := LoadEmbeddedColorsConfig(path)
			if err != nil {
				return nil, err
			}
			return cfg.Colors()
		}
	}
	return nil, fmt.Errorf("not found")
}

func ListUserColors(dir string) ([]string, error) {
	colorsDir := filepath.Join(dir, "colors")
	entries, err := os.ReadDir(colorsDir)
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
