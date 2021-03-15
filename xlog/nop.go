package xlog

import (
	"context"
)

type nop struct{}

func (l *nop) With(key, value string) Entry {
	return l
}

func (l *nop) Warn(ctx context.Context, format string, a ...interface{}) error {
	return nil
}

func (l *nop) Error(ctx context.Context, format string, a ...interface{}) error {
	return nil
}

func (l *nop) Fatal(ctx context.Context, format string, a ...interface{}) error {
	return nil
}

func (l *nop) Notice(ctx context.Context, format string, a ...interface{}) error {
	return nil
}

func (l *nop) Debug(ctx context.Context, format string, a ...interface{}) error {
	return nil
}

// Close 关闭打开的文件
func (l *nop) Close() {}

func (l *nop) SetConfig(conf Config) {}
