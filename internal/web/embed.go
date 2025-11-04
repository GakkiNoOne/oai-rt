package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed all:dist
var distFS embed.FS

// GetDistFS 获取嵌入的前端文件系统
func GetDistFS() fs.FS {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}

// SetupStaticRoutes 设置静态文件路由
func SetupStaticRoutes(r *gin.Engine) error {
	distFileSystem := GetDistFS()

	// 处理静态文件
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		
		// API 路由不处理
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"msg":     "API路径不存在",
			})
			return
		}

		// 如果是根路径或没有文件扩展名，返回 index.html
		if path == "/" || filepath.Ext(path) == "" {
			data, err := fs.ReadFile(distFileSystem, "index.html")
			if err != nil {
				c.String(http.StatusInternalServerError, "读取index.html失败")
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
			return
		}

		// 尝试读取请求的文件
		filePath := strings.TrimPrefix(path, "/")
		data, err := fs.ReadFile(distFileSystem, filePath)
		if err != nil {
			// 文件不存在，返回 index.html（支持前端路由）
			data, err = fs.ReadFile(distFileSystem, "index.html")
			if err != nil {
				c.String(http.StatusInternalServerError, "读取index.html失败")
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
			return
		}

		// 根据文件扩展名设置Content-Type
		ext := filepath.Ext(filePath)
		contentType := getContentType(ext)
		c.Data(http.StatusOK, contentType, data)
	})

	return nil
}

// getContentType 根据文件扩展名返回Content-Type
func getContentType(ext string) string {
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	default:
		return "application/octet-stream"
	}
}

