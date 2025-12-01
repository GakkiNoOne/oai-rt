package handler

import (
	"fmt"
	"net/http"

	"rt-manage/internal/model"
	"rt-manage/internal/service"
	"rt-manage/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RTHandler RT 处理器
type RTHandler struct {
	rtService     service.RTService
	configService service.ConfigService
}

// NewRTHandler 创建 RT 处理器实例
func NewRTHandler(rtService service.RTService, configService service.ConfigService) *RTHandler {
	return &RTHandler{
		rtService:     rtService,
		configService: configService,
	}
}

// APIResponse 统一响应格式
type APIResponse struct {
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

// ListRTs 获取RT列表 - POST /api/rts/list
func (h *RTHandler) ListRTs(c *gin.Context) {
	var req struct {
		Page       int    `json:"page"`
		PageSize   int    `json:"page_size"`
		BizId      string `json:"biz_id"`
		Tag        string `json:"tag"`
		Email      string `json:"email"`
		Type       string `json:"type"`
		Enabled    *bool  `json:"enabled"`
		CreateDate string `json:"create_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("获取RT列表 - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     "参数错误: " + err.Error(),
		})
		return
	}

	// 默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	logger.Info("获取RT列表 - 请求", "page", req.Page, "page_size", req.PageSize, "biz_id", req.BizId, "tag", req.Tag, "email", req.Email, "type", req.Type, "enabled", req.Enabled)

	rts, total, err := h.rtService.List(req.Page, req.PageSize, req.BizId, req.Tag, req.Email, req.Type, req.Enabled, req.CreateDate)
	if err != nil {
		logger.Error("获取RT列表失败", "error", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Msg:     "获取列表失败: " + err.Error(),
		})
		return
	}

	logger.Info("获取RT列表成功", "total", total, "returned", len(rts))

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "获取成功",
		Data: gin.H{
			"items":     rts,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// CreateRT 创建RT - POST /api/rts/create
func (h *RTHandler) CreateRT(c *gin.Context) {
	var req struct {
		BizId    string `json:"biz_id"`
		RTToken  string `json:"rt_token" binding:"required"`
		Proxy    string `json:"proxy"`
		ClientID string `json:"client_id"`
		Tag      string `json:"tag"`
		Enabled  bool   `json:"enabled"`
		Memo     string `json:"memo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("创建RT - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     "参数错误: " + err.Error(),
		})
		return
	}

	// 记录请求参数（完整RT Token）
	logger.Info("创建RT - 请求", "biz_id", req.BizId, "rt_token", req.RTToken, "proxy", req.Proxy, "client_id", req.ClientID, "tag", req.Tag, "enabled", req.Enabled)

	rt := &model.RT{
		BizId:    req.BizId,
		Rt:       req.RTToken,
		Proxy:    req.Proxy,
		ClientID: req.ClientID,
		Tag:      req.Tag,
		Enabled:  req.Enabled,
		Memo:     req.Memo,
	}

	if err := h.rtService.Create(rt); err != nil {
		logger.Error("创建RT失败", "biz_id", req.BizId, "rt_token", req.RTToken, "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}

	logger.Info("创建RT成功", "id", rt.ID, "biz_id", rt.BizId, "rt_token", rt.Rt)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "创建成功",
		Data:    rt,
	})
}

// UpdateRT 更新RT - POST /api/rts/update
func (h *RTHandler) UpdateRT(c *gin.Context) {
	var req struct {
		ID      int64                  `json:"id" binding:"required"`
		Updates map[string]interface{} `json:"updates" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("更新RT - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     "参数错误: " + err.Error(),
		})
		return
	}

	logger.Info("更新RT - 请求", "id", req.ID, "updates", req.Updates)

	rt, err := h.rtService.Update(req.ID, req.Updates)
	if err != nil {
		logger.Error("更新RT失败", "id", req.ID, "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}

	logger.Info("更新RT成功", "id", req.ID, "name", rt.BizId)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "更新成功",
		Data:    rt,
	})
}

// DeleteRT 删除RT - POST /api/rts/delete
func (h *RTHandler) DeleteRT(c *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("删除RT - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     "参数错误: " + err.Error(),
		})
		return
	}

	logger.Info("删除RT - 请求", "id", req.ID)

	if err := h.rtService.Delete(req.ID); err != nil {
		logger.Error("删除RT失败", "id", req.ID, "error", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Msg:     "删除失败: " + err.Error(),
		})
		return
	}

	logger.Info("删除RT成功", "id", req.ID)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "删除成功",
	})
}

// BatchDeleteRTs 批量删除RT - POST /api/rts/batch-delete
func (h *RTHandler) BatchDeleteRTs(c *gin.Context) {
	var req struct {
		IDs []int64 `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("批量删除RT - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     "参数错误: " + err.Error(),
		})
		return
	}

	logger.Info("批量删除RT - 请求", "ids", req.IDs, "count", len(req.IDs))

	successCount, failCount, err := h.rtService.BatchDelete(req.IDs)

	logger.Info("批量删除RT完成", "success_count", successCount, "fail_count", failCount, "error", err)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "批量删除完成",
		Data: gin.H{
			"success_count": successCount,
			"fail_count":    failCount,
			"error":         err,
		},
	})
}

