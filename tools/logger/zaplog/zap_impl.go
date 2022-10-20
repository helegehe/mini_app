package zaplog

import (
	"github.com/helegehe/mini_app/internal/pkg"
	"go.uber.org/zap"
)

// LogLevel 日志等级
type LogLevel int

const (
	logDebug LogLevel = iota
	logInfo
	logWarn
	logErr
	logOff
)

var (
	// 字符串和级别映射
	levelMapping = map[string]LogLevel{
		"debug": logDebug,
		"info":  logInfo,
		"warn":  logWarn,
		"err":   logErr,
		"off":   logOff,
	}
)

// InitLogger 初始化日志框架
func InitLogger(projectName string, conf *pkg.Logger) *ZapLogger {
	myLogger := &ZapLogger{}
	myLogger.Logger = InitZap(projectName, conf.Path, conf.MaxAge, conf.RotationTime)
	myLogger.LogLevel = levelMapping[conf.Level]
	myLogger.Slog = myLogger.Logger.Sugar()
	return myLogger
}

// InitLoggerForTest 初始化日志不需要传入参数，用于测试，不会打印文件
func InitLoggerForTest() *ZapLogger {
	myLogger := &ZapLogger{}
	l, _ := zap.NewDevelopment()
	myLogger.Logger = l
	myLogger.Slog = l.Sugar()
	myLogger.LogLevel = logDebug
	return myLogger
}

// ZapLogger zap具体实现接口
type ZapLogger struct {
	Logger   *zap.Logger
	Slog     *zap.SugaredLogger
	LogLevel LogLevel
}

// Debug log
func (z *ZapLogger) Debug(v ...interface{}) {
	if z.LogLevel <= logDebug {
		z.Slog.Debug(v...)
	}
}

// Debugf log
func (z *ZapLogger) Debugf(format string, v ...interface{}) {
	if z.LogLevel <= logDebug {
		z.Slog.Debugf(format, v...)
	}
}

// Info log
func (z *ZapLogger) Info(v ...interface{}) {
	if z.LogLevel <= logInfo {
		z.Slog.Info(v...)
	}
}

// Infof log
func (z *ZapLogger) Infof(format string, v ...interface{}) {
	if z.LogLevel <= logInfo {
		z.Slog.Infof(format, v...)
	}
}

// Warn log
func (z *ZapLogger) Warn(v ...interface{}) {
	if z.LogLevel <= logWarn {
		z.Slog.Warn(v...)
	}
}

// Warnf log
func (z *ZapLogger) Warnf(format string, v ...interface{}) {
	if z.LogLevel <= logWarn {
		z.Slog.Warnf(format, v...)
	}
}

// Error log
func (z *ZapLogger) Error(v ...interface{}) {
	if z.LogLevel <= logErr {
		z.Slog.Error(v...)
	}
}

// Errorf log
func (z *ZapLogger) Errorf(format string, v ...interface{}) {
	if z.LogLevel <= logErr {
		z.Slog.Errorf(format, v...)
	}
}
