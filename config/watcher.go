package config

import "fmt"

type Watcher interface {
	Watch(notifyChan chan *ConfigNotify)
	GetPriority() int
	SetPriority(priority int)
	SetWatchErrorFunc(errorFunc WatchErrorFunc)
	Stop()
}

// WatchErrorFunc 监听出现错误的函数,用户自己定义
type WatchErrorFunc func(err error)

// 默认的监听错误报错
func DefaultWatchError(err error) {
	fmt.Println("xconfig watch error:", err)
	return
}
