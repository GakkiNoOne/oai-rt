package database

import (
	"fmt"
	"time"

	"rt-manage/internal/config"
	"rt-manage/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	
	// 使用纯 Go 的 SQLite 驱动（不需要 CGO）
	sqlite "github.com/glebarez/sqlite"
)

var db *gorm.DB

// Init 初始化数据库连接
func Init(cfg *config.DatabaseConfig) error {
	var dialector gorm.Dialector
	var err error

	// 设置表前缀
	model.SetTablePrefix(cfg.TablePrefix)

	// 根据数据库类型选择驱动
	switch cfg.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
		dialector = mysql.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", cfg.Type)
	}

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// 自动迁移数据库表（只创建表和列，不处理索引变更）
	migrator := db.Migrator()
	
	// 检查表是否存在，不存在才自动迁移
	if !migrator.HasTable(&model.RT{}) {
		if err := db.AutoMigrate(&model.RT{}); err != nil {
			return fmt.Errorf("创建 rt_rts 表失败: %w", err)
		}
	}
	
	if !migrator.HasTable(&model.SystemConfig{}) {
		if err := db.AutoMigrate(&model.SystemConfig{}); err != nil {
			return fmt.Errorf("创建 system_configs 表失败: %w", err)
		}
	}

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}

// Close 关闭数据库连接
func Close() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
