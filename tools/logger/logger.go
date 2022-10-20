package logger

// Logger 系统logger接口
// 后续如果需要替换日志框架，只要实现当前接口，替换具体实现和初始化就可以了
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}
