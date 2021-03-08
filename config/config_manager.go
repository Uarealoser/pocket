package config

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type ConfigManager struct {
	config atomic.Value

	// 事件监听
	watchers     []Watcher
	watchersLock sync.RWMutex

	sectionPriority map[string]int

	notifyChan chan *ConfigNotify

	WatchError WatchErrorFunc

	// 监听到配置后，出发更新
	listenerList []*listenerManager
}

func NewConfigManager() (cm *ConfigManager) {
	cm = &ConfigManager{
		watchers:        make([]Watcher, 0),
		notifyChan:      make(chan *ConfigNotify, 1),
		sectionPriority: make(map[string]int),
	}

}

// 监听事件
func (cm *ConfigManager) watchJob() {
	for {
		// 监听到配置变更
		notify, ok := <-cm.notifyChan
		if !ok {
			return
		}

		// 获取旧配置
		oldConfig, _ := cm.GetCfg()

	}
}

func (cm *ConfigManager) GetCfg() (cfg Config, err error) {
	switch t := cm.config.Load().(type) {
	case Config:
		cfg = t
		return
	}
	err = fmt.Errorf("config is not initialized")
	return
}
