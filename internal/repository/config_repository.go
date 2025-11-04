package repository

import (
	"errors"

	"rt-manage/internal/model"
	"gorm.io/gorm"
)

// ConfigRepository 配置数据仓库接口
type ConfigRepository interface {
	GetByKey(key string) (*model.SystemConfig, error)
	Set(key, value string) error
	GetAll() (map[string]string, error)
	BatchSet(configs map[string]string) error
}

type configRepository struct {
	db *gorm.DB
}

// NewConfigRepository 创建配置仓库实例
func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepository{db: db}
}

// GetByKey 根据key获取配置
func (r *configRepository) GetByKey(key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	err := r.db.Where("config_key = ?", key).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

// Set 设置配置
func (r *configRepository) Set(key, value string) error {
	existing, err := r.GetByKey(key)
	if err != nil {
		return err
	}

	if existing != nil {
		// 更新
		existing.ConfigValue = value
		return r.db.Save(existing).Error
	}

	// 创建
	config := &model.SystemConfig{
		ConfigKey:   key,
		ConfigValue: value,
	}
	return r.db.Create(config).Error
}

// GetAll 获取所有配置
func (r *configRepository) GetAll() (map[string]string, error) {
	var configs []model.SystemConfig
	if err := r.db.Find(&configs).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, config := range configs {
		result[config.ConfigKey] = config.ConfigValue
	}
	return result, nil
}

// BatchSet 批量设置配置
func (r *configRepository) BatchSet(configs map[string]string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for key, value := range configs {
			var existing model.SystemConfig
			err := tx.Where("config_key = ?", key).First(&existing).Error
			
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 创建
				newConfig := &model.SystemConfig{
					ConfigKey:   key,
					ConfigValue: value,
				}
				if err := tx.Create(newConfig).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else {
				// 更新
				existing.ConfigValue = value
				if err := tx.Save(&existing).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
