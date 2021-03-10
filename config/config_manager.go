package config

import (
	"Uarealoser/pocket/config/ini"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type ConfigManagerInterface interface {
	AddWatcher(w Watcher, priority ...int)
	GetCfg() (cfg Config, err error)
	WatchFile(path string, priority ...int) error
	SetWatchError(w WatchErrorFunc) (err error)
	RegistListener(l Listener, sections ...string)
	RemoveConfigListener(l Listener) bool
	Merge(config Config, watcher Watcher)
	Close()
}

type configManager struct {
	config atomic.Value

	// 事件监听
	Watchers     []Watcher
	watchersLock sync.RWMutex

	sectionPriority map[string]int

	notifyChan chan *ConfigNotify

	WatchError WatchErrorFunc

	// 监听到配置后，发出更新
	listenerList []*listenerManager
}

func NewconfigManager() ConfigManagerInterface {
	cm := &configManager{
		Watchers:        make([]Watcher, 0),
		notifyChan:      make(chan *ConfigNotify, 1),
		sectionPriority: make(map[string]int),
	}
	go cm.watch()
	return cm
}

// 监听事件
func (cm *configManager) watch() {
	for {
		notify, ok := <-cm.notifyChan
		if !ok {
			return
		}

		oldConfig, err := cm.GetCfg()
		cm.Merge(notify.Config, notify.Watcher)
		newConfig, _ := cm.GetCfg()
		// 配置全局listtenner(plugin.OnAfterRealod()),不传入sections
		var rootListeners = make([]*listenerManager, 0)

		for _, listener := range cm.listenerList {
			if len(listener.Sections) == 0 {
				rootListeners = append(rootListeners, listener)
				continue
			}
			// 对比是否存在配置改动
			if err == nil && compareSections(oldConfig, newConfig, listener.Sections) {
				continue
			}

			// 如果section有改动,那么执行配置变更
			listener.OnConfigChange(newConfig, listener.Sections...)
		}
		// 触发全局reload
		for _, rl := range rootListeners {
			rl.OnConfigChange(newConfig, rl.Sections...)
		}

		if notify.Wg != nil {
			notify.Wg.Done()
		}
	}
}

func (cm *configManager) GetCfg() (cfg Config, err error) {
	switch t := cm.config.Load().(type) {
	case Config:
		cfg = t
		return
	}
	err = fmt.Errorf("config is not initialized")
	return
}

// WatchFile 监听文件,文件修改后更新配置
func (cm *configManager) WatchFile(path string, priority ...int) error {

	// 注册
	watcher := &FileWatcher{
		Path:    path,
		Stopped: false,
	}
	cm.AddWatcher(watcher, priority...)

	if err := watcher.Init(cm.notifyChan); err != nil {
		return err
	}
	go watcher.Watch(cm.notifyChan)
	return nil
}

func (cm *configManager) AddWatcher(w Watcher, priority ...int) {
	cm.watchersLock.Lock()
	defer cm.watchersLock.Unlock()
	if priority == nil || len(priority) == 0 {
		priority = append(priority, len(cm.Watchers))
	}
	w.SetPriority(priority[0])
	if we := cm.WatchError; we != nil {
		w.SetWatchErrorFunc(we)
	}
	cm.Watchers = append(cm.Watchers, w)
}

func (cm *configManager) SetWatchError(w WatchErrorFunc) (err error) {
	cm.WatchError = w
	cm.watchersLock.RLock()
	defer cm.watchersLock.RUnlock()
	for _, watcher := range cm.Watchers {
		watcher.SetWatchErrorFunc(w)
	}
	return
}

// RegistListener 注册listener,当配置更新的时候,调用此函数进行更新
func (cm *configManager) RegistListener(l Listener, sections ...string) {

	lm := &listenerManager{
		Listener: l,
		Sections: sections,
	}
	cm.listenerList = append(cm.listenerList, lm)

	return
}

// RemoveConfigListener 删除注册listener
func (cm *configManager) RemoveConfigListener(l Listener) bool {

	for i, v := range cm.listenerList {
		if l == v.Listener {
			cm.listenerList = append(cm.listenerList[:i], cm.listenerList[i+1:]...)
			return true
		}
	}

	return false
}

func (cm *configManager) Merge(config Config, watcher Watcher) {
	priority := 0
	if watcher != nil {
		priority = watcher.GetPriority()
	}
	newConfig := ini.NewIniConfig()
	currentConfig, err := cm.GetCfg()
	if err == nil {
		for _, section := range currentConfig.GetSections() {
			newConfig.Sections[section], _ = currentConfig.GetSection(section)
		}
	}

	for _, section := range config.GetSections() {
		oldPriority, ok := cm.sectionPriority[section]
		if !ok || priority <= oldPriority {
			newConfig.Sections[section], _ = config.GetSection(section)
			cm.sectionPriority[section] = priority
		}
	}
	cm.config.Store(newConfig)
}

func (cm *configManager) Close() {
	for _, watcher := range cm.Watchers {
		if watcher != nil {
			watcher.Stop()
		}
	}
}

/*
	以下为Config接口实现方法
*/
// GetValue _
func (cm *configManager) GetValue(section, key string) (value interface{}, err error) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetValue(section, key)

}

