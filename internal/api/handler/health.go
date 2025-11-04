package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health 健康检查
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "service is healthy",
	})
}

// Ping 测试接口
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

