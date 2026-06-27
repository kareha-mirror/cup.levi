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
var CustomConfig string

//go:embed colors/*.yaml
var BuiltinsFS embed.FS

type Colors struct {
	Buffer  termi.ColorPair
	Status  termi.ColorPair
	Current termi.ColorPair
	Border  termi.ColorPair
}

type Config struct {
	Buffer  string `yaml:"buffer"`
	Status  string `yaml:"status"`
	Current string `yaml:"current"`
	Border  string `yaml:"border"`
}

func (cfg *Config) Colors() (*Colors, error) {
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
	cfgDir string

	Customs  map[string]struct{}
	Builtins map[string]struct{}

	Names []string
}

func ListCustoms(cfgDir string) ([]string, error) {
	colorsDir := filepath.Join(cfgDir, "colors")
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

func ListBuiltins() ([]string, error) {
	entries, err := fs.ReadDir(BuiltinsFS, "colors")
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

func LoadList(cfgDir string) (*List, error) {
	customDir := filepath.Join(cfgDir, "colors")
	os.Mkdir(customDir, 0777)
	customPath := filepath.Join(customDir, "custom.yaml")
	_, err := os.Stat(customPath)
	if err != nil {
		os.WriteFile(customPath, []byte(CustomConfig), 0666)
	}

	cList, err := ListCustoms(cfgDir)
	if err != nil {
		return nil, err
	}
	bList, err := ListBuiltins()
	if err != nil {
		return nil, err
	}

	customs := map[string]struct{}{}
	for _, name := range cList {
		customs[name] = struct{}{}
	}
	builtins := map[string]struct{}{}
	for _, name := range bList {
		builtins[name] = struct{}{}
	}

	tMap := map[string]struct{}{}
	for name := range customs {
		tMap[name] = struct{}{}
	}
	for name := range builtins {
		tMap[name] = struct{}{}
	}
	total := []string{}
	for name := range tMap {
		total = append(total, name)
	}
	sort.Strings(total)

	return &List{cfgDir, customs, builtins, total}, nil
}

func ParseConfig(b []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadCustom(cfgDir string, name string) (*Config, error) {
	path := filepath.Join(cfgDir, "colors", name+".yaml")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseConfig(b)
}

func LoadBuiltin(name string) (*Config, error) {
	path := filepath.Join("colors", name+".yaml")
	b, err := fs.ReadFile(BuiltinsFS, path)
	if err != nil {
		return nil, err
	}
	return ParseConfig(b)
}

func (list *List) Load(name string) (*Colors, error) {
	if _, ok := list.Customs[name]; ok {
		cfg, err := LoadCustom(list.cfgDir, name)
		if err != nil {
			return nil, err
		}
		return cfg.Colors()
	}
	if _, ok := list.Builtins[name]; ok {
		cfg, err := LoadBuiltin(name)
		if err != nil {
			return nil, err
		}
		return cfg.Colors()
	}
	return nil, fmt.Errorf("not found")
}

func Load(cfgDir string, name string) (*Colors, error) {
	list, err := LoadList(cfgDir)
	if err != nil {
		return nil, err
	}
	return list.Load(name)
}

func Parse(s string) (*Colors, error) {
	cfg, err := ParseConfig([]byte(s))
	if err != nil {
		return nil, err
	}
	return cfg.Colors()
}
