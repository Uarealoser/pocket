package benchmark

import (
	"context"
	"fmt"
	"os"
	"testing"

	"micode.be.xiaomi.com/kangaroo/pkg/xlog"

	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func BenchmarkNormal(b *testing.B) {

	fileConf := xlog.FileConfig{
		LogPath: "/tmp/xlogtest",
		LogName: "xlog_exported",
		Level:   xlog.NoticeLevel,
	}
	cConfig := xlog.ConsoleConfig{
		Level: xlog.DebugLevel,
	}

	xlog.InitWithConfig(xlog.Config{
		Service:    "xlog_exported",
		WriterType: "file",
		File:       fileConf,
		Console:    cConfig,
	})

	b.Log("start benchmark normal")
	b.Run("xlog-ctx-normal", func(b *testing.B) {
		ctx := context.WithValue(context.Background(), xlog.LOGID_TAG, "1234567890123")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				xlog.Notice(ctx, "test banchmark %s", "xlog")
			}
		})
	})

	b.Run("xlog-ctx-runtime", func(b *testing.B) {
		ctx := context.WithValue(context.Background(), xlog.LOGID_TAG, "1234567890123")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				xlog.Fatal(ctx, "test banchmark %s", "xlog")
			}
		})
	})

	b.Run("xlog-NopLog", func(b *testing.B) {
		nopLog := xlog.NewNopLogger()
		ctx := context.WithValue(context.Background(), xlog.LOGID_TAG, "1234567890123")
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				nopLog.Notice(ctx, "test banchmark %s", "xlog")
			}
		})
	})
	fileOut, err := os.OpenFile("/tmp/xlogtest/Logrus.file", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer fileOut.Close()
	var logr = logrus.New()
	logr.Out = fileOut
	logr.Level = logrus.TraceLevel
	logr.Formatter = &logrus.TextFormatter{DisableColors: true}
	//logrus.SetOutput(fileOut)
	b.Run("Logrus-normal", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logr.Debug(fmt.Sprintf("test banchmark %s", "Logrus"))
			}
		})
	})
	fileOut1, err := os.OpenFile("/tmp/xlogtest/zerolog.file", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer fileOut1.Close()
	b.Run("zerolog-normal", func(b *testing.B) {
		log.Logger = log.Output(fileOut1)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				log.Debug().Msg(fmt.Sprintf("test banchmark %s", "zerolog"))
			}
		})
	})
	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(zap.DebugLevel)
	cfg.OutputPaths = []string{
		"/tmp/xlogtest/zap.file",
	}
	logger, err := cfg.Build()
	sLogger := logger.Sugar()
	if err != nil {
		b.Fatal(err)
	}
	b.Run("uberZap-runtime", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Debug(fmt.Sprintf("test banchmark %s", "zerolog"))
			}
		})
	})

	b.Run("sugarZap-runtime", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sLogger.Debug(fmt.Sprintf("test banchmark %s", "zerolog"))
			}
		})
	})
}
