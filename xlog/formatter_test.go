package xlog

import (
	"testing"
	"time"
)

func TestDefaultFormat(t *testing.T) {
	f := &DefaultFormatter{
		hasRuntime: true,
		isJSON:     true,
	}

	l := NewLoggerWithConsole()
	e := newEntry(l)
	e.skip = 2
	e.timestamp, _ = time.Parse("2006-01-02 15:04:05", "0")
	e.logger.Service = "TestService"
	e.level = NoticeLevel
	e.logID = "123123123123"
	e.message = "it is a test log"
	//e.WithString("user_id", "123")
	//e.WithString("order_id", "abcd efgh")
	//m := make(map[string]string)
	//m["mapa"] = "a"
	//m["mapb"] = "b"
	//e.WithMap(m)
	r, err := f.Format(*e)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(r))
}
