package discovery

import (
	"Uarealoser/pocket/xlog"
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net"
	"strconv"
	"strings"
)

const (
	NacosDefaultClientTimeout = 5000
	NacosDefaultLogLevel      = "dubug"
	NacosNamespace            = "nuzar_namespace"
	NacosDefaultSysConfNode   = "sysconf"
)

type NacosDiscovery struct {
	namingClient  naming_client.INamingClient
	configClient  config_client.IConfigClient
	naocsUserName string
	nacosPasswd   string

	groupName   string
	serviceName string
	env         string // 环境(onlinec/prev/test)
	idc         string // 进程所在的机房

	nacosAddr     string
	nacosPort     int
	discoveryType string
}

func NewNacosDiscovery(group, service, environ, nacosAddr, auth, idc string) (discovery *NacosDiscovery, err error) {
	dc := new(NacosDiscovery)
	var userName, password string
	//设置权限验证
	if len(auth) > 0 {
		userName, password, err = dc.GetAuth(auth)
		if err != nil {
			xlog.Warn(nil, "NewXNacosRegister failed, cause:GetAuth, err:%#v", err)
			return
		}
	}
	clientConfig := constant.ClientConfig{
		TimeoutMs:           NacosDefaultClientTimeout,
		NotLoadCacheAtStart: true,
		MaxAge:              3,
		LogLevel:            NacosDefaultLogLevel,
		Username:            userName,
		Password:            password,
		NamespaceId:         NacosNamespace,
	}
	ip, port, err := net.SplitHostPort(strings.TrimRight(nacosAddr, "/"))
	if err != nil {
		ip = strings.TrimRight(nacosAddr, "/")
		port = "8848"
	}
	iport, err := strconv.ParseUint(port, 10, 64)
	if err != nil {
		return
	}
	serverConfig := []constant.ServerConfig{
		{
			IpAddr: ip,
			Port:   iport,
			Scheme: "http",
		},
	}
	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfig,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		xlog.Error(nil, "NewXNacosRegister failed, cause:create namingClient failed, err: %#v", err)
		return
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfig,
		"clientConfig":  clientConfig,
	})

	dc.namingClient = namingClient
	dc.configClient = configClient
	dc.discoveryType = DiscoveryTypeNacos
	dc.groupName = group
	dc.serviceName = service
	dc.env = environ
	dc.idc = idc
	dc.nacosAddr = nacosAddr
	return dc, nil
}

func (nd *NacosDiscovery) Open() error {
	return nil
}

// 获取本服务订阅的所有服务
func (nd *NacosDiscovery) GetSubscribe() ([]ServiceInfo, error) {
	return nil, nil
}

// 获取指定服务名的服务
func (nd *NacosDiscovery) GetService(groupName, serviceName string) (*Service, error) {
	service, err := nd.namingClient.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
		Clusters:    []string{nd.idc},
		GroupName:   groupName,
	})
	if err != nil {
		xlog.Error(nil, "get service failed, err:%#v", err)
		return nil, err
	}
	sysConf, err := nd.configClient.GetConfig(vo.ConfigParam{
		DataId: serviceName + "." + nd.env + "." + NacosDefaultSysConfNode,
		Group:  groupName,
	})
	if err != nil {
		xlog.Error(nil, "get sysconf from nacos failed, err:%#v", err)
		return nil, err
	}
	nvservice, err := nd.initService(serviceName, groupName, sysConf)
	if err != nil {
		err = fmt.Errorf("init service failed, err:%v", err)
		return nil, err
	}
	for _, instance := range service.Hosts {
		nvservice.serverList = append(nvservice.serverList, NewServer(instance.Ip, uint32(instance.Port), instance.Weight, instance.ClusterName))
	}
	return nvservice, nil
}

func (nd *NacosDiscovery) Close() {
	return
}

func (nd *NacosDiscovery) GetDiscoveryType() string {
	return nd.discoveryType
}

func (nd *NacosDiscovery) GetAuth(auth string) (userName, password string, err error) {
	if !strings.Contains(auth, ":") {
		xlog.Warn(nil, "wrong auth format,err: %#v", err)
		return
	}
	userNameAndPwd := strings.SplitN(auth, ":", 2)
	userName = userNameAndPwd[0]
	password = userNameAndPwd[1]
	return
}

func (nd *NacosDiscovery) initService(serviceName, groupName, data string) (service *Service, err error) {
	var sysConf map[string]interface{}
	if err := json.Unmarshal([]byte(data), &sysConf); err != nil {
		xlog.Error(nil, "json unmarshal failed, data:%#v,"+
			" err:%#v", data, err)
		return
	}
	service = NewService(groupName, serviceName)
	if err = service.Init(sysConf); err != nil {
		xlog.Error(nil, "init service from sysconf failed,cause:%#v", err)
		return
	}
	return
}
