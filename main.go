package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/Yulian302/qugopy/config"
	"github.com/Yulian302/qugopy/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var rdb *redis.Client
var ctx = context.Background()

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "ok",
	})
}
func TaskEnqueueHandler(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}
	var mode models.Mode
	mode = models.Mode(c.Query("mode"))
	if mode == "" {
		log.Printf("MODE is empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "MODE is absent in query parameters. Please provide MODE",
		})
		return
	}
	if !mode.IsValid() {
		log.Printf("MODE is not valid")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "MODE is not valid. Check docs for more information",
		})
		return
	}

	internalTask := &models.IntTask{
		Task: task,
		ID:   uuid.New().String(),
	}
	userTaskJson, err := json.Marshal(internalTask)
	if err != nil {
		log.Printf("Failed to marshal task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if mode == "redis" {
		intCmd := rdb.LPush("task_queue", userTaskJson)
		if err := intCmd.Err(); err != nil {
			log.Printf("Redis LPush failed: %v", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Task queue unavailable",
				"details": err.Error(),
			})
			return
		}
	} else {
		// enqueue locally
		fmt.Print("enqueue locally")
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":   "Task enqueued!",
		"priority": task.Priority,
		"type":     task.Type,
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
