package editor

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const ConfigFilename = "editor.yaml"

type Config struct {
	AutoIndent bool `yaml:"auto-indent"`
	TabStop    int  `yaml:"tab-stop"`

	Colors string `yaml:"colors"`
	Silent bool   `yaml:"silent"`
	CRLF   bool   `yaml:"crlf"`
	Depth  int    `yaml:"depth"`
	Shared string `yaml:"shared"`

	EscapeTimeout int `yaml:"escape-timeout"`
}

func ConfigPath(cfgDir string) string {
	return filepath.Join(cfgDir, ConfigFilename)
}

func LoadConfig(cfgDir string) (*Config, error) {
	path := ConfigPath(cfgDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(cfgDir string, cfg *Config) error {
	path := ConfigPath(cfgDir)
	data, err := yaml.Marshal(cfg)
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

// return default on error
func PrepareConfig(cfgDir string) (*Config, error) {
	path := ConfigPath(cfgDir)
	_, err := os.Stat(path)
	if err == nil {
		cfg, err := LoadConfig(cfgDir)
		if err != nil {
			return DefaultConfig(), err
		}
		return cfg, nil
	} else {
		cfg := DefaultConfig()
		err := SaveConfig(cfgDir, cfg)
		return cfg, err
	}
}
