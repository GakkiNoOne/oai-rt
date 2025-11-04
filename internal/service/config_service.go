package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"rt-manage/internal/repository"
	"rt-manage/internal/scheduler"
	"rt-manage/pkg/logger"
)

// ConfigService 配置服务接口
type ConfigService interface {
	GetSystemConfigs() (map[string]string, map[string]string, error)
	SaveSystemConfigs(configs map[string]string) error
	GetProxyList() ([]string, error)
	GetClientIdList() ([]string, error)
	GetConfig(key string) (string, error)
	SetConfig(key, value string) error
}

type configService struct {
	repo   repository.ConfigRepository
	rtRepo repository.RTRepository
}

// NewConfigService 创建配置服务实例
func NewConfigService(repo repository.ConfigRepository, rtRepo repository.RTRepository) ConfigService {
	return &configService{
		repo:   repo,
		rtRepo: rtRepo,
	}
}

// GetSystemConfigs 获取系统配置
func (s *configService) GetSystemConfigs() (map[string]string, map[string]string, error) {
	// 获取数据库配置
	dbConfigs, err := s.repo.GetAll()
	if err != nil {
		return nil, nil, err
	}

	// 设置默认值
	if dbConfigs == nil {
		dbConfigs = make(map[string]string)
	}

	// 应用默认值
	defaults := map[string]string{
		"proxy_list":            "[]",
		"client_id_list":        "[]",
		"auto_refresh_enabled":  "false",
		"auto_refresh_interval": "60",
	}

	for key, defaultValue := range defaults {
		if _, exists := dbConfigs[key]; !exists {
			dbConfigs[key] = defaultValue
		}
	}

	// 获取环境变量配置
	envConfigs := map[string]string{
		"API_PREFIX":     getEnv("API_PREFIX", "/api"),
		"API_SECRET":     getEnv("API_SECRET", "your-api-secret"),
		"ADMIN_PREFIX":   getEnv("ADMIN_PREFIX", "/admin"),
		"ADMIN_USERNAME": getEnv("ADMIN_USERNAME", "admin"),
		"ADMIN_PASSWORD": getEnv("ADMIN_PASSWORD", "admin123"),
	}

	return dbConfigs, envConfigs, nil
}

// SaveSystemConfigs 保存系统配置
func (s *configService) SaveSystemConfigs(configs map[string]string) error {
	// 验证配置
	if err := s.validateConfigs(configs); err != nil {
		return err
	}

	// 获取旧配置用于对比
	oldEnabled, _ := s.GetConfig("auto_refresh_enabled")
	oldInterval, _ := s.GetConfig("auto_refresh_interval")
	oldProxyList, _ := s.GetConfig("proxy_list")
	oldClientIdList, _ := s.GetConfig("client_id_list")

	// 保存配置
	if err := s.repo.BatchSet(configs); err != nil {
		return err
	}

	// 检查 proxy_list 或 client_id_list 是否变化
	newProxyList := configs["proxy_list"]
	newClientIdList := configs["client_id_list"]

	proxyChanged := newProxyList != oldProxyList
	clientIdChanged := newClientIdList != oldClientIdList

	if proxyChanged || clientIdChanged {
		logger.Info("代理或ClientID配置发生变化，开始批量更新RT",
			"proxy_changed", proxyChanged,
			"client_id_changed", clientIdChanged,
		)

		// 批量更新所有RT的proxy和clientId
		if err := s.updateAllRTConfigs(newProxyList, newClientIdList); err != nil {
			logger.Error("批量更新RT配置失败", "error", err)
			// 不返回错误，配置已保存
		}
	}

	// 检查自动刷新配置是否变化，动态更新调度器
	newEnabled := configs["auto_refresh_enabled"]
	newInterval := configs["auto_refresh_interval"]

	if newEnabled != oldEnabled || newInterval != oldInterval {
		logger.Info("自动刷新配置变化",
			"old_enabled", oldEnabled,
			"new_enabled", newEnabled,
			"old_interval", oldInterval,
			"new_interval", newInterval,
		)

		// 动态更新调度器
		if err := s.updateScheduler(newEnabled, newInterval); err != nil {
			logger.Error("更新调度器失败", "error", err)
			// 不返回错误，配置已保存，调度器更新失败不影响配置保存
		}
	}

	return nil
}

// updateScheduler 更新调度器状态
func (s *configService) updateScheduler(enabled, intervalStr string) error {
	return scheduler.GetManager().UpdateFromConfig(enabled, intervalStr)
}

// GetProxyList 获取代理列表
func (s *configService) GetProxyList() ([]string, error) {
	config, err := s.repo.GetByKey("proxy_list")
	if err != nil {
		return nil, err
	}

	if config == nil || config.ConfigValue == "" {
		return []string{}, nil
	}

	var proxyList []string
	if err := json.Unmarshal([]byte(config.ConfigValue), &proxyList); err != nil {
		return nil, fmt.Errorf("解析代理列表失败: %w", err)
	}

	return proxyList, nil
}

