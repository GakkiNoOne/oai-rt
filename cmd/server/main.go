package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"rt-manage/internal/api"
	"rt-manage/internal/config"
	"rt-manage/internal/database"
	"rt-manage/internal/repository"
	"rt-manage/internal/scheduler"
	"rt-manage/internal/service"
	"rt-manage/pkg/logger"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 初始化日志
	if err := logger.Init(); err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}
	defer logger.Sync()

	// 获取配置并打印（用于调试）
	cfg := config.Get()
	logger.Info("=== 当前配置信息 ===")
	logger.Info("服务器配置", 
		"host", cfg.Server.Host, 
		"port", cfg.Server.Port, 
		"mode", cfg.Server.Mode)
	logger.Info("数据库配置", 
		"type", cfg.Database.Type, 
		"host", cfg.Database.Host, 
		"port", cfg.Database.Port, 
		"database", cfg.Database.Database,
		"user", cfg.Database.User,
		"table_prefix", cfg.Database.TablePrefix)
	logger.Info("OpenAI配置", 
		"client_id", cfg.OpenAI.ClientID, 
		"refresh_interval", cfg.OpenAI.RefreshInterval,
		"schedule_enabled", cfg.OpenAI.ScheduleEnabled,
		"proxy", cfg.OpenAI.Proxy)
	logger.Info("认证配置", 
		"username", cfg.Auth.Username, 
		"jwt_expire_hours", cfg.Auth.JWTExpireHours,
		"api_secret_length", len(cfg.Auth.APISecret))
	logger.Info("==================")

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		logger.Fatal("数据库初始化失败", "error", err)
	}
	defer database.Close()

	logger.Info("数据库初始化成功", "type", cfg.Database.Type)

	// 初始化调度器管理器
	db := database.GetDB()
	rtRepo := repository.NewRTRepository(db)
	configRepo := repository.NewConfigRepository(db)
	rtService := service.NewRTService(rtRepo, configRepo)
	configService := service.NewConfigService(configRepo, rtRepo)
	
	// 初始化全局调度器管理器
	scheduler.InitManager(rtService, cfg.OpenAI.RefreshInterval)
	defer scheduler.GetManager().Stop()
	
	// 从数据库读取配置，决定是否启动调度器
	dbConfigs, _, err := configService.GetSystemConfigs()
	if err == nil {
		autoRefreshEnabled := dbConfigs["auto_refresh_enabled"]
		autoRefreshInterval := dbConfigs["auto_refresh_interval"]
		
		// 优先使用数据库配置，如果数据库没有配置则使用config.yaml
		if autoRefreshEnabled == "true" {
			scheduler.GetManager().UpdateFromConfig(autoRefreshEnabled, autoRefreshInterval)
			logger.Info("根据数据库配置启动调度器", "enabled", true, "interval", autoRefreshInterval)
		} else if cfg.OpenAI.ScheduleEnabled {
			// 如果数据库没有明确禁用，且config.yaml启用了，则启动
			scheduler.GetManager().Start(cfg.OpenAI.RefreshInterval)
			logger.Info("根据配置文件启动调度器", "interval_hours", cfg.OpenAI.RefreshInterval)
		}
	} else {
		logger.Error("读取数据库配置失败，使用config.yaml配置", "error", err)
		if cfg.OpenAI.ScheduleEnabled {
			scheduler.GetManager().Start(cfg.OpenAI.RefreshInterval)
			logger.Info("根据配置文件启动调度器", "interval_hours", cfg.OpenAI.RefreshInterval)
		}
	}

	// 创建路由
	router := api.NewRouter()

	// 启动服务
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("服务启动", "address", addr)

	// 优雅关闭
	go func() {
		if err := router.Run(addr); err != nil {
			logger.Fatal("服务启动失败", "error", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.Info("服务正在关闭...")
}

