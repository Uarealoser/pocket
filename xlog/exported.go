package xlog

import (
	"Uarealoser/pocket/config"
	"context"
	"log"
	"strings"
)

var (
	defaultLogger Logger = NewLoggerWithConsole()

	TypeFile    = "file"
	TypeConsole = "console"
)

// Config logger配置
type Config struct {
	WriterType string
	Service    string // 服务名称
	File       FileConfig
	Console    ConsoleConfig
}

// FileConfig writer配置
type FileConfig struct {
	LogPath string // 日志文件路径
	LogName string // 日志文件名，无后缀
	Level   Level  // 日志输出等级
	IsJson  bool
}

// ConsoleConfig writer配置
type ConsoleConfig struct {
	Level  Level // 日志输出等级
	IsJson bool
}

// FormatterConfig formatter配置
type FormatterConfig struct {
	HasRuntime bool // 是否打印runtime信息，fatal以上一定会打
	IsJSON     bool // 是否json输出数据，原soa日志格式
}

func GetDefaultLogger() Logger {
	return defaultLogger
}

func SetDefaultLogger(l Logger) {
	if l == nil {
		return
	}
	defaultLogger = l
}

func TakeOverGoLog(l Logger) {
	log.SetFlags(0)
	log.SetOutput(logWriter{l})
}

// InitWithConfig 初始化logger配置
func InitWithConfig(conf Config) Logger {
	defaultLogger = NewLoggerWithConsole()
	defaultLogger.(*XLogger).Service = conf.Service

	if len(conf.WriterType) <= 0 {
		return NewLogger()
	}

	defaultLogger.(*XLogger).flushWriter()
	if strings.Contains(conf.WriterType, TypeFile) {
		defaultLogger.(*XLogger).Writers[TypeFile] = NewFileWriter(conf.File)
	}
	if strings.Contains(conf.WriterType, TypeConsole) {
		defaultLogger.(*XLogger).Writers[TypeConsole] = NewConsoleWriter(conf.Console)
	}
	return defaultLogger
}

func addWriterByXconfig(config config.Config) (file FileConfig, console ConsoleConfig, err error) {
	section := "log"
	if _, err = config.GetSection(section); err != nil {
		return
	}

	file.LogPath = config.GetStringDefault(section, "path", defaultLogPath)
	file.LogName = config.GetStringDefault(section, "log_name", defaultLogName)
	level, err2 := config.GetString(section, "level")
	if err2 != nil {
		file.Level = DefaultLevel
		console.Level = DefaultLevel
	} else {
		file.Level = StringToLevel(level)
		console.Level = StringToLevel(level)
	}
	return
}

// InitWithXConfig .
func InitWithXConfig(config config.Config, groupName, serviceName string) (l Logger, err error) {
	file, console, err := addWriterByXconfig(config)
	if err != nil {
		return
	}
	return InitWithConfig(Config{
		WriterType: config.GetStringDefault("log", "type", ""),
		Service:    groupName + "." + serviceName,
		File:       file,
		Console:    console,
	}), nil
}

func Warn(ctx context.Context, format string, a ...interface{}) error {
	l, ok := defaultLogger.(*XLogger)
	if !ok {
		return defaultLogger.Warn(ctx, format, a...)
	}
	return l.getEntry().Warn(ctx, format, a...)
}

func Error(ctx context.Context, format string, a ...interface{}) error {
	l, ok := defaultLogger.(*XLogger)
	if !ok {
		return defaultLogger.Error(ctx, format, a...)
	}
	return l.getEntry().Error(ctx, format, a...)
}

func Fatal(ctx context.Context, format string, a ...interface{}) error {
	l, ok := defaultLogger.(*XLogger)
	if !ok {
		return defaultLogger.Fatal(ctx, format, a...)
	}
	return l.getEntry().Fatal(ctx, format, a...)
}

func Notice(ctx context.Context, format string, a ...interface{}) error {
	l, ok := defaultLogger.(*XLogger)
	if !ok {
		return defaultLogger.Notice(ctx, format, a...)
	}
	return l.getEntry().Notice(ctx, format, a...)
}

func Debug(ctx context.Context, format string, a ...interface{}) error {
	l, ok := defaultLogger.(*XLogger)
	if !ok {
		return defaultLogger.Debug(ctx, format, a...)
	}
	return l.getEntry().Debug(ctx, format, a...)
}

func With(key, value string) Entry {
	return defaultLogger.With(key, value)
}

func Close() {
	defaultLogger.Close()
}

func NewNopLogger() Logger {
	return &nop{}
}

type logWriter struct {
	l Logger
}

func (w logWriter) Write(p []byte) (n int, err error) {
	logger := w.l
	if logger == nil {
		logger = defaultLogger
	}
	if logger == nil {
		logger = NewLogger()
	}
	logger.Notice(nil, "%s", string(p))
	n = len(p)
	return
}
