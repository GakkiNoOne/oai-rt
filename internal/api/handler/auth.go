package handler

import (
	"rt-manage/internal/config"
	jwtutil "rt-manage/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	Username string `json:"username"`
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	// 从配置读取用户名和密码
	cfg := config.Get()
	if req.Username != cfg.Auth.Username || req.Password != cfg.Auth.Password {
		c.JSON(401, gin.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 生成JWT token
	token, err := jwtutil.GenerateToken(req.Username, cfg.Auth.JWTSecret, cfg.Auth.JWTExpireHours)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "生成令牌失败",
		})
		return
	}

	c.JSON(200, LoginResponse{
		Token: token,
		Username: req.Username,
	})
}

// GetCurrentUser 获取当前用户信息 - POST /api/user/info
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"msg":     "未认证",
		})
		return
	}

	c.JSON(200, gin.H{
		"success":  true,
		"msg":      "获取成功",
		"username": username,
	})
}