// RefreshRT 刷新单个RT - POST /api/rts/refresh
func (h *RTHandler) RefreshRT(c *gin.Context) {
	var req struct {
		ID                 int64 `json:"id" binding:"required"`
		RefreshUserInfo    bool  `json:"refresh_user_info"`
		RefreshAccountInfo bool  `json:"refresh_account_info"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("刷新RT - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     "参数错误: " + err.Error(),
		})
		return
	}

	logger.Info("刷新RT - 请求", "id", req.ID, "refresh_user_info", req.RefreshUserInfo, "refresh_account_info", req.RefreshAccountInfo)

	rt, err := h.rtService.Refresh(req.ID, req.RefreshUserInfo, req.RefreshAccountInfo)
	if err != nil {
		logger.Error("刷新RT失败", "id", req.ID, "error", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Msg:     "刷新失败: " + err.Error(),
		})
		return
	}

	logger.Info("刷新RT成功", "id", rt.ID, "name", rt.BizId, "rt", rt.Rt, "at", rt.At, "email", rt.Email, "user_name", rt.UserName, "type", rt.Type)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "刷新成功",
		Data:    rt,
	})
}

// BatchRefreshRTs 批量刷新RT - POST /api/rts/batch-refresh
func (h *RTHandler) BatchRefreshRTs(c *gin.Context) {
	var req struct {
		IDs []int64 `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("批量刷新RT - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     "参数错误: " + err.Error(),
		})
		return
	}

	logger.Info("批量刷新RT - 请求", "ids", req.IDs, "count", len(req.IDs))

	successCount, failCount, results, err := h.rtService.BatchRefresh(req.IDs)
	if err != nil {
		logger.Error("批量刷新RT失败", "ids", req.IDs, "error", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Msg:     "批量刷新失败: " + err.Error(),
		})
		return
	}

	logger.Info("批量刷新RT完成", "total_count", len(req.IDs), "success_count", successCount, "fail_count", failCount)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "批量刷新完成",
		Data: gin.H{
			"total_count":   len(req.IDs),
			"success_count": successCount,
			"fail_count":    failCount,
			"results":       results,
		},
	})
}

// BatchImportRTs 批量导入RT - POST /api/rts/batch-import
func (h *RTHandler) BatchImportRTs(c *gin.Context) {
	var req struct {
		BatchName string   `json:"batch_name"`
		Tag       string   `json:"tag"`
		Proxy     string   `json:"proxy"`
		ClientID  string   `json:"client_id"`
		RTTokens  []string `json:"rt_tokens" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("批量导入RT - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     "参数错误: " + err.Error(),
		})
		return
	}

	// 记录请求参数（RT Tokens脱敏）
	tokenPreviews := make([]string, 0, len(req.RTTokens))
	for _, token := range req.RTTokens {
		if len(token) > 20 {
			tokenPreviews = append(tokenPreviews, token[:20]+"...")
		} else {
			tokenPreviews = append(tokenPreviews, token)
		}
	}
	logger.Info("批量导入RT - 请求", "batch_name", req.BatchName, "tag", req.Tag, "proxy", req.Proxy, "client_id", req.ClientID, "token_count", len(req.RTTokens), "token_previews", tokenPreviews)

	// 获取代理列表和 Client ID 列表
	proxyList, _ := h.configService.GetProxyList()
	clientIdList, _ := h.configService.GetClientIdList()

	successCount, failCount, err := h.rtService.BatchImport(req.BatchName, req.Tag, req.Proxy, req.ClientID, req.RTTokens, proxyList, clientIdList)
	if err != nil {
		logger.Error("批量导入RT失败", "batch_name", req.BatchName, "tag", req.Tag, "token_count", len(req.RTTokens), "error", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Msg:     "批量导入失败: " + err.Error(),
		})
		return
	}

	logger.Info("批量导入RT完成", "batch_name", req.BatchName, "tag", req.Tag, "total_count", len(req.RTTokens), "success_count", successCount, "fail_count", failCount)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     fmt.Sprintf("批量导入完成: 成功 %d 个, 失败/跳过 %d 个", successCount, failCount),
		Data: gin.H{
			"total_count":   len(req.RTTokens),
			"success_count": successCount,
			"fail_count":    failCount,
		},
	})
}

// RefreshUserInfo 刷新用户信息 - POST /api/rts/refresh-user-info
func (h *RTHandler) RefreshUserInfo(c *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("刷新用户信息 - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     "参数错误: " + err.Error(),
		})
		return
	}

	logger.Info("刷新用户信息 - 请求", "id", req.ID)

	rt, err := h.rtService.RefreshUserInfo(req.ID)
	if err != nil {
		logger.Error("刷新用户信息失败", "id", req.ID, "error", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}

	logger.Info("刷新用户信息成功", "id", req.ID, "user_name", rt.UserName, "email", rt.Email)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "刷新用户信息成功",
		Data:    rt,
	})
}

// RefreshAccountInfo 刷新账号信息 - POST /api/rts/refresh-account-info
func (h *RTHandler) RefreshAccountInfo(c *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("刷新账号信息 - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Msg:     "参数错误: " + err.Error(),
		})
		return
	}

	logger.Info("刷新账号信息 - 请求", "id", req.ID)

	rt, err := h.rtService.RefreshAccountInfo(req.ID)
	if err != nil {
		logger.Error("刷新账号信息失败", "id", req.ID, "error", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}

	logger.Info("刷新账号信息成功", "id", req.ID, "type", rt.Type)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Msg:     "刷新账号信息成功",
		Data:    rt,
	})
}
