package xlog

import (
	"bytes"
	"sync"
)

var (
	bufferPool *sync.Pool
)

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}
func newFormatter(conf FormatterConfig) Formatter {
	return NewDefaultFormatter(conf)
}

// formatter
type Formatter interface {
	Format(entry) ([]byte, error)
}
