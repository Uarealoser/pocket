package xlog

import (
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
)

func TestLoggerFunc(t *testing.T) {
	outputs := make([]interface{}, 0)
	lf := LoggerFunc(func(a ...interface{}) error {
		outputs = append(outputs, a...)
		return nil
	})
	lf.Error(nil, "", 1, "error", 2.3)
	//lf.WithString("with", "string").Warn(nil, "", "warn")
	lf.Debug(nil, "", "debug")
	lf.Fatal(nil, "", "fatal")
	lf.Notice(nil, "", "notice")
	re := fmt.Sprintf("%+v", outputs)
	assert.Equal(t, "[1 error 2.3 debug fatal notice]", re)
}
