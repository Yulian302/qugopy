package api

import (
	"github.com/Yulian302/qugopy/internal/api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func NewRouter(rdb *redis.Client) *gin.Engine {
	router := gin.Default()

	// router.Use(
	// 	gin.Logger(),
	// 	gin.Recovery(),
	// )

	router.GET("/test", handlers.HealthCheckHandler)
	router.POST("/tasks", handlers.TaskEnqueueHandler(rdb))

	return router
}
