package scheduler

import (
	"context"
	"time"

	"rt-manage/pkg/logger"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	rtService RTServiceInterface
	interval  time.Duration
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewScheduler 创建调度器实例
func NewScheduler(rtService RTServiceInterface, intervalDays int) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		rtService: rtService,
		interval:  time.Duration(intervalDays) * 24 * time.Hour,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start 启动定时任务
func (s *Scheduler) Start() {
	logger.Info("启动定时刷新任务", "interval", s.interval)

	// 立即执行一次
	go func() {
		if err := s.rtService.AutoRefreshAll(); err != nil {
			logger.Error("自动刷新失败", "error", err)
		}
	}()

	// 定时执行
	ticker := time.NewTicker(s.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				logger.Info("执行定时刷新任务")
				if err := s.rtService.AutoRefreshAll(); err != nil {
					logger.Error("自动刷新失败", "error", err)
				}
			case <-s.ctx.Done():
				ticker.Stop()
				logger.Info("定时刷新任务已停止")
				return
			}
		}
	}()
}

// Stop 停止定时任务
func (s *Scheduler) Stop() {
	s.cancel()
}
