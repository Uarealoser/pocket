package xlog

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestConsoleWriter(t *testing.T) {
	norWriter := new(bytes.Buffer)
	errWriter := new(bytes.Buffer)
	w := &ConsoleWriter{
		name:      "console",
		stdout:    norWriter,
		stderr:    errWriter,
		level:     DebugLevel,
		formatter: &DefaultFormatter{},
	}
	wantLevel := DebugLevel
	if w.Level() != wantLevel {
		t.Errorf("consoleWriter level expected: %s, get: %s", wantLevel, w.Level())
	}

	l := NewLoggerWithConsole()
	e := l.getEntry()
	e.message = "test"
	if err := w.Output(*e); err != nil {
		t.Error(err)
	}
	got := string(norWriter.Bytes())
	want := fmt.Sprintf("[%s] [%s] [%s] [%s] [%s] [] %s\n", e.timestamp.Format("2006-01-02 15:04:05"), e.logger.Service, e.logger.Host, e.level.String(), e.logID, e.message)
	if got != want {
		t.Errorf("got:'%s',want:'%s'", got, want)
	}

	// test json
	norWriter.Reset()
	w.formatter = &DefaultFormatter{isJSON: true}
	e.level = NoticeLevel
	//e.WithString("user_id", "123")
	//e.WithString("order_id", "xxx")
	//e.WithString("price", "12.3")
	if err := w.Output(*e); err != nil {
		t.Error(err)
	}
	got = string(norWriter.Bytes())
	fmt.Println(got)
	// jsonMsg := `{"user_id":123,"order_id":"xxx","price":12.3,"notice":"test"}`
	// want = fmt.Sprintf("[%s] [%s] [%s] [%s] [%s] [%s] %s\n", e.timestamp.Format("2006-01-02 15:04:05"), e.logger.Service, e.logger.Host, e.level.String(), e.cause, e.logID, jsonMsg)
	// if got != want {
	// 	t.Errorf("got:'%s',want:'%s'", got, want)
	// }
}

type myWriteCloser struct {
	io.Writer
}

func (myWriteCloser) Close() error {
	return nil
}
func TestFileWriter(t *testing.T) {
	norWriter := new(bytes.Buffer)
	errWriter := new(bytes.Buffer)
	w := &FileWriter{
		name:      "file",
		filepath:  defaultLogPath,
		filename:  defaultLogName,
		file:      myWriteCloser{norWriter},
		errFile:   myWriteCloser{errWriter},
		formatter: &DefaultFormatter{},
	}

	expectedName := "file"
	if w.Name() != expectedName {
		t.Errorf("fileWriter name expected: %s, get: %s", expectedName, w.Name())
	}

	l := NewLoggerWithConsole()
	e := l.getEntry()
	e.message = "it is test log."
	if err := w.Output(*e); err != nil {
		t.Error(err)
	}
	got := string(norWriter.Bytes())
	want := fmt.Sprintf("[%s] [%s] [%s] [%s] [%s] [] %s\n", e.timestamp.Format("2006-01-02 15:04:05"), e.logger.Service, e.logger.Host, e.level.String(), e.logID, e.message)
	if got != want {
		t.Errorf("got:'%s',want:'%s'", got, want)
	}

	e.level = WarnLevel
	if err := w.Output(*e); err != nil {
		t.Error(err)
	}
	got = string(errWriter.Bytes())
	want = fmt.Sprintf("[%s] [%s] [%s] [%s] [%s] [] %s\n", e.timestamp.Format("2006-01-02 15:04:05"), e.logger.Service, e.logger.Host, e.level.String(), e.logID, e.message)
	if got != want {
		t.Errorf("got:'%s',want:'%s'", got, want)
	}
}
