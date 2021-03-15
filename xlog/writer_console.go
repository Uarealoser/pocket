package xlog

import (
	"io"
	"os"
)

type ConsoleWriter struct {
	name      string
	stdout    io.Writer
	stderr    io.Writer
	level     Level
	formatter Formatter
}

func NewConsoleWriter(conf ConsoleConfig) XlogWriter {
	return &ConsoleWriter{
		stdout:    os.Stdout,
		stderr:    os.Stderr,
		level:     conf.Level,
		formatter: newFormatter(FormatterConfig{IsJSON: conf.IsJson}),
	}
}

func (w *ConsoleWriter) Name() string {
	return w.name
}

func (w *ConsoleWriter) Level() Level {
	return w.level
}

func (w *ConsoleWriter) Output(e entry) error {
	file := w.stdout
	if w.level >= WarnLevel {
		file = w.stderr
	}
	byteMsg, err := w.formatter.Format(e)
	if err != nil {
		return err
	}
	file.Write(byteMsg)
	return nil
}
