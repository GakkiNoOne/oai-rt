package middleware

import (
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"rt-manage/pkg/logger"
)

// 静态文件扩展名列表
var staticFileExtensions = map[string]bool{
	".js":   true,
	".css":  true,
	".html": true,
	".ico":  true,
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".svg":  true,
	".woff": true,
	".woff2": true,
	".ttf":  true,
	".eot":  true,
	".map":  true,
}

// isStaticFile 判断是否为静态文件请求
func isStaticFile(path string) bool {
	ext := filepath.Ext(path)
	return staticFileExtensions[ext]
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		// 跳过静态文件的请求日志
		if isStaticFile(path) {
			return
		}

		cost := time.Since(start)
		
		logger.Info("HTTP请求",
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"status", c.Writer.Status(),
			"cost", cost,
			"ip", c.ClientIP(),
			"user-agent", c.Request.UserAgent(),
		)
	}
}

