package handlers

import (
	"net/http"

	"github.com/Yulian302/qugopy/internal/tasks"
	"github.com/Yulian302/qugopy/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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
		err := tasks.EnqueueTask(task, rdb)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":   "Task enqueued!",
			"priority": task.Priority,
			"type":     task.Type,
		})
	}

}
