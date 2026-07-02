package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const Filename = "editor.yaml"

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

func Path(cfgDir string) string {
	return filepath.Join(cfgDir, Filename)
}

func Load(cfgDir string) (*Config, error) {
	path := Path(cfgDir)
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

func Save(cfgDir string, cfg *Config) error {
	path := Path(cfgDir)
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
func Prepare(cfgDir string) (*Config, error) {
	path := Path(cfgDir)
	_, err := os.Stat(path)
	if err == nil {
		cfg, err := Load(cfgDir)
		if err != nil {
			return Default(), err
		}
		return cfg, nil
	} else {
		cfg := Default()
		err := Save(cfgDir, cfg)
		return cfg, err
	}
}
