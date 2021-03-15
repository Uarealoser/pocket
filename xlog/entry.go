package xlog

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

const (
	DEFAULT_LOG_ID = "900000001"
	LOGID_TAG      = "logID"
	DEFAULT_SKIP   = 6 // 默认是6级调用
)

// entry
type entry struct {
	logger    *XLogger
	level     Level
	timestamp time.Time
	logID     string
	message   string
	skip      int // 获取runtimeinfo时参数
	isFill    bool
	data      map[string]string
}

func newEntry(l *XLogger) *entry {
	e, ok := l.entryPool.Get().(*entry)
	if ok {
		e.reset()
		return e
	}
	return &entry{
		logger:    l,
		logID:     DEFAULT_LOG_ID,
		timestamp: time.Now(),
		skip:      DEFAULT_SKIP,
		data:      make(map[string]string),
	}
}

func (e *entry) Warn(ctx context.Context, format string, a ...interface{}) error {
	return e.output(ctx, WarnLevel, format, a...)
}

func (e *entry) Error(ctx context.Context, format string, a ...interface{}) error {
	return e.output(ctx, ErrorLevel, format, a...)
}

func (e *entry) Fatal(ctx context.Context, format string, a ...interface{}) error {
	return e.output(ctx, FatalLevel, format, a...)
}

func (e *entry) Notice(ctx context.Context, format string, a ...interface{}) error {
	return e.output(ctx, NoticeLevel, format, a...)
}

func (e *entry) Debug(ctx context.Context, format string, a ...interface{}) error {
	return e.output(ctx, DebugLevel, format, a...)
}

func (e *entry) Close() {
	e.logger.Close()
}

func (e *entry) reset() {
	e.logID = DEFAULT_LOG_ID
	e.timestamp = time.Now()
	e.message = ""
	e.isFill = false
	e.skip = DEFAULT_SKIP
	for k := range e.data {
		delete(e.data, k)
	}
}

func (e *entry) output(ctx context.Context, lv Level, format string, args ...interface{}) error {
	e.level = lv
	var err error
	// 读锁，防止打日志的时候修改日志writer
	e.logger.Lock.RLock()
	for _, w := range e.logger.Writers {
		if !e.level.IsEnabled(w.Level()) {
			continue
		}

		e.fillFields(ctx, format, args...)
		if errOut := w.Output(*e); errOut != nil {
			err = errOut
		}
	}
	e.logger.Lock.RUnlock()

	// 释放entry到pool
	e.release()
	return err
}

func (e *entry) fillFields(ctx context.Context, format string, args ...interface{}) {
	if e.isFill {
		return
	}
	if ctx != nil {
		if v := ctx.Value(LOGID_TAG); v != nil {
			if vi, ok := v.(int); ok {
				e.logID = strconv.Itoa(vi)
			} else if vi, ok := v.(int64); ok {
				e.logID = strconv.FormatInt(vi, 10)
			} else if vs, ok := v.(string); ok {
				e.logID = vs
			} else {
				e.logID = fmt.Sprintf("%v", v)
			}
		}
	}
	e.message = fmt.Sprintf(format, args...)
	e.isFill = true
}

func (e *entry) release() {
	e.logger.entryPool.Put(e)
}

func (e *entry) With(key, value string) Entry {
	e.data[key] = value
	return e
}
