package config

import (
	"testing"
)

func (cfg *RootConfig) isValid() bool {
	return cfg.HOST != "" &&
		cfg.PORT != "" &&
		cfg.REDIS.HOST != "" &&
		cfg.REDIS.PORT != ""
}

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
	}
	if !cfg.isValid() {
		t.Error("Config is invalid - required fields are empty")
	}
}
