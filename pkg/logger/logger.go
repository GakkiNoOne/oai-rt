package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"rt-manage/internal/config"
)

var log *zap.Logger

// Init 初始化日志
func Init() error {
	cfg := config.Get()
	
	var zapConfig zap.Config
	if cfg.Log.Encoding == "json" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		// 自定义 console 编码配置，让堆栈信息更清晰
		zapConfig.EncoderConfig.StacktraceKey = "stacktrace"
		zapConfig.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
	}

	zapConfig.Level = zap.NewAtomicLevelAt(getLogLevel(cfg.Log.Level))
	zapConfig.Encoding = cfg.Log.Encoding
	zapConfig.OutputPaths = cfg.Log.OutputPaths
	zapConfig.ErrorOutputPaths = cfg.Log.ErrorOutputPaths

	var err error
	log, err = zapConfig.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return err
	}

	return nil
}

// getLogLevel 获取日志级别
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// Debug 记录debug日志
func Debug(msg string, fields ...interface{}) {
	log.Sugar().Debugw(msg, fields...)
}

// Info 记录info日志
func Info(msg string, fields ...interface{}) {
	log.Sugar().Infow(msg, fields...)
}

// Warn 记录warn日志
func Warn(msg string, fields ...interface{}) {
	log.Sugar().Warnw(msg, fields...)
}

// Error 记录error日志
func Error(msg string, fields ...interface{}) {
	log.Sugar().Errorw(msg, fields...)
}

// Fatal 记录fatal日志
func Fatal(msg string, fields ...interface{}) {
	log.Sugar().Fatalw(msg, fields...)
}

// Sync 刷新日志缓冲
func Sync() {
	log.Sync()
}

