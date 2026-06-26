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
	Shared string `yaml:"shared"`

	EscTimeout int `yaml:"esc-timeout"`
}

func DefaultConfig() *Config {
	return &Config{
		AutoIndent: true,
		TabStop:    4,

		Colors: "standard",
		Silent: false,
		Shared: "xyz",

		EscTimeout: 100,
	}
}

func ConfigPath(dir string) string {
	return filepath.Join(dir, ConfigFilename)
}

func LoadConfig(dir string) (*Config, error) {
	path := ConfigPath(dir)
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

func SaveConfig(dir string, cfg *Config) error {
	path := ConfigPath(dir)
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
func PrepareConfig(dir string) (*Config, error) {
	path := ConfigPath(dir)
	_, err := os.Stat(path)
	if err == nil {
		cfg, err := LoadConfig(dir)
		if err != nil {
			return DefaultConfig(), err
		}
		return cfg, nil
	} else {
		cfg := DefaultConfig()
		err := SaveConfig(dir, cfg)
		return cfg, err
	}
}
