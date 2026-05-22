package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var cfg *Config

func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	paths := []string{
		filepath.Join(os.Getenv("HOME"), ".config", "veego", "config.yaml"),
		filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "veego", "config.yaml"),
		"config/config.yaml",
	}

	var data []byte
	var err error
	for _, path := range paths {
		if path == "" {
			continue
		}
		data, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	cfg = &c
	return cfg, nil
}

func Get() *Config {
	return cfg
}
