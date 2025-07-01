package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Yulian302/qugopy/config"
	"github.com/gin-gonic/gin"
)



func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"health": "ok",
	})
}

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()


	r.GET("/test", healthCheckHandler)

	if err := r.Run(fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT)); err != nil {
		log.Fatal("Failed to start server: %w", err)
	}
}
