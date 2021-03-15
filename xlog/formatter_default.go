package xlog

import (
	"bytes"
	"strings"
)

type DefaultFormatter struct {
	hasRuntime bool
	isJSON     bool
}

func NewDefaultFormatter(conf FormatterConfig) Formatter {
	return &DefaultFormatter{
		hasRuntime: conf.HasRuntime,
		isJSON:     conf.IsJSON,
	}
}

func (f *DefaultFormatter) Format(e entry) ([]byte, error) {
	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)

	tags := []string{
		e.timestamp.Format("2006-01-02 15:04:05"),
		e.logger.Service,
		e.logger.Host,
		e.level.String(),
		e.logID,
	}
	tagsFormat(buffer, tags...)

	if f.hasRuntime || e.level >= FatalLevel {
		buffer.WriteByte('[')
		function, filename, lineno := getRuntimeInfo(e.skip)
		runtimeFormat(buffer, function, filename, lineno)
		buffer.WriteByte(']')
	}
	buffer.WriteByte(' ')
	if f.isJSON {
		e.With(strings.ToLower(e.level.String()), e.message)
		fieldsJsonFormat(buffer, e.data)
	} else if len(e.data) > 0 {
		e.With(strings.ToLower(e.level.String()), e.message)
		fieldsJsonFormat(buffer, e.data)
	} else {
		fieldsNormalFormat(buffer, e.message, e.data)
	}
	buffer.WriteByte('\n')

	return buffer.Bytes(), nil
}
