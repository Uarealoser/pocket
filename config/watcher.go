package config

import "sync"

type ConfigNotify struct {
	Watcher
	Config
	wg *sync.WaitGroup
}

type Watcher interface {
	Watch(notifyChan chan *ConfigNotify)
	Priority() int
	SetPriority(priority int)
	SetWatchErrorFunc(errorFunc WatchErrorFunc)
	Stop()
}

// WatchErrorFunc 监听出现错误的函数,用户自己定义
type WatchErrorFunc func(err error)

// OnConfigChangeFunc 配置更新后,用于重新加载资源,用户可以自定义
type OnConfigChangeFunc func(Config) error

type Listener interface {
	OnConfigChange(cfg Config, sections ...string)
}

type listenerManager struct {
	Sections []string
	Listener
}
