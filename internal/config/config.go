package config

import (
	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Log      LogConfig      `mapstructure:"log"`
	Database DatabaseConfig `mapstructure:"database"`
	OpenAI   OpenAIConfig   `mapstructure:"openai"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level            string   `mapstructure:"level"`
	Encoding         string   `mapstructure:"encoding"`
	OutputPaths      []string `mapstructure:"output_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            string `mapstructure:"type"` // sqlite 或 mysql
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	TablePrefix     string `mapstructure:"table_prefix"` // 表名前缀，如 "rt" 则表名为 rt_rts, rt_system_configs
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// OpenAIConfig OpenAI 配置
type OpenAIConfig struct {
	ClientID        string `mapstructure:"client_id"`
	Proxy           string `mapstructure:"proxy"`
	RefreshInterval int    `mapstructure:"refresh_interval"` // 自动刷新间隔（天）
	ScheduleEnabled bool   `mapstructure:"schedule_enabled"` // 是否启用定时刷新
}

// AuthConfig 认证配置
type AuthConfig struct {
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
	JWTSecret      string `mapstructure:"jwt_secret"`
	JWTExpireHours int    `mapstructure:"jwt_expire_hours"`
	APISecret      string `mapstructure:"api_secret"`
}

var cfg *Config

// Init 初始化配置
func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")

	// 日志默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.encoding", "json")
	viper.SetDefault("log.output_paths", []string{"stdout"})
	viper.SetDefault("log.error_output_paths", []string{"stderr"})

	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.database", "./data/tokens.db")
	viper.SetDefault("database.table_prefix", "")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 3600)
	viper.SetDefault("openai.client_id", "app_WXrF1LSkiTtfYqiL6XtjygvX")
	viper.SetDefault("openai.refresh_interval", 2) // 默认2天
	viper.SetDefault("openai.schedule_enabled", false)
	viper.SetDefault("auth.username", "admin")
	viper.SetDefault("auth.password", "admin123")
	viper.SetDefault("auth.jwt_secret", "your-secret-key-change-this-in-production")
	viper.SetDefault("auth.jwt_expire_hours", 24)
	viper.SetDefault("auth.api_secret", "my-api-secret-2025")

	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，使用默认值
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}

// Get 获取配置
func Get() *Config {
	return cfg
}
