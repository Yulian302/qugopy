package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Yulian302/qugopy/config"
	"github.com/Yulian302/qugopy/internal/queue"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"github.com/stretchr/testify/assert"
)

var (
	cfg *config.RootConfig
	err error
	r   *gin.Engine
	rdb *redis.Client
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal("config is not set up. check config file")
		os.Exit(1)
	}
	cfg.MODE = "redis"
	r = gin.New()
	code := m.Run()
	os.Exit(code)
}

func TestHealth(t *testing.T) {
	path := "/test"
	r.GET(path, HealthCheckHandler)

	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"health"`)
}
func TestEnqueueHandlerRedis_EdgeCases(t *testing.T) {
	config.AppConfig.MODE = "redis"
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.REDIS.HOST, cfg.REDIS.PORT),
	})

	_ = rdb.Del("task_queue").Err()
	t.Cleanup(func() {
		_ = rdb.Del("task_queue").Err()
	})

	r.POST("/tasks", TaskEnqueueHandler(rdb))

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "valid task",
			body:       `{"type": "email", "payload": "test", "priority": 10}`,
			wantStatus: 201,
			wantBody:   "Task enqueued",
		},
		{
			name:       "missing type field",
			body:       `{"payload": "test", "priority": 10}`,
			wantStatus: 400,
			wantBody:   "Invalid request payload",
		},
		{
			name:       "invalid priority",
			body:       `{"payload": "test", "priority": -1}`,
			wantStatus: 400,
			wantBody:   "Invalid request payload",
		},
		{
			name:       "invalid json",
			body:       `{"type": "email",`,
			wantStatus: 400,
			wantBody:   "Invalid request payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/tasks"
			req, _ := http.NewRequest("POST", path, bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != "" {
				assert.Contains(t, w.Body.String(), tt.wantBody)
			}
			if tt.wantStatus == 201 {
				tasks, err := rdb.LRange("task_queue", 0, -1).Result()
				assert.NoError(t, err)
				assert.NotEmpty(t, tasks)
				assert.Contains(t, tasks[0], `"type":"email"`)
			}
		})
	}

}

func TestEnqueueHandlerLocal_EdgeCases(t *testing.T) {

	config.AppConfig.MODE = "local"
	r.POST("/tasks", TaskEnqueueHandler(rdb))

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "valid task",
			body:       `{"type": "email", "payload": "test", "priority": 10}`,
			wantStatus: 201,
			wantBody:   "Task enqueued",
		},
		{
			name:       "missing type field",
			body:       `{"payload": "test", "priority": 10}`,
			wantStatus: 400,
			wantBody:   "Invalid request payload",
		},
		{
			name:       "invalid priority",
			body:       `{"payload": "test", "priority": -1}`,
			wantStatus: 400,
			wantBody:   "Invalid request payload",
		},
		{
			name:       "invalid json",
			body:       `{"type": "email",`,
			wantStatus: 400,
			wantBody:   "Invalid request payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/tasks"
			req, _ := http.NewRequest("POST", path, bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != "" {
				assert.Contains(t, w.Body.String(), tt.wantBody)
			}
			if tt.wantStatus == 201 {
				head, exists := queue.DefaultLocalQueue.PQ.Pop()
				assert.Equal(t, true, exists)
				assert.Equal(t, head.Task.Type, "email")
			}
		})
	}

}
