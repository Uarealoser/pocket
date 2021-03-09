package config

type Listener interface {
	OnConfigChange(cfg Config, sections ...string)
}

type listenerManager struct {
	Listener
	Sections []string
}

// OnConfigChangeFunc 配置更新后,用于重新加载资源,用户可以自定义
type OnConfigChangeFunc func(Config) error

type ListenerFunc func(Config, ...string)

func (l ListenerFunc) OnConfigChange(cfg Config, sections ...string) {
	l(cfg, sections...)
}
