package discovery

const (
	DiscoveryTypeNacos = "nacos"
)

type ServiceInfo struct {
	group   string
	service string
}

type Discovery interface {
	Open() error
	// 获取本服务订阅的所有服务
	GetSubscribe() ([]ServiceInfo, error)
	// 获取指定服务名的服务
	GetService(groupName, serviceName string) (*Service, error)
	Close()
	GetDiscoveryType() string
}

func NewDiscovery(group, service, env, rootNode, discoveryHost, auth, idc string) (discovery Discovery, err error) {
	return NewNacosDiscovery(group, service, env, auth, auth, idc)
}
