package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rt-manage/internal/service"
	"rt-manage/pkg/logger"
)

// ConfigHandler 配置处理器
type ConfigHandler struct {
	configService service.ConfigService
}

// NewConfigHandler 创建配置处理器实例
func NewConfigHandler(configService service.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
	}
}

// GetSystemConfigs 获取系统配置 - POST /api/configs/get-system
func (h *ConfigHandler) GetSystemConfigs(c *gin.Context) {
	logger.Info("获取系统配置 - 请求")

	dbConfigs, envConfigs, err := h.configService.GetSystemConfigs()
	if err != nil {
		logger.Error("获取系统配置失败", "error", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Msg:     "获取配置失败: " + err.Error(),
		})
		return
	}

	logger.Info("获取系统配置成功", "db_configs_count", len(dbConfigs), "env_configs_count", len(envConfigs))

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "获取成功",
		Data: gin.H{
			"configs":     dbConfigs,
			"env_configs": envConfigs,
		},
	})
}

// SaveSystemConfigs 保存系统配置 - POST /api/configs/save-system
func (h *ConfigHandler) SaveSystemConfigs(c *gin.Context) {
	var req struct {
		Configs map[string]string `json:"configs" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("保存系统配置 - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     "参数错误: " + err.Error(),
		})
		return
	}

	logger.Info("保存系统配置 - 请求", "configs", req.Configs)

	if err := h.configService.SaveSystemConfigs(req.Configs); err != nil {
		logger.Error("保存系统配置失败", "configs", req.Configs, "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     "保存失败: " + err.Error(),
		})
		return
	}

	logger.Info("保存系统配置成功", "configs", req.Configs)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "保存成功",
	})
}

// GetProxyList 获取代理列表 - POST /api/configs/get-proxy-list
func (h *ConfigHandler) GetProxyList(c *gin.Context) {
	logger.Info("获取代理列表 - 请求")

	proxyList, err := h.configService.GetProxyList()
	if err != nil {
		logger.Error("获取代理列表失败", "error", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Msg:     "获取代理列表失败: " + err.Error(),
		})
		return
	}

	logger.Info("获取代理列表成功", "proxy_count", len(proxyList), "proxies", proxyList)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "获取成功",
		Data:    proxyList,
	})
}

// GetClientIDList 获取 Client ID 列表 - POST /api/configs/get-clientid-list
func (h *ConfigHandler) GetClientIDList(c *gin.Context) {
	logger.Info("获取 Client ID 列表 - 请求")

	clientIDList, err := h.configService.GetClientIdList()
	if err != nil {
		logger.Error("获取 Client ID 列表失败", "error", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Msg:     "获取 Client ID 列表失败: " + err.Error(),
		})
		return
	}

	logger.Info("获取 Client ID 列表成功", "clientid_count", len(clientIDList), "clientids", clientIDList)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "获取成功",
		Data:    clientIDList,
	})
}
