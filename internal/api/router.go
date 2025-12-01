package api

import (
	"github.com/gin-gonic/gin"
	"rt-manage/internal/api/handler"
	"rt-manage/internal/config"
	"rt-manage/internal/database"
	"rt-manage/internal/middleware"
	"rt-manage/internal/repository"
	"rt-manage/internal/service"
	"rt-manage/internal/web"
	"rt-manage/pkg/logger"
)

// NewRouter 创建路由
func NewRouter() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(config.Get().Server.Mode)

	r := gin.New()

	// 使用中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", handler.Health)

	// 初始化仓库
	db := database.GetDB()
	rtRepo := repository.NewRTRepository(db)
	configRepo := repository.NewConfigRepository(db)

	// 初始化服务
	rtService := service.NewRTService(rtRepo, configRepo)
	configService := service.NewConfigService(configRepo, rtRepo)

	// 初始化处理器
	rtHandler := handler.NewRTHandler(rtService, configService)
	configHandler := handler.NewConfigHandler(configService)
	authHandler := handler.NewAuthHandler()
	publicAPIHandler := handler.NewPublicAPIHandler(rtService)

	// 对外公开API路由组（使用API Secret认证）
	// 从配置文件读取路由前缀，默认为 "/public-api"
	publicAPIPrefix := config.Get().Auth.PublicAPIPrefix
	if publicAPIPrefix == "" {
		publicAPIPrefix = "/public-api"
	}
	
	// 验证前缀，防止与内部管理 API 冲突
	if publicAPIPrefix == "/internalweb" || publicAPIPrefix == "/health" {
		logger.Error("public_api_prefix 配置错误：不能使用 /internalweb 或 /health，这些路径已被系统占用", "prefix", publicAPIPrefix)
		panic("public_api_prefix 配置错误：不能使用 /internalweb 或 /health")
	}
	
	publicAPI := r.Group(publicAPIPrefix)
	publicAPI.Use(middleware.APISecret())
	{
		publicAPI.GET("/health", handler.Health)                          // 健康检查
		publicAPI.POST("/refresh", publicAPIHandler.RefreshAndGetAT)      // 刷新RT并获取AT
		publicAPI.POST("/get-at", publicAPIHandler.GetAT)                 // 获取AT（不刷新）
	}
	
	logger.Info("对外API路由前缀", "prefix", publicAPIPrefix)

	// API路由组（内部管理使用）
	api := r.Group("/internalweb/v1")
	{
		// 认证路由（不需要JWT验证）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)  // 登录
		}

		// 需要JWT认证的路由（全部使用POST + JSON Body）
		authorized := api.Group("")
		authorized.Use(middleware.JWTAuth())
		{
			// 用户信息
			user := authorized.Group("/user")
			{
				user.POST("/info", authHandler.GetCurrentUser)  // 获取当前用户
			}

			// RT管理路由
			rts := authorized.Group("/rts")
			{
				rts.POST("/list", rtHandler.ListRTs)                // 获取列表
				rts.POST("/create", rtHandler.CreateRT)             // 创建RT
				rts.POST("/update", rtHandler.UpdateRT)             // 更新RT
				rts.POST("/delete", rtHandler.DeleteRT)             // 删除RT
				rts.POST("/batch-delete", rtHandler.BatchDeleteRTs) // 批量删除
			rts.POST("/batch-refresh", rtHandler.BatchRefreshRTs) // 批量刷新
			rts.POST("/batch-import", rtHandler.BatchImportRTs) // 批量导入
			rts.POST("/refresh", rtHandler.RefreshRT)           // 单个刷新
			rts.POST("/refresh-user-info", rtHandler.RefreshUserInfo)       // 刷新用户信息
			rts.POST("/refresh-account-info", rtHandler.RefreshAccountInfo) // 刷新账号信息
		}

			// 配置管理路由
			configs := authorized.Group("/configs")
			{
				configs.POST("/get-system", configHandler.GetSystemConfigs)       // 获取系统配置
				configs.POST("/save-system", configHandler.SaveSystemConfigs)     // 保存系统配置
				configs.POST("/get-proxy-list", configHandler.GetProxyList)       // 获取代理列表
				configs.POST("/get-clientid-list", configHandler.GetClientIDList) // 获取 Client ID 列表
			}
		}
	}

	// 设置静态文件服务（前端页面）
	if err := web.SetupStaticRoutes(r); err != nil {
		logger.Error("设置静态文件路由失败", "error", err)
	}

	return r
}

