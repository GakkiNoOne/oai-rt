package repository

import (
	"errors"
	"time"

	"rt-manage/internal/model"

	"gorm.io/gorm"
)

// RTRepository RT 数据仓库接口
type RTRepository interface {
	Create(rt *model.RT) error
	Update(rt *model.RT) error
	GetByID(id int64) (*model.RT, error)
	GetByBizId(bizId string) (*model.RT, error)
	GetByEmail(email string) (*model.RT, error)
	List(page, pageSize int, bizId string, tag string, email string, typeStr string, enabled *bool, createDate string) ([]*model.RT, int64, error)
	Delete(id int64) error
	BatchDelete(ids []int64) (int, int, error)
	GetByIDs(ids []int64) ([]*model.RT, error)
	GetByToken(token string) (*model.RT, error)
}

type rtRepository struct {
	db *gorm.DB
}

// NewRTRepository 创建 RT 仓库实例
func NewRTRepository(db *gorm.DB) RTRepository {
	return &rtRepository{db: db}
}

// Create 创建新的 RT 记录
func (r *rtRepository) Create(rt *model.RT) error {
	return r.db.Create(rt).Error
}

// Update 更新 RT 记录
func (r *rtRepository) Update(rt *model.RT) error {
	return r.db.Save(rt).Error
}

// GetByID 根据 ID 获取 RT
func (r *rtRepository) GetByID(id int64) (*model.RT, error) {
	var rt model.RT
	err := r.db.Where("id = ?", id).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

// GetByBizId 根据业务ID获取 RT
func (r *rtRepository) GetByBizId(bizId string) (*model.RT, error) {
	var rt model.RT
	err := r.db.Where("biz_id = ?", bizId).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

// GetByEmail 根据邮箱获取 RT
func (r *rtRepository) GetByEmail(email string) (*model.RT, error) {
	var rt model.RT
	err := r.db.Where("email = ?", email).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

// List 获取 RT 列表
func (r *rtRepository) List(page, pageSize int, bizId string, tag string, email string, typeStr string, enabled *bool, createDate string) ([]*model.RT, int64, error) {
	var rts []*model.RT
	var total int64

	query := r.db.Model(&model.RT{})

	// 应用筛选条件
	if bizId != "" {
		query = query.Where("biz_id LIKE ?", "%"+bizId+"%")
	}
	if tag != "" {
		query = query.Where("tag LIKE ?", "%"+tag+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}
	if typeStr != "" {
		query = query.Where("type LIKE ?", "%"+typeStr+"%")
	}
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}
	if createDate != "" {
		// 按日期筛选（忽略时间部分）
		startTime, _ := time.Parse("2006-01-02", createDate)
		endTime := startTime.Add(24 * time.Hour)
		query = query.Where("create_time >= ? AND create_time < ?", startTime, endTime)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&rts).Error; err != nil {
		return nil, 0, err
	}

	return rts, total, nil
}

// Delete 删除 RT 记录
func (r *rtRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.RT{}).Error
}

// BatchDelete 批量删除 RT
func (r *rtRepository) BatchDelete(ids []int64) (int, int, error) {
	successCount := 0
	failCount := 0
	var lastErr error

	for _, id := range ids {
		if err := r.Delete(id); err != nil {
			failCount++
			lastErr = err
		} else {
			successCount++
		}
	}

	return successCount, failCount, lastErr
}

// GetByIDs 根据ID列表获取RT
func (r *rtRepository) GetByIDs(ids []int64) ([]*model.RT, error) {
	var rts []*model.RT
	if err := r.db.Where("id IN ?", ids).Find(&rts).Error; err != nil {
		return nil, err
	}
	return rts, nil
}

// GetByToken 根据token获取RT
func (r *rtRepository) GetByToken(token string) (*model.RT, error) {
	var rt model.RT
	err := r.db.Where("rt = ?", token).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}
