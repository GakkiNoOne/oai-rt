package model

import (
	"time"
)

var tablePrefix string

// SetTablePrefix 设置表名前缀
func SetTablePrefix(prefix string) {
	tablePrefix = prefix
}

// GetTablePrefix 获取表名前缀
func GetTablePrefix() string {
	return tablePrefix
}

// withPrefix 为表名添加前缀
func withPrefix(tableName string) string {
	if tablePrefix == "" {
		return tableName
	}
	return tablePrefix + "_" + tableName
}

// RT 存储 RT token 的模型
type RT struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	BizId           string    `json:"biz_id" gorm:"type:varchar(255);uniqueIndex:idx_biz_id;not null"`
	UserName        string    `json:"user_name" gorm:"type:varchar(255)"`
	Email           string    `json:"email" gorm:"type:varchar(255)"`
	Type            string    `json:"type" gorm:"type:varchar(50)"`
	Rt              string    `json:"rt" gorm:"type:text;not null"`
	At              string    `json:"at" gorm:"type:text"`
	Proxy           string    `json:"proxy" gorm:"type:varchar(255)"`
	ClientID        string    `json:"client_id" gorm:"type:varchar(255)"`
	Tag             string    `json:"tag" gorm:"type:varchar(255)"`
	Enabled         bool      `json:"enabled" gorm:"default:true;not null"`
	LastRT          string     `json:"last_rt" gorm:"type:text"`
	RefreshResult   string     `json:"refresh_result" gorm:"type:text"`
	UserInfo        string     `json:"user_info" gorm:"type:text"`
	AccountInfo     string     `json:"account_info" gorm:"type:text"`
	LastRefreshTime *time.Time `json:"last_refresh_time" gorm:"type:datetime;default:null"`
	Memo            string     `json:"memo" gorm:"type:text"`
	CreateTime      time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime      time.Time `json:"update_time" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (RT) TableName() string {
	return withPrefix("rts")
}

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ConfigKey   string    `json:"config_key" gorm:"type:varchar(255);uniqueIndex:idx_config_key;not null"`
	ConfigValue string    `json:"config_value" gorm:"type:text"`
	CreateTime  time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return withPrefix("system_configs")
}
