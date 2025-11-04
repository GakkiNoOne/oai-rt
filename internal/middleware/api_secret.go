package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rt-manage/internal/config"
	"rt-manage/pkg/logger"
)

// APISecret API密钥认证中间件
func APISecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiSecret := c.GetHeader("X-API-Secret")
		
		// 获取配置的 API Secret
		cfg := config.Get()
		expectedSecret := cfg.Auth.APISecret
		
		// 如果配置为空，使用默认值
		if expectedSecret == "" {
			expectedSecret = "your-api-secret-here"
		}
		
		if apiSecret == "" {
			logger.Warn("API请求缺少密钥", "path", c.Request.URL.Path, "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"msg":     "缺少API密钥",
			})
			c.Abort()
			return
		}
		
		if apiSecret != expectedSecret {
			logger.Warn("API密钥验证失败", "path", c.Request.URL.Path, "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"msg":     "API密钥无效",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

