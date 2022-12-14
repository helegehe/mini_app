package zaplog

import (
	"os"
	"path"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitZap logPath: 日志打印目录
// maxAge: 日志最大存在时间，单位：天
// rotationTime: 日志切分时间，单位：小时
// projectName: 项目名称
func InitZap(projectName, logPath string, maxAge, rotationTime time.Duration) *zap.Logger {
	if len(projectName) == 0 {
		panic("logger init fail, project name is empty")
	}
	maxAge = maxAge * 24 * time.Hour
	rotationTime = rotationTime * time.Hour

	// 创建日志存放目录
	_, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(logPath, os.ModePerm)
		if err != nil{
			panic(err)
		}
	}
	logPath = path.Join(logPath, projectName)

	// error日志文件配置
	errWriter, err := rotatelogs.New(
		logPath+"_err_%Y-%m-%d.log",
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotationTime),
	)
	if err != nil {
		panic(err)
	}

	// info日志文件配置
	infoWriter, err := rotatelogs.New(
		logPath+"_info_%Y-%m-%d.log",
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotationTime),
	)
	if err != nil {
		panic(err)
	}

	// 优先级设置（一个日志输出全部信息，一个日志输出error信息）
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	// 控制台输出设置
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoderConfig.EncodeTime = timeEncoder
	consoleEncoderConfig.EncodeCaller = customCallerEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// 文件输出设置
	errorCore := zapcore.AddSync(errWriter)
	infoCore := zapcore.AddSync(infoWriter)
	fileEncodeConfig := zap.NewProductionEncoderConfig()
	fileEncodeConfig.EncodeTime = timeEncoder
	fileEncodeConfig.EncodeCaller = customCallerEncoder
	fileEncoder := zapcore.NewConsoleEncoder(fileEncodeConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, errorCore, highPriority),
		zapcore.NewCore(fileEncoder, infoCore, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, zapcore.DebugLevel),
	)

	// 显示行号
	caller := zap.AddCaller()
	// 设置打印堆栈深度
	callerSkip := zap.AddCallerSkip(2)

	development := zap.Development()
	logger := zap.New(core, caller, callerSkip, development)

	// 替换全局日志
	zap.ReplaceGlobals(logger)

	// 将系统输出重定向到zap中，保证所有出现异常均能打印到文件中
	if _, err := zap.RedirectStdLogAt(logger, zapcore.ErrorLevel); err != nil {
		panic(err)
	}

	return logger
}

// 自定义打印路径，减少输出日志打印路径长度，根据输入项目名进行减少
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	str := caller.String()
	index := strings.Index(str, "src")
	if index == -1 {
		enc.AppendString(caller.FullPath())
	} else {
		index = index + len("src") + 1
		enc.AppendString(str[index:])
	}
}

// 格式化日志时间，官方的不好看
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
