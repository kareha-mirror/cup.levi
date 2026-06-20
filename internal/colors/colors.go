package colors

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
	"tea.kareha.org/cup/termi"
)

//go:embed custom.yaml
var customConfig string

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

	customs  map[string]bool
	builtins map[string]bool

	Total []string
}

func LoadList(dir string) *List {
	customDir := filepath.Join(dir, "colors")
	os.Mkdir(customDir, 0777)
	customPath := filepath.Join(customDir, "custom.yaml")
	_, err := os.Stat(customPath)
	if err != nil {
		os.WriteFile(customPath, []byte(customConfig), 0666)
	}

	cList, _ := listCustom(dir)
	bList, _ := listBuiltin()
	customs := map[string]bool{}
	for _, name := range cList {
		customs[name] = true
	}
	builtins := map[string]bool{}
	for _, name := range bList {
		builtins[name] = true
	}

	tMap := map[string]bool{}
	for name := range customs {
		tMap[name] = true
	}
	for name := range builtins {
		tMap[name] = true
	}
	total := []string{}
	for name := range tMap {
		total = append(total, name)
	}
	sort.Strings(total)

	return &List{dir, customs, builtins, total}
}

func (list *List) Load(name string) (*Colors, error) {
	if list.customs[name] {
		path := filepath.Join(list.dir, "colors", name+".yaml")
		cfg, err := loadCustom(path)
		if err != nil {
			return nil, err
		}
		return cfg.Colors()
	}
	if list.builtins[name] {
		path := filepath.Join("colors", name+".yaml")
		cfg, err := loadBuiltin(path)
		if err != nil {
			return nil, err
		}
		return cfg.Colors()
	}
	return nil, fmt.Errorf("not found")
}

func listCustom(dir string) ([]string, error) {
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

func loadCustom(path string) (*config, error) {
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
