package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rt-manage/internal/model"
	"rt-manage/internal/service"
	"rt-manage/pkg/logger"
)

// PublicAPIHandler 对外API处理器
type PublicAPIHandler struct {
	rtService service.RTService
}

// NewPublicAPIHandler 创建对外API处理器
func NewPublicAPIHandler(rtService service.RTService) *PublicAPIHandler {
	return &PublicAPIHandler{
		rtService: rtService,
	}
}

// RefreshAndGetAT 刷新RT并获取AT - POST /public-api/refresh
func (h *PublicAPIHandler) RefreshAndGetAT(c *gin.Context) {
	var req struct {
		BizId string `json:"biz_id"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("RefreshAndGetAT - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "参数错误: " + err.Error(),
		})
		return
	}

	// 优先使用 biz_id，其次 email
	if req.BizId == "" && req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "biz_id 和 email 至少提供一个",
		})
		return
	}

	logger.Info("RefreshAndGetAT - 请求", "biz_id", req.BizId, "email", req.Email)

	// 根据条件查找RT
	var rt *model.RT
	var err error
	
	if req.BizId != "" {
		// 优先使用 biz_id
		rt, err = h.rtService.GetByBizId(req.BizId)
	} else {
		// 使用 email
		rt, err = h.rtService.GetByEmail(req.Email)
	}

	if err != nil {
		logger.Error("RefreshAndGetAT - 查找RT失败", "biz_id", req.BizId, "email", req.Email, "error", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"msg":     "RT不存在",
		})
		return
	}

	// 刷新RT（不刷新用户信息和账号信息）
	refreshedRT, err := h.rtService.Refresh(rt.ID, false, false)
	if err != nil {
		logger.Error("RefreshAndGetAT - 刷新失败", "id", rt.ID, "biz_id", rt.BizId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"msg":     "刷新失败: " + err.Error(),
			"data": gin.H{
				"refresh_result": rt.RefreshResult,
			},
		})
		return
	}

	logger.Info("RefreshAndGetAT - 成功", "id", refreshedRT.ID, "biz_id", refreshedRT.BizId)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "刷新成功",
		"data": gin.H{
			"biz_id":      refreshedRT.BizId,
			"email":       refreshedRT.Email,
			"access_token": refreshedRT.At,
			"refresh_token": refreshedRT.Rt,
			"type":        refreshedRT.Type,
			"user_name":   refreshedRT.UserName,
		},
	})
}

// GetAT 获取AT（不刷新）- POST /public-api/get-at
func (h *PublicAPIHandler) GetAT(c *gin.Context) {
	var req struct {
		BizId string `json:"biz_id"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("GetAT - 参数错误", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "参数错误: " + err.Error(),
		})
		return
	}

	// 优先使用 biz_id，其次 email
	if req.BizId == "" && req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "biz_id 和 email 至少提供一个",
		})
		return
	}

	logger.Info("GetAT - 请求", "biz_id", req.BizId, "email", req.Email)

	// 根据条件查找RT
	var rt *model.RT
	var err error
	
	if req.BizId != "" {
		// 优先使用 biz_id
		rt, err = h.rtService.GetByBizId(req.BizId)
	} else {
		// 使用 email
		rt, err = h.rtService.GetByEmail(req.Email)
	}

	if err != nil {
		logger.Error("GetAT - 查找RT失败", "biz_id", req.BizId, "email", req.Email, "error", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"msg":     "RT不存在",
		})
		return
	}

	logger.Info("GetAT - 成功", "id", rt.ID, "biz_id", rt.BizId)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "获取成功",
		"data": gin.H{
			"biz_id":      rt.BizId,
			"email":       rt.Email,
			"access_token": rt.At,
			"refresh_token": rt.Rt,
			"type":        rt.Type,
			"user_name":   rt.UserName,
		},
	})
}

