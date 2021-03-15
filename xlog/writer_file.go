package xlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileWriter struct {
	name      string
	filepath  string
	filename  string
	level     Level
	file      io.WriteCloser //io.Writer
	errFile   io.WriteCloser //io.Writer
	formatter Formatter
}

var (
	defaultLogPath = "/home/work/logs/applogs"
	defaultLogName = "app"
)

// NewFileWriter 初始化失败panic
func NewFileWriter(conf FileConfig) XlogWriter {
	w := &FileWriter{
		name:      "file",
		filepath:  defaultLogPath,
		filename:  defaultLogName,
		level:     conf.Level,
		formatter: newFormatter(FormatterConfig{IsJSON: conf.IsJson}),
	}
	if conf.LogPath != "" {
		w.filepath = conf.LogPath
	}

	var err error
	pExist, err := pathExists(w.filepath)
	if err != nil {
		panic(fmt.Sprintf("err check filepath[%s] exists failed err=%s", w.filepath, err.Error()))
	}
	if !pExist {
		if err := os.Mkdir(w.filepath, 0755); err != nil {
			panic(fmt.Sprintf("err make filepath[%s] failed err=%s", w.filepath, err.Error()))
		}
	}

	if conf.LogName != "" {
		w.filename = conf.LogName
	}
	logFullPath := filepath.Join(w.filepath, w.filename+".log")
	w.file, err = w.openFile(logFullPath)
	if err != nil {
		panic(fmt.Sprintf("log file[%s] open failed err=%s", logFullPath, err.Error()))
	}
	errlogFullPath := logFullPath + ".wf"
	w.errFile, err = w.openFile(errlogFullPath)
	if err != nil {
		panic(fmt.Sprintf("err log file[%s] open failed err=%s", errlogFullPath, err.Error()))
	}
	return w
}

func (w *FileWriter) openFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
}

func (w *FileWriter) Name() string {
	return w.name
}

func (w *FileWriter) Level() Level {
	return w.level
}

func (w *FileWriter) Output(e entry) error {
	file := w.file
	if e.level >= WarnLevel {
		file = w.errFile
	}
	byteMsg, err := w.formatter.Format(e)
	if err != nil {
		return err
	}
	file.Write(byteMsg)
	return nil
}

func (w *FileWriter) Close() error {
	if w.errFile != nil {
		w.errFile.Close()
	}
	if w.file != nil {
		w.file.Close()
	}
	return nil
}
