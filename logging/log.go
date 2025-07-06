package logging

import (
	"log"

	"github.com/gin-gonic/gin"
)

func DebugLog(msg string) {
	if gin.Mode() == gin.DebugMode {
		log.Println(msg)
	}
}
