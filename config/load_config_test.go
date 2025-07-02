package config

import (
	"testing"
)

func (cfg *AppConfig) isValid() bool {
	return cfg.HOST != "" &&
		cfg.PORT != "" &&
		cfg.REDIS.HOST != "" &&
		cfg.REDIS.PORT != ""
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "valid config",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cfg.isValid() {
				t.Error("Config is invalid - required fields are empty")
			}
		})
	}
}
