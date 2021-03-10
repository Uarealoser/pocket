package config

import (
	"Uarealoser/pocket/config"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWatch_File(t *testing.T) {
	cm := config.NewconfigManager()
	defer cm.Close()
	// 每10s 执行一次，有更改则更新配置
	err := cm.WatchFile("./test.ini")
	if err != nil {
		t.Fatal()
	}
	cfg, err := cm.GetCfg()
	if err != nil {
		t.Fatal()
	}
	str, err := cfg.GetString("something", "int64_key")
	if err != nil {
		t.Fatal()
	}
	assert.Equal(t, str, "789541")
}

func TestListener_File(t *testing.T) {
	cm := config.NewconfigManager()
	defer cm.Close()
	// 每10s 执行一次，有更改则更新配置
	err := cm.WatchFile("./test.ini")
	if err != nil {
		t.Fatal()
	}
	var l config.ListenerFunc
	l = func(cfg config.Config, s ...string) {
		str, err := cfg.GetString("something", "int64_key")
		if err != nil {
			t.Fatal()
		}
		fmt.Println("str:", str)
	}
	cm.RegistListener(l, "something")
	time.Sleep(5 * time.Second)
}
