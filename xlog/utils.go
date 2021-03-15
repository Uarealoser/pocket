package xlog

import (
	"Uarealoser/pocket/xlog/rawjson"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

// GenLogID 生成LogID
func GenLogID() int64 {
	rand.Seed(time.Now().UnixNano())
	num := rand.Int63n(899999999) + 100000000
	return num
}

func tagsFormat(buffer *bytes.Buffer, tags ...string) {
	for _, tag := range tags {
		buffer.WriteByte('[')
		buffer.WriteString(tag)
		buffer.WriteByte(']')
		buffer.WriteByte(' ')
	}
}
func getRuntimeInfo(skip int) (function, filename string, lineno int) {
	function = "???"
	pc, filename, lineno, ok := runtime.Caller(skip)
	if ok {
		function = runtime.FuncForPC(pc).Name()
	}
	return
}
func runtimeFormat(buffer *bytes.Buffer, function, filename string, lineno int) {
	buffer.WriteString(filepath.Base(function))
	buffer.WriteByte(':')
	buffer.WriteString(filepath.Base(filename))
	buffer.WriteByte(':')
	buffer.WriteString(strconv.FormatInt(int64(lineno), 10))
}
func logPanic(format string, args ...interface{}) {
	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)

	tags := []string{
		time.Now().Format("2006-01-02 15:04:05"),
		panicLevel.String(),
		"xlog.panic",
		DEFAULT_LOG_ID,
	}
	tagsFormat(buffer, tags...)

	function, filename, lineno := getRuntimeInfo(3)
	runtimeFormat(buffer, function, filename, lineno)

	buffer.WriteString(fmt.Sprintf(format, args...))

	os.Stderr.Write(buffer.Bytes())
}
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func fieldsNormalFormat(buffer *bytes.Buffer, message string, data map[string]string) {
	if len(data) > 0 {
		buffer.WriteByte('[')
		i := 0
		for k, v := range data {
			if i != 0 {
				buffer.WriteByte(',')
			}
			buffer.WriteString(k)
			buffer.WriteByte(':')
			buffer.WriteString(v)
			i++
		}
		buffer.WriteByte(']')
		buffer.WriteByte(' ')
	}
	buffer.WriteString(message)
}

func fieldsJsonFormat(buffer *bytes.Buffer, data map[string]string) {
	tmp := make([]byte, 0, 500)
	tmp = append(tmp, '{')
	for key, v := range data {
		tmp = rawjson.AppendKey(tmp, key)
		tmp = rawjson.AppendString(tmp, v)
	}
	tmp = append(tmp, '}')
	buffer.Write(tmp)
}
