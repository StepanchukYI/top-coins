package config

import (
	"encoding/json"
	"os"
)

// Application holds application configuration values
type Config struct {
	BindAddr    string
	DatabaseDSN string
	LogLevel    string
}

func NewConfig(file string) (cfg *Config, err error) {
	cfg = &Config{}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
