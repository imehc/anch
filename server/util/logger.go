package util

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Logger 全局日志记录器
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

// LogConfig 日志配置
type LogConfig struct {
	LogDir     string // 日志目录
	MaxSizeMB  int    // 单个文件最大大小（MB）
	MaxBackups int    // 保留的旧日志文件数量
	MaxAgeDays int    // 保留日志文件的最大天数
	Compress   bool   // 是否压缩旧日志
	Level      string // 日志级别: debug, info, warn, error
}

// InitLogger 初始化日志系统
func InitLogger(config LogConfig) error {
	// 设置默认值
	if config.LogDir == "" {
		config.LogDir = "logs"
	}
	if config.MaxSizeMB == 0 {
		config.MaxSizeMB = 10 // 默认 10MB
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 30 // 默认保留30个备份
	}
	if config.MaxAgeDays == 0 {
		config.MaxAgeDays = 30 // 默认保留30天
	}
	if config.Level == "" {
		config.Level = "info"
	}

	// 创建日志目录
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// 获取当前日期作为日志文件名的一部分
	today := time.Now().Format("2006-01-02")

	// 配置日志级别
	var level zapcore.Level
	switch config.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Info 日志文件配置
	infoWriter := &lumberjack.Logger{
		Filename:   filepath.Join(config.LogDir, fmt.Sprintf("info-%s.log", today)),
		MaxSize:    config.MaxSizeMB,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAgeDays,
		Compress:   config.Compress,
		LocalTime:  true,
	}

	// Error 日志文件配置
	errorWriter := &lumberjack.Logger{
		Filename:   filepath.Join(config.LogDir, fmt.Sprintf("error-%s.log", today)),
		MaxSize:    config.MaxSizeMB,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAgeDays,
		Compress:   config.Compress,
		LocalTime:  true,
	}

	// 创建 core
	// Info 及以上级别写入 info 日志文件
	infoCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(infoWriter),
		level,
	)

	// Error 及以上级别写入 error 日志文件
	errorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(errorWriter),
		zapcore.ErrorLevel,
	)

	// 同时输出到控制台（开发环境）
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	// 合并多个 core
	core := zapcore.NewTee(infoCore, errorCore, consoleCore)

	// 创建 logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	Sugar = Logger.Sugar()

	Sugar.Info("Logger initialized successfully")
	return nil
}

// Info 记录信息日志
func Info(format string, args ...any) {
	if Sugar != nil {
		Sugar.Infof(format, args...)
	}
}

// Error 记录错误日志
func Error(format string, args ...any) {
	if Sugar != nil {
		Sugar.Errorf(format, args...)
	}
}

// Debug 记录调试日志
func Debug(format string, args ...any) {
	if Sugar != nil {
		Sugar.Debugf(format, args...)
	}
}

// Warn 记录警告日志
func Warn(format string, args ...any) {
	if Sugar != nil {
		Sugar.Warnf(format, args...)
	}
}

// LogError 记录错误及其上下文信息
func LogError(operation string, err error, context map[string]any) {
	if Sugar == nil {
		return
	}

	fields := []any{
		"operation", operation,
		"error", err.Error(),
	}

	for k, v := range context {
		fields = append(fields, k, v)
	}

	Sugar.Errorw("Operation failed", fields...)
}

// Close 关闭日志系统
func CloseLogger() {
	if Logger != nil {
		Logger.Sync()
	}
}
