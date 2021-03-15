package xlog

import (
	"Uarealoser/pocket/config"
	"context"
	"io"
	"os"
	"strings"
	"sync"
)

const (
	DefaultService = "default_service"
	DefaultHost    = "127.0.0.1"
)

// XLogger 对Logger的实现，主要提供服务的组件
type XLogger struct {
	Writers   map[string]XlogWriter
	Service   string
	Host      string
	Lock      sync.RWMutex
	entryPool *sync.Pool
	//mu        sync.Mutex
	//IsEnvDev  bool // 只有用soa初始化为开发环境才会为true
}

func (l *XLogger) With(key, value string) Entry {
	e := l.getEntry().With(key, value)
	e.(*entry).skip -= 1
	return e
}

// NewLogger 创建一个默认配置的XLogger
func NewLogger() *XLogger {
	writers := make(map[string]XlogWriter, 2)
	l := &XLogger{
		Writers:   writers,
		Service:   DefaultService,
		Host:      DefaultHost,
		entryPool: new(sync.Pool),
	}
	if host, err := os.Hostname(); err == nil {
		l.Host = host
	} else {
		logPanic("NewLogger, host not found err=%s", err.Error())
	}
	return l
}

// NewLoggerWithConsole 创建一个默认配置的XLogger，包含一个控制台writer
func NewLoggerWithConsole() *XLogger {
	l := NewLogger()
	l.Writers["console"] = NewConsoleWriter(ConsoleConfig{})
	return l
}

func (l *XLogger) getEntry() *entry {
	return newEntry(l)
}

// Warn warn级别的log
func (l *XLogger) Warn(ctx context.Context, format string, a ...interface{}) error {
	return l.getEntry().Warn(ctx, format, a...)
}

// Error error级别的log
func (l *XLogger) Error(ctx context.Context, format string, a ...interface{}) error {
	return l.getEntry().Error(ctx, format, a...)
}

// Fatal fatal级别的log
func (l *XLogger) Fatal(ctx context.Context, format string, a ...interface{}) error {
	return l.getEntry().Fatal(ctx, format, a...)
}

// Notice notice级别的log
func (l *XLogger) Notice(ctx context.Context, format string, a ...interface{}) error {
	return l.getEntry().Notice(ctx, format, a...)
}

// Debug debug级别的log
func (l *XLogger) Debug(ctx context.Context, format string, a ...interface{}) error {
	return l.getEntry().Debug(ctx, format, a...)
}

// Close 关闭打开的文件
func (l *XLogger) Close() {
	l.Lock.Lock()
	defer l.Lock.Unlock()
	for _, w := range l.Writers {
		if c, ok := w.(io.Closer); ok {
			c.Close()
		}
	}
}

// OnConfigChange 配置变更回调，依赖了xconfig
func (l *XLogger) OnConfigChange(config config.Config, sections ...string) {
	writerType := config.GetStringDefault("log", "type", "")
	file, console, err := addWriterByXconfig(config)
	if err != nil {
		if len(sections) == 0 {
			l.Fatal(nil, "xlog reload error with xconfig : %+v, sections : %+v, err : %+v", config, sections, err)
			return
		}
		s, _ := config.GetSection(sections[0])
		l.Fatal(nil, "xlog reload error with xconfig.%v(log) : %+v, err : %+v", sections[0], s, err)
		return
	}
	l.Close()
	l.SetConfig(Config{
		WriterType: writerType,
		Service:    "",
		File:       file,
		Console:    console,
	})
	return
}

func (l *XLogger) SetConfig(conf Config) {
	if len(conf.Service) > 0 {
		l.Service = conf.Service
	}
	if len(conf.File.LogName) == 0 {
		conf.File.LogName = defaultLogName
	}
	if len(conf.File.LogPath) == 0 {
		conf.File.LogPath = defaultLogPath
	}
	if host, err := os.Hostname(); err == nil {
		l.Host = host
	}
	l.flushWriter()
	if len(conf.WriterType) <= 0 {
		return
	}
	l.Lock.Lock()
	defer l.Lock.Unlock()
	if strings.Contains(conf.WriterType, TypeFile) {
		l.Writers[TypeFile] = NewFileWriter(conf.File)
	}
	if strings.Contains(conf.WriterType, TypeConsole) {
		l.Writers[TypeConsole] = NewConsoleWriter(conf.Console)
	}
}

func (l *XLogger) flushWriter() {
	l.Lock.Lock()
	for k := range l.Writers {
		delete(l.Writers, k)
	}
	l.Lock.Unlock()
}