// GetString _
func (cm *configManager) GetString(section, key string) (value string, err error) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetString(section, key)
}

// GetStringDefault _
func (cm *configManager) GetStringDefault(section, key string, defval string) (value string) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetStringDefault(section, key, defval)
}

// GetInt _
func (cm *configManager) GetInt(section, key string) (value int, err error) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetInt(section, key)
}

// GetIntDefault _
func (cm *configManager) GetIntDefault(section, key string, defval int) (value int) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetIntDefault(section, key, defval)
}

// GetInt64 _
func (cm *configManager) GetInt64(section, key string) (value int64, err error) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetInt64(section, key)
}

// GetInt64Default _
func (cm *configManager) GetInt64Default(section, key string, defval int64) (value int64) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetInt64Default(section, key, defval)
}

// GetSlice _
func (cm *configManager) GetSlice(section, key string) (value []string, err error) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetSlice(section, key)
}

// GetSliceDefault _
func (cm *configManager) GetSliceDefault(section, key string, defval []string) (value []string) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return defval
	}
	return xconfig.GetSliceDefault(section, key, defval)
}

// GetMap _
func (cm *configManager) GetMap(section, key string) (value map[string]string, err error) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetMap(section, key)
}

// GetMapDefault _
func (cm *configManager) GetMapDefault(section, key string, defval map[string]string) (value map[string]string) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return defval
	}
	return xconfig.GetMapDefault(section, key, defval)
}

// GetBool _
func (cm *configManager) GetBool(section, key string) (value bool, err error) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetBool(section, key)
}

// GetBoolDefault _
func (cm *configManager) GetBoolDefault(section, key string, defval bool) (value bool) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return defval
	}
	return xconfig.GetBoolDefault(section, key, defval)
}

// GetFloat64 _
func (cm *configManager) GetFloat64(section, key string) (value float64, err error) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetFloat64(section, key)
}

// GetFloat64Default _
func (cm *configManager) GetFloat64Default(section, key string, defval float64) (value float64) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return defval
	}
	return xconfig.GetFloat64Default(section, key, defval)
}

// GetDuration _
func (cm *configManager) GetDuration(section, key string) (value time.Duration, err error) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetDuration(section, key)
}

// GetDurationDefault _
func (cm *configManager) GetDurationDefault(section, key string, defval time.Duration) (value time.Duration) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return defval
	}
	return xconfig.GetDurationDefault(section, key, defval)
}

// GetSectionKeys _
func (cm *configManager) GetSectionKeys(section string) (keys []string) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetSectionKeys(section)
}

// GetSection _
func (cm *configManager) GetSection(section string) (sec map[string]interface{}, err error) {

	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetSection(section)
}

func (cm *configManager) GetSections() (sections []string) {
	xconfig, err := cm.GetCfg()
	if err != nil {
		return
	}
	return xconfig.GetSections()
}
