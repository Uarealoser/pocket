package discovery

import "strings"

type Server struct {
	idc  string // 机房标志
	ip   string
	port uint32

	weight       float64 //权重，round robin使用
	curWeight    float64
	effectWeight float64
	weightFactor uint32

	enable      uint32 // 是否可用
	offline     int    // 是否下线
	lastActive  int64  // 存活事件戳
	failedCount int32  // 连续失败的次数
}

func NewServer(ip string, port uint32, weight float64, idc string) *Server {
	return &Server{
		ip:           ip,
		port:         port,
		weight:       weight,
		effectWeight: weight,
		enable:       1,
		idc:          strings.ToLower(idc),
		offline:      0,
	}
}
