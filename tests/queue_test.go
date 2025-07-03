package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/Yulian302/qugopy/config"
	"github.com/go-redis/redis"
)

var r *redis.Client

func cleanup() {
	if r != nil {
		_ = r.Close()
	}
}

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("LoadConfig() error = %v\n", err)
		os.Exit(1)
	}

	r = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.REDIS.HOST, cfg.REDIS.PORT),
	})

	code := m.Run()
	cleanup()
	os.Exit(code)
}

func TestRedisHealth(t *testing.T) {

	if _, err := r.Ping().Result(); err != nil {
		t.Fatalf("Redis ping failed: %v", err)
	}
}

func TestCreateInitQueue(t *testing.T) {

	_ = r.Del("init_list").Err()
	defer r.Del("init_list")

	values := []interface{}{1, 2, 3}
	if err := r.LPush("init_list", values...).Err(); err != nil {
		t.Fatal("Could not create list: init_list")
	}

	length, err := r.LLen("init_list").Result()
	if err != nil {
		t.Fatalf("Could not get list length: %v", err)
	}
	if length != 3 {
		t.Fatalf("Expected list length 3, got %d", length)
	}
}
