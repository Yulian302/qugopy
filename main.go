package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Yulian302/qugopy/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var rdb *redis.Client
var ctx = context.Background()

type Task struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "ok",
	})
}
func TaskEnqueueHandler(c *gin.Context) {
	var task Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invdalid request payload",
			"details": err.Error(),
		})
		return
	}

	taskJson, err := json.Marshal(task)
	if err != nil {
		log.Printf("Failed to marshal task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	intCmd := rdb.LPush("task_queue", taskJson)
	if err := intCmd.Err(); err != nil {
		log.Printf("Redis LPush failed: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Task queue unavailable",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "Task enqueued!",
		"id":     task.ID,
	})

}

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.REDIS.HOST, cfg.REDIS.PORT),
	})

	r.GET("/test", healthCheckHandler)
	r.POST("/task", TaskEnqueueHandler)

	if err := r.Run(fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT)); err != nil {
		log.Fatal("Failed to start server: %w", err)
	}
}
