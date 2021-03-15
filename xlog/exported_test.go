package xlog

import (
	"context"
	"testing"
)

func TestExported(t *testing.T) {
	Notice(context.Background(), "it is a exported %s log again", NoticeLevel)

	Notice(context.Background(), "it is a exported %s log third", NoticeLevel)

	Error(context.Background(), "it is a exported %s log third", ErrorLevel)

	data := make(map[string]string, 0)
	data["method"] = "efgh"
	data["user_id"] = "666"
	data["a"] = "b"

	fileConf := FileConfig{
		LogPath: "/tmp/xlogtest",
		LogName: "xlog_exported",
		Level:   NoticeLevel,
	}
	cConfig := ConsoleConfig{
		Level: DebugLevel,
	}
	InitWithConfig(Config{
		Service:    "xlog_exported",
		WriterType: "file|console",
		File:       fileConf,
		Console:    cConfig,
	})

	var ctx context.Context
	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890122")
	Debug(ctx, "it is a exported %s log", DebugLevel)

	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890123")
	Notice(ctx, "it is a exported %s log", NoticeLevel)

	//ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890120")
	//WithString("user_id", "12345").WithString("order_id", "xxxx").Notice(ctx, "it is a exported %s log", NoticeLevel)
	//
	//ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890120")
	//WithMap(data).Notice(ctx, "it is a exported %s log", NoticeLevel)

	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890121")
	Warn(ctx, "it is a exported %s log", WarnLevel)

	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890126")
	Fatal(ctx, "it is a exported %s log", FatalLevel)

	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890125")
	Notice(ctx, "it is a exported %s log again", NoticeLevel)

	Notice(ctx, "it is a exported %s log again", NoticeLevel)

	Notice(ctx, "it is a exported %s log third", NoticeLevel)

	With("asdf", "asdfasdf").Fatal(ctx, "asdf%v", "asdf")
	l := GetDefaultLogger()
	l.Fatal(ctx, "it is a exported %s log", FatalLevel)
	l.With("asdf", "asdfasdf").Fatal(ctx, "it is a exported %s log", FatalLevel)
	//config := make(map[string]string)
	//config["path"] = "./logs"
	//config["filename"] = "xlog_old"
	//config["level"] = "debug"
	//config["service"] = "xlog_old"
	//xlogold.InitLogger("file", config)
	//xlogold.Noticex("1234567890123", "test xlogold debug %s", "xlog")
	//xlogold.Debugx("1234567890123", "test xlogold debug %s", "xlog")
	//xlogold.Fatalx("1234567890123", "test xlogold debug %s", "xlog")
}

func TestExportedJson(t *testing.T) {
	Notice(context.Background(), "it is a exported %s log again", NoticeLevel)

	Notice(context.Background(), "it is a exported %s log third", NoticeLevel)

	Error(context.Background(), "it is a exported %s log third", ErrorLevel)

	data := make(map[string]string, 0)
	data["method"] = "efgh"
	data["user_id"] = "666"
	data["a"] = "b"

	fileConf := FileConfig{
		LogPath: "/tmp/xlogtest",
		LogName: "xlog_exported",
		Level:   NoticeLevel,
		IsJson:  true,
	}
	cConfig := ConsoleConfig{
		Level:  DebugLevel,
		IsJson: true,
	}
	InitWithConfig(Config{
		Service:    "xlog_exported",
		WriterType: "file|console",
		File:       fileConf,
		Console:    cConfig,
	})

	var ctx context.Context
	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890122")
	Debug(ctx, "it is a exported %s log", DebugLevel)

	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890123")
	Notice(ctx, "it is a exported %s log", NoticeLevel)

	//ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890120")
	//WithString("user_id", "12345").WithString("order_id", "xxxx").Notice(ctx, "it is a exported %s log", NoticeLevel)
	//
	//ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890120")
	//WithMap(data).Notice(ctx, "it is a exported %s log", NoticeLevel)

	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890121")
	Warn(ctx, "it is a exported %s log", WarnLevel)

	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890126")
	Fatal(ctx, "it is a exported %s log", FatalLevel)

	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890125")
	Notice(ctx, "it is a exported %s log again", NoticeLevel)

	Notice(ctx, "it is a exported %s log again", NoticeLevel)

	Notice(ctx, "it is a exported %s log third", NoticeLevel)

	With("asdf", "asdfasdf").Fatal(ctx, "asdf%v", "asdf")
	l := GetDefaultLogger()
	l.Fatal(ctx, "it is a exported %s log", FatalLevel)
	l.With("asdf", "asdfasdf").Fatal(ctx, "it is a exported %s log", FatalLevel)
	//config := make(map[string]string)
	//config["path"] = "./logs"
	//config["filename"] = "xlog_old"
	//config["level"] = "debug"
	//config["service"] = "xlog_old"
	//xlogold.InitLogger("file", config)
	//xlogold.Noticex("1234567890123", "test xlogold debug %s", "xlog")
	//xlogold.Debugx("1234567890123", "test xlogold debug %s", "xlog")
	//xlogold.Fatalx("1234567890123", "test xlogold debug %s", "xlog")
}

//func ExampleNotice() {
//	writersConf := make([]WriterConfig, 2)
//	writersConf = append(writersConf, WriterConfig{
//		Name:    "file",
//		LogPath: "./logs",
//		LogName: "xlog_exported",
//		Level:   NoticeLevel,
//	})
//	writersConf = append(writersConf, WriterConfig{
//		Name:  "console",
//		Level: DebugLevel,
//		FormatterConf: FormatterConfig{
//			Name:       "console",
//			HasRuntime: true,
//			IsJSON:     true,
//		},
//	})
//	InitWithConfig(Config{
//		Service: "xlog_test_service",
//		Writers: writersConf,
//	})
//
//	var ctx context.Context
//
//	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890123")
//	Ctx(ctx).Cause("xlog.cause").Notice("it is a exported %s log", NoticeLevel)
//
//	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890120")
//	Ctx(ctx).Cause("xlog.cause").WithInt64("user_id", 12345).WithString("order_id", "xxxx").Notice("it is a exported %s log", NoticeLevel)
//
//	ctx = context.WithValue(context.Background(), LOGID_TAG, "1234567890125")
//	Ctx(ctx).Notice("it is a exported %s log again", NoticeLevel)
//
//	Cause("xlog.cause").Notice("it is a exported %s log again", NoticeLevel)
//
//	Notice("it is a exported %s log third", NoticeLevel)
//}
