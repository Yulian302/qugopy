package handlers

import (
	"bytes"
	"encoding/json"
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

func newTestRouter(rdb *redis.Client) *gin.Engine {
	r := gin.New()
	r.POST("/tasks", TaskEnqueueHandler(rdb))
	return r
}

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

	_ = rdb.Del("go_queue").Err()
	t.Cleanup(func() {
		_ = rdb.Del("go_queue").Err()
	})
	r := newTestRouter(rdb)

	tests := []struct {
		queueType  string
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			queueType:  "go_queue",
			name:       "valid task",
			body:       `{"type": "download_file", "payload": "test", "priority": 10}`,
			wantStatus: 201,
			wantBody:   "Task enqueued",
		},
		{
			queueType:  "go_queue",
			name:       "missing type field",
			body:       `{"payload": "download_file", "priority": 10}`,
			wantStatus: 400,
			wantBody:   "Invalid request payload",
		},
		{
			queueType:  "python_queue",
			name:       "invalid priority",
			body:       `{"payload": "payload", "priority": -1}`,
			wantStatus: 400,
			wantBody:   "Invalid request payload",
		},
		{
			queueType:  "python_queue",
			name:       "invalid json",
			body:       `{"type": "download_file",`,
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
				res, err := rdb.ZPopMin(tt.queueType, 1).Result()
				assert.NoError(t, err)
				var task queue.IntTask
				if err := json.Unmarshal([]byte(res[0].Member.(string)), &task); err != nil {
					log.Fatal(err)
				}
				assert.Equal(t, task.Task.Type, "download_file")
			}
		})
	}

}

func TestEnqueueHandlerLocal_EdgeCases(t *testing.T) {

	config.AppConfig.MODE = "local"
	r := newTestRouter(rdb)

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "valid task",
			body:       `{"type": "download_file", "payload": "test", "priority": 10}`,
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
			body:       `{"type": "download_file",`,
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
				head, exists := queue.GoLocalQueue.PQ.Pop()
				assert.Equal(t, true, exists)
				assert.Equal(t, head.Task.Type, "download_file")
			}
		})
	}

}
