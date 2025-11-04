package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rt-manage/pkg/logger"
)

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("发生panic", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "内部服务器错误",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

