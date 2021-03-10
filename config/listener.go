package config

type Listener interface {
	OnConfigChange(cfg Config, sections ...string)
}

type listenerManager struct {
	Listener
	Sections []string // 以section为粒度，进行监听
}

// OnConfigChangeFunc 配置更新后,用于重新加载资源,用户可以自定义
type OnConfigChangeFunc func(Config) error

type ListenerFunc func(Config, ...string)

func (l ListenerFunc) OnConfigChange(cfg Config, sections ...string) {
	l(cfg, sections...)
}
