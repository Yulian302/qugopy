package config

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
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

func TestRedisHealth(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() error = %v", err)
		return
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.REDIS.HOST, cfg.REDIS.PORT),
	})

	defer rdb.Close()

	if _, err := rdb.Ping().Result(); err != nil {
		t.Fatalf("Redis ping failed: %v", err)
	}
}
