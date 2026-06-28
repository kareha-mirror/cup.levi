package color

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

const colorsDir = "colors"
const dotYAML = ".yaml"

//go:embed custom.yaml
var CustomSchemeConfig string

//go:embed colors/*.yaml
var BuiltinsFS embed.FS

type Scheme struct {
	Buffer  termi.ColorPair
	Status  termi.ColorPair
	Current termi.ColorPair
	Border  termi.ColorPair
}

type SchemeConfig struct {
	Buffer  string `yaml:"buffer"`
	Status  string `yaml:"status"`
	Current string `yaml:"current"`
	Border  string `yaml:"border"`
}

func (cfg *SchemeConfig) Scheme() (*Scheme, error) {
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

	return &Scheme{
		Buffer:  buffer,
		Status:  status,
		Current: current,
		Border:  border,
	}, nil
}

type SchemeList struct {
	cfgDir string

	Customs  map[string]struct{}
	Builtins map[string]struct{}

	Names []string
}

func ListCustomSchemes(cfgDir string) ([]string, error) {
	dir := filepath.Join(cfgDir, colorsDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, dotYAML) {
			continue
		}
		name = name[:len(name)-len(dotYAML)]
		names = append(names, name)
	}
	return names, nil
}

func ListBuiltinSchemes() ([]string, error) {
	entries, err := fs.ReadDir(BuiltinsFS, colorsDir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, dotYAML) {
			continue
		}
		name = name[:len(name)-len(dotYAML)]
		names = append(names, name)
	}
	return names, nil
}

func LoadSchemeList(cfgDir string) (*SchemeList, error) {
	// ensure directory and example file exist
	customDir := filepath.Join(cfgDir, colorsDir)
	os.Mkdir(customDir, 0777)
	customPath := filepath.Join(customDir, "custom.yaml")
	_, err := os.Stat(customPath)
	if err != nil { // file not exists
		os.WriteFile(customPath, []byte(CustomSchemeConfig), 0666)
	}

	// load list of names
	cList, err := ListCustomSchemes(cfgDir)
	if err != nil {
		return nil, err
	}
	bList, err := ListBuiltinSchemes()
	if err != nil {
		return nil, err
	}

	// create indices
	customs := map[string]struct{}{}
	for _, name := range cList {
		customs[name] = struct{}{}
	}
	builtins := map[string]struct{}{}
	for _, name := range bList {
		builtins[name] = struct{}{}
	}

	// merge and sort names
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

	return &SchemeList{cfgDir, customs, builtins, total}, nil
}

func ParseSchemeConfig(b []byte) (*SchemeConfig, error) {
	var cfg SchemeConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadCustomScheme(cfgDir string, name string) (*SchemeConfig, error) {
	path := filepath.Join(cfgDir, colorsDir, name+dotYAML)
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseSchemeConfig(b)
}

func LoadBuiltinScheme(name string) (*SchemeConfig, error) {
	path := filepath.Join(colorsDir, name+dotYAML)
	b, err := fs.ReadFile(BuiltinsFS, path)
	if err != nil {
		return nil, err
	}
	return ParseSchemeConfig(b)
}

func (list *SchemeList) Load(name string) (*Scheme, error) {
	// first, search customs
	if _, ok := list.Customs[name]; ok {
		cfg, err := LoadCustomScheme(list.cfgDir, name)
		if err != nil {
			return nil, err
		}
		return cfg.Scheme()
	}
	// second, search builtins
	if _, ok := list.Builtins[name]; ok {
		cfg, err := LoadBuiltinScheme(name)
		if err != nil {
			return nil, err
		}
		return cfg.Scheme()
	}
	return nil, fmt.Errorf("not found")
}

// for convenience on startup
func LoadScheme(cfgDir string, name string) (*Scheme, error) {
	list, err := LoadSchemeList(cfgDir)
	if err != nil {
		return nil, err
	}
	return list.Load(name)
}

// for loading directly from buffer
func ParseScheme(s string) (*Scheme, error) {
	cfg, err := ParseSchemeConfig([]byte(s))
	if err != nil {
		return nil, err
	}
	return cfg.Scheme()
}
