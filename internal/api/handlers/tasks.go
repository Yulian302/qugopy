package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Yulian302/qugopy/config"
	"github.com/Yulian302/qugopy/internal/queue"
	"github.com/Yulian302/qugopy/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

func TaskEnqueueHandler(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task models.Task

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request payload",
				"details": err.Error(),
			})
			return
		}
		mode := config.AppConfig.MODE

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

		// push to redis
		if mode == "redis" {
			intCmd := rdb.ZAdd("task_queue", redis.Z{
				Score:  float64(task.Priority),
				Member: userTaskJson,
			})
			if err := intCmd.Err(); err != nil {
				log.Printf("Redis ZAdd failed: %v", err)
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error":   "Task queue unavailable",
					"details": err.Error(),
				})
				return
			}
		} else {
			// enqueue locally
			queue.DefaultLocalQueue.Lock.Lock()
			defer queue.DefaultLocalQueue.Lock.Unlock()
			queue.DefaultLocalQueue.PQ.Push(*internalTask)
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":   "Task enqueued!",
			"priority": task.Priority,
			"type":     task.Type,
		})
	}

}