// GetClientIdList 获取 Client ID 列表
func (s *configService) GetClientIdList() ([]string, error) {
	config, err := s.repo.GetByKey("client_id_list")
	if err != nil {
		return nil, err
	}

	if config == nil || config.ConfigValue == "" {
		// 返回默认 Client ID
		return []string{"app_WXrF1LSkiTtfYqiL6XtjygvX"}, nil
	}

	var clientIdList []string
	if err := json.Unmarshal([]byte(config.ConfigValue), &clientIdList); err != nil {
		return nil, fmt.Errorf("解析 Client ID 列表失败: %w", err)
	}

	// 如果列表为空，返回默认值
	if len(clientIdList) == 0 {
		return []string{"app_WXrF1LSkiTtfYqiL6XtjygvX"}, nil
	}

	return clientIdList, nil
}

// GetConfig 获取单个配置
func (s *configService) GetConfig(key string) (string, error) {
	config, err := s.repo.GetByKey(key)
	if err != nil {
		return "", err
	}
	if config == nil {
		return "", nil
	}
	return config.ConfigValue, nil
}

// SetConfig 设置单个配置
func (s *configService) SetConfig(key, value string) error {
	return s.repo.Set(key, value)
}

// validateConfigs 验证配置
func (s *configService) validateConfigs(configs map[string]string) error {
	// 验证代理列表格式
	if proxyListStr, ok := configs["proxy_list"]; ok {
		var proxyList []string
		if err := json.Unmarshal([]byte(proxyListStr), &proxyList); err != nil {
			return fmt.Errorf("代理列表格式错误: %w", err)
		}
	}

	// 验证并确保 Client ID 列表至少有一个值
	if clientIdListStr, ok := configs["client_id_list"]; ok {
		var clientIdList []string
		if err := json.Unmarshal([]byte(clientIdListStr), &clientIdList); err != nil {
			return fmt.Errorf("Client ID 列表格式错误: %w", err)
		}

		// 如果列表为空，添加默认值
		if len(clientIdList) == 0 {
			clientIdList = []string{"app_WXrF1LSkiTtfYqiL6XtjygvX"}
			defaultClientIdStr, _ := json.Marshal(clientIdList)
			configs["client_id_list"] = string(defaultClientIdStr)
			logger.Warn("Client ID 列表为空，自动添加默认值")
		}
	}

	return nil
}

// updateAllRTConfigs 批量更新所有RT的proxy和clientId配置
func (s *configService) updateAllRTConfigs(proxyListStr, clientIdListStr string) error {
	// 解析代理列表
	var proxyList []string
	if err := json.Unmarshal([]byte(proxyListStr), &proxyList); err != nil {
		return fmt.Errorf("解析代理列表失败: %v", err)
	}

	// 解析ClientID列表
	var clientIdList []string
	if err := json.Unmarshal([]byte(clientIdListStr), &clientIdList); err != nil {
		return fmt.Errorf("解析ClientID列表失败: %v", err)
	}

	// 确保ClientID列表至少有一个默认值
	if len(clientIdList) == 0 {
		clientIdList = []string{"app_WXrF1LSkiTtfYqiL6XtjygvX"}
		logger.Warn("ClientID列表为空，使用默认值")
	}

	// 获取所有RT（不限制启用状态）
	rts, total, err := s.rtRepo.List(1, 100000, "", "", "", "", nil, "")
	if err != nil {
		return fmt.Errorf("获取RT列表失败: %v", err)
	}

	logger.Info("开始批量更新RT配置", "total_count", total, "proxy_count", len(proxyList), "client_id_count", len(clientIdList))

	updatedCount := 0
	for _, rt := range rts {
		needUpdate := false
		oldProxy := rt.Proxy
		oldClientId := rt.ClientID

		// 检查并更新Proxy
		if len(proxyList) == 0 {
			// 如果代理列表为空，清空所有RT的代理
			if rt.Proxy != "" {
				rt.Proxy = ""
				needUpdate = true
			}
		} else {
			// 检查RT的代理是否在新列表中
			proxyInList := false
			for _, proxy := range proxyList {
				if rt.Proxy == proxy {
					proxyInList = true
					break
				}
			}
			// 如果不在列表中，随机选择一个新代理
			if !proxyInList {
				rt.Proxy = proxyList[rand.Intn(len(proxyList))]
				needUpdate = true
			}
		}

		// 检查并更新ClientID
		clientIdInList := false
		for _, clientId := range clientIdList {
			if rt.ClientID == clientId {
				clientIdInList = true
				break
			}
		}
		// 如果不在列表中，随机选择一个新ClientID
		if !clientIdInList {
			rt.ClientID = clientIdList[rand.Intn(len(clientIdList))]
			needUpdate = true
		}

		// 如果需要更新，保存到数据库
		if needUpdate {
			if err := s.rtRepo.Update(rt); err != nil {
				logger.Error("更新RT配置失败", "rt_id", rt.ID, "name", rt.BizId, "error", err)
			} else {
				updatedCount++
				logger.Info("更新RT配置成功",
					"rt_id", rt.ID,
					"name", rt.BizId,
					"old_proxy", oldProxy,
					"new_proxy", rt.Proxy,
					"old_client_id", oldClientId,
					"new_client_id", rt.ClientID,
				)
			}
		}
	}

	logger.Info("批量更新RT配置完成", "updated_count", updatedCount, "total_count", total)
	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// boolPtr 返回bool指针
func boolPtr(b bool) *bool {
	return &b
}
