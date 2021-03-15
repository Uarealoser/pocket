/*
Package xlog是封装了log功能的包
Logger 为定义的interface
Xlogger, nop和loggerfunc是对Logger的实现
*/
package xlog

import "context"

// Logger 接口定义
type Logger interface {
	//报警级别从高到低
	Fatal(ctx context.Context, format string, a ...interface{}) error
	Error(ctx context.Context, format string, a ...interface{}) error
	Warn(ctx context.Context, format string, a ...interface{}) error
	Notice(ctx context.Context, format string, a ...interface{}) error
	Debug(ctx context.Context, format string, a ...interface{}) error
	With(key, value string) Entry
	//file类型关闭文件句柄
	Close()

	SetConfig(conf Config)
}

// Entry 接口定义
type Entry interface {
	//报警级别从高到低
	Fatal(ctx context.Context, format string, a ...interface{}) error
	Error(ctx context.Context, format string, a ...interface{}) error
	Warn(ctx context.Context, format string, a ...interface{}) error
	Notice(ctx context.Context, format string, a ...interface{}) error
	Debug(ctx context.Context, format string, a ...interface{}) error
	With(key, value string) Entry
}
