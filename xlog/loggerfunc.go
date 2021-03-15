package xlog

import (
	"context"
)

// LoggerFunc 实现go-kit 自带的方法，传入log的处理方法，生成一个自定义输出的logger，在此情况下前两个参数(ctx,format)作废
type LoggerFunc func(...interface{}) error

// Log implements Logger by calling f(keyvals...). 保留原方法
func (f LoggerFunc) Log(keyvals ...interface{}) error {
	return f(keyvals...)
}

// Warn 过滤前两个参数
func (f LoggerFunc) Warn(_ context.Context, _ string, a ...interface{}) error {
	return f(a...)
}

func (f LoggerFunc) Error(_ context.Context, _ string, a ...interface{}) error {
	return f(a...)
}

func (f LoggerFunc) Fatal(_ context.Context, _ string, a ...interface{}) error {
	return f(a...)
}

func (f LoggerFunc) Notice(ctx context.Context, format string, a ...interface{}) error {
	return f(a...)
}

func (f LoggerFunc) Debug(ctx context.Context, format string, a ...interface{}) error {
	return f(a...)
}

func (f LoggerFunc) With(key, value string) Entry {
	return f
}

// Close 关闭打开的文件
func (f LoggerFunc) Close() {}

func (f LoggerFunc) SetConfig(conf Config) {}
