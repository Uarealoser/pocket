package discovery

import (
	"Uarealoser/pocket/xlog"
	"fmt"
)

const (
	ServiceDefaultTimeout  = 1000 // ms
	ServiceMaxIdleTimeout  = 30   // ms
	ServiceTypeApplication = "application"
)

var (
	ServiceMaxIdleConn uint32 = 4096
)

type Timeout struct {
	ConnTimeoutMs  uint32 // ms
	ReadTimeoutMs  uint32 // ms
	WriteTimeoutMs uint32 // ms
}

type Service struct {
	group      string
	service    string
	serverList []*Server
	idcServers map[string][]*Server // key:idc, value serverList

	transport string
	protocol  string
	timeout   Timeout
	balance   string
	host      string
	maxRetry  uint32

	timeoutMap     map[string]Timeout // key method
	maxIdleConn    uint32
	maxIdleTimeout uint32

	serviceType string //服务类型为application/mysql/reids
}

type RedisService struct {
	Service
	auth string
}

type MysqlService struct {
	Service
	// mysql
	userName string
	passwd   string
	database string
}

func NewService(group, service string) *Service {
	return &Service{
		group:          group,
		service:        service,
		serverList:     make([]*Server, 0),
		idcServers:     make(map[string][]*Server),
		timeoutMap:     make(map[string]Timeout),
		maxIdleConn:    ServiceMaxIdleConn,
		maxIdleTimeout: ServiceMaxIdleTimeout,
	}
}

func (s *Service) Init(sysConf map[string]interface{}) (err error) {
	global, ok := sysConf["global"]
	if !ok {
		xlog.Warn(nil, "Init [%s] failed,not found global config", s.service)
		return fmt.Errorf("not found sys conf")
	}
	globalMap, ok := global.(map[string]interface{})
	if !ok {
		xlog.Warn(nil, "Init [%s] failed,invalid sys conf[%+v]", s.service, sysConf)
		return fmt.Errorf("sys conf format error")
	}
	serviceType, err := getStringFromMap(globalMap, "sign")
	if err != nil {
		xlog.Warn(nil, "init [%s] failed, err: %v", s.service, err)
		return
	}
	s.serviceType = serviceType
	transport, err := getStringFromMap(globalMap, "transport")
	if err != nil {
		if s.serviceType == ServiceTypeApplication {
			xlog.Warn(nil, "Init [%s] failed, err:%v", s.service, err)
			return
		}
	}
	protocol, err := getStringFromMap(globalMap, "protocol")
	if err != nil {
		if s.serviceType == ServiceTypeApplication {
			xlog.Warn(nil, "Init [%s] failed, err:%v", s.service, err)
			return
		}
	}
	balance, err := getStringFromMap(globalMap, "balance")
	if err != nil {
		xlog.Warn(nil, "Init [%s] failed, err:%v", s.service, err)
		balance = "random"
	}
	connTimeoutMs, err := getIntFromMap(globalMap, "conn_timeout_ms")
	if err != nil {
		connTimeoutMs = ServiceDefaultTimeout
		xlog.Warn(nil, "connTimeoutMs %v, use default value[%v]", err, connTimeoutMs)
	}

	readTimeoutMs, err := getIntFromMap(globalMap, "read_timeout_ms")
	if err != nil {
		readTimeoutMs = ServiceDefaultTimeout
		xlog.Warn(nil, "readTimeoutMs %v, use default value[%v]", err, readTimeoutMs)
	}

	writeTimeoutMs, err := getIntFromMap(globalMap, "write_timeout_ms")
	if err != nil {
		writeTimeoutMs = ServiceDefaultTimeout
		xlog.Warn(nil, "writeTimeoutMs %v, use default value[%v]", err, writeTimeoutMs)
	}

	host, _ := getStringFromMap(globalMap, "hostname")
	maxRetry, err := getIntFromMap(globalMap, "max_retry")
	if err != nil {
		maxRetry = 0
		xlog.Warn(nil, "maxRetyr %v, use default value[%v]", err, maxRetry)
	}

	s.transport = transport
	s.protocol = protocol
	s.timeout = Timeout{
		ConnTimeoutMs:  uint32(connTimeoutMs),
		ReadTimeoutMs:  uint32(readTimeoutMs),
		WriteTimeoutMs: uint32(writeTimeoutMs),
	}

	s.balance = balance
	s.host = host
	s.maxRetry = uint32(maxRetry)

	// 初始化method conf
	s.initMethodConf(sysConf)

	return nil
}

func (s *Service) initMethodConf(sysConf map[string]interface{}) (err error) {
	method, ok := sysConf["method"]
	if !ok {
		return
	}

	methodArray, ok := method.([]interface{})
	if !ok {
		return
	}

	for _, v := range methodArray {
		vMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		connTimeoutMs, err := getIntFromMap(vMap, "conn_timeout_ms")
		if err != nil {
			connTimeoutMs = ServiceDefaultTimeout
		}
		readTimeoutMs, err := getIntFromMap(vMap, "read_timeout_ms")
		if err != nil {
			readTimeoutMs = ServiceDefaultTimeout
		}

		writeTimeoutMs, err := getIntFromMap(vMap, "write_timeout_ms")
		if err != nil {
			writeTimeoutMs = ServiceDefaultTimeout
		}

		methodName, err := getStringFromMap(vMap, "method")
		if err != nil {
			continue
		}

		s.timeoutMap[methodName] = Timeout{
			ConnTimeoutMs:  uint32(connTimeoutMs),
			ReadTimeoutMs:  uint32(readTimeoutMs),
			WriteTimeoutMs: uint32(writeTimeoutMs),
		}
	}

	return
}
