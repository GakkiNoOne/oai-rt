package scheduler

import (
	"strconv"
	"sync"

	"rt-manage/pkg/logger"
)

var (
	globalManager *Manager
	once          sync.Once
)

// RTServiceInterface 定义RT服务接口，避免循环导入
type RTServiceInterface interface {
	AutoRefreshAll() error
}

// Manager 调度器管理器（单例模式）
type Manager struct {
	scheduler   *Scheduler
	rtService   RTServiceInterface
	mu          sync.RWMutex
	running     bool
	intervalDay int
}

// InitManager 初始化全局管理器（只调用一次）
func InitManager(rtService RTServiceInterface, intervalDays int) {
	once.Do(func() {
		globalManager = &Manager{
			rtService:   rtService,
			intervalDay: intervalDays,
		}
		logger.Info("调度器管理器初始化完成", "interval_days", intervalDays)
	})
}

// GetManager 获取全局管理器实例
func GetManager() *Manager {
	if globalManager == nil {
		panic("调度器管理器未初始化，请先调用 InitManager")
	}
	return globalManager
}

// Start 启动调度器
func (m *Manager) Start(intervalDays int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		logger.Warn("调度器已在运行中，无需重复启动")
		return nil
	}

	// 创建新的调度器
	m.scheduler = NewScheduler(m.rtService, intervalDays)
	m.intervalDay = intervalDays
	m.scheduler.Start()
	m.running = true

	logger.Info("调度器已启动", "interval_days", intervalDays)
	return nil
}

// Stop 停止调度器
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		logger.Warn("调度器未运行，无需停止")
		return nil
	}

	if m.scheduler != nil {
		m.scheduler.Stop()
		m.scheduler = nil
	}
	m.running = false

	logger.Info("调度器已停止")
	return nil
}

// Restart 重启调度器（更新间隔）
func (m *Manager) Restart(intervalDays int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果未运行，直接返回
	if !m.running {
		logger.Info("调度器未运行，无需重启")
		return nil
	}

	// 停止旧调度器
	if m.scheduler != nil {
		m.scheduler.Stop()
		logger.Info("旧调度器已停止", "old_interval_days", m.intervalDay)
	}

	// 启动新调度器
	m.scheduler = NewScheduler(m.rtService, intervalDays)
	m.intervalDay = intervalDays
	m.scheduler.Start()

	logger.Info("调度器已重启", "new_interval_days", intervalDays)
	return nil
}

// IsRunning 检查调度器是否运行
func (m *Manager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// GetInterval 获取当前刷新间隔（天）
func (m *Manager) GetInterval() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.intervalDay
}

// UpdateFromConfig 根据配置更新调度器状态
func (m *Manager) UpdateFromConfig(enabled string, intervalStr string) error {
	isEnabled := enabled == "true"
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		logger.Error("解析刷新间隔失败", "interval", intervalStr, "error", err)
		interval = 2 // 默认2天
	}

	wasRunning := m.IsRunning()
	oldInterval := m.GetInterval()

	logger.Info("配置变化检测",
		"enabled", isEnabled,
		"was_running", wasRunning,
		"new_interval", interval,
		"old_interval", oldInterval,
	)

	// 根据状态变化进行操作
	if isEnabled && !wasRunning {
		// 需要启动
		return m.Start(interval)
	} else if !isEnabled && wasRunning {
		// 需要停止
		return m.Stop()
	} else if isEnabled && wasRunning && interval != oldInterval {
		// 间隔变化，需要重启
		return m.Restart(interval)
	}

	logger.Info("调度器状态无需变更")
	return nil
}

