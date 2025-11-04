package middleware

import (
	"strings"

	"rt-manage/internal/config"
	jwtutil "rt-manage/pkg/jwt"
	"rt-manage/pkg/logger"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件（支持JWT Token和API Secret两种方式）
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"success": false,
				"msg":     "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 检查格式: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(401, gin.H{
				"success": false,
				"msg":     "认证令牌格式错误，应为: Bearer <token>",
			})
			c.Abort()
			return
		}

		token := parts[1]
		cfg := config.Get()

		// 方式1: 检查是否为 API Secret（固定密码）
		if cfg.Auth.APISecret != "" && token == cfg.Auth.APISecret {
			logger.Debug("使用API Secret认证通过")
			// API Secret 认证成功，设置默认用户名
			c.Set("username", "api_user")
			c.Set("auth_type", "api_secret")
			c.Next()
			return
		}

		// 方式2: 尝试作为 JWT Token 解析
		claims, err := jwtutil.ParseToken(token, cfg.Auth.JWTSecret)
		if err != nil {
			logger.Debug("JWT解析失败，且不匹配API Secret", "error", err)
			c.JSON(401, gin.H{
				"success": false,
				"msg":     "认证失败：令牌无效或已过期",
			})
			c.Abort()
			return
		}

		// JWT 认证成功
		logger.Debug("使用JWT Token认证通过", "username", claims.Username)
		c.Set("username", claims.Username)
		c.Set("auth_type", "jwt")
		c.Next()
	}
}
