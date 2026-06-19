package colors

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

//go:embed user.yaml
var userConfig string

//go:embed colors/*.yaml
var embedFS embed.FS

type Colors struct {
	Buffer  termi.ColorPair
	Status  termi.ColorPair
	Current termi.ColorPair
	Border  termi.ColorPair
}

type config struct {
	Buffer  string `yaml:"buffer"`
	Status  string `yaml:"status"`
	Current string `yaml:"current"`
	Border  string `yaml:"border"`
}

func (cfg *config) Colors() (*Colors, error) {
	buffer, err := termi.ParseColorPair(cfg.Buffer)
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
		Buffer:  buffer,
		Status:  status,
		Current: current,
		Border:  border,
	}, nil
}

type List struct {
	dir string

	Users    []string
	Builtins []string

	uMap map[string]bool
	bMap map[string]bool
}

func LoadList(dir string) *List {
	userDir := filepath.Join(dir, "colors")
	os.Mkdir(userDir, 0777)
	userPath := filepath.Join(userDir, "user.yaml")
	_, err := os.Stat(userPath)
	if err != nil {
		os.WriteFile(userPath, []byte(userConfig), 0666)
	}

	users, _ := listUser(dir)
	builtins, _ := listBuiltin()
	uMap := map[string]bool{}
	for _, name := range users {
		uMap[name] = true
	}
	bMap := map[string]bool{}
	for _, name := range builtins {
		bMap[name] = true
	}
	return &List{dir, users, builtins, uMap, bMap}
}

func (list *List) Load(name string) (*Colors, error) {
	if list.uMap[name] {
		path := filepath.Join(list.dir, "colors", name+".yaml")
		cfg, err := loadUser(path)
		if err != nil {
			return nil, err
		}
		return cfg.Colors()
	}
	if list.bMap[name] {
		path := filepath.Join("colors", name+".yaml")
		cfg, err := loadBuiltin(path)
		if err != nil {
			return nil, err
		}
		return cfg.Colors()
	}
	return nil, fmt.Errorf("not found")
}

func listUser(dir string) ([]string, error) {
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

func listBuiltin() ([]string, error) {
	entries, err := fs.ReadDir(embedFS, "colors")
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

func parseConfig(b []byte) (*config, error) {
	var cfg config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func loadUser(path string) (*config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseConfig(b)
}

func loadBuiltin(path string) (*config, error) {
	b, err := fs.ReadFile(embedFS, path)
	if err != nil {
		return nil, err
	}
	return parseConfig(b)
}

func Parse(s string) (*Colors, error) {
	cfg, err := parseConfig([]byte(s))
	if err != nil {
		return nil, err
	}
	return cfg.Colors()
}
