package xlog

import (
	"context"
	"testing"
)

func TestLogger(t *testing.T) {
	service := "xlog_test_service"

	log := NewLoggerWithConsole()
	conf := FileConfig{
		LogPath: "/tmp/xlogtest",
		LogName: service,
		Level:   DebugLevel,
	}
	log.Writers["file"] = NewFileWriter(conf)

	ctx := context.WithValue(context.Background(), LOGID_TAG, "1234567890123")
	if err := Debug(ctx, "test ctx"); err != nil {
		t.Errorf("logger ctx debug err:%v", err)
	}
	if err := log.Notice(ctx, "test cause"); err != nil {
		t.Errorf("logger ctx notice err:%v", err)
	}

	data := make(map[string]string, 0)
	data["method"] = "efgh"
	data["user_id"] = "666"
	data["a"] = "b"
	//if err := log.WithString("method", "abcd").WithString("user_id", "1234").Notice(ctx, "test cause"); err != nil {
	//	t.Errorf("logger ctx notice err:%v", err)
	//}
	//if err := log.WithMap(data).Notice(ctx, "test cause"); err != nil {
	//	t.Errorf("logger ctx notice err:%v", err)
	//}
	if err := log.Warn(ctx, "test warn"); err != nil {
		t.Errorf("logger ctx warn err:%v", err)
	}
	if err := log.Fatal(ctx, "test fatal"); err != nil {
		t.Errorf("logger ctx fatal err:%v", err)
	}

	// e := l.getEntry()
	// e.Ctx(ctx)
	// e.Cause("xlog.test")
	// // e.Skip(2)
	// if err := e.Notice("it is a test entry notice: %s %d", "string", 1); err != nil {
	// 	t.Errorf("entry notice err: %v", err)
	// }
	// if err := e.Error("it is a test entry error: %s %d", "string", 1); err != nil {
	// 	t.Errorf("entry error err: %v", err)
	// }
}
