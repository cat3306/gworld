package conf

import (
	"github.com/cat3306/gocommon/confutil"
	"github.com/cat3306/goworld/util"
)

var (
	GlobalConf ClusterConf
)

type ClusterConf struct {
	Servers map[util.ClusterType][]ServerConf `json:"servers"`
}

func (s *ClusterConf) Select(t util.ClusterType, index int) *ServerConf {
	if list, ok := s.Servers[t]; ok {
		return &list[index]
	}
	panic("not found cluster type")
}
func (s *ClusterConf) ClusterList(t util.ClusterType) []ServerConf {
	if list, ok := s.Servers[t]; ok {
		return list
	}
	panic("not found cluster type")
}
func (s *ClusterConf) ClusterIdxs(t util.ClusterType) (idx []int) {
	if list, ok := s.Servers[t]; ok {
		for i := range list {
			idx = append(idx, i)
		}
	}
	panic("not found cluster type")
	return
}

type ServerConf struct {
	Logic           string                 `json:"logic"`
	Ip              string                 `json:"host"`     //ip
	Port            int                    `json:"tcp_port"` //port
	MaxConn         int                    `json:"max_conn"` //最大连接数
	ConnWriteBuffer int                    `json:"conn_write_buffer"`
	ConnReadBuffer  int                    `json:"conn_read_buffer"`
	KV              map[string]interface{} `json:"kv"`
}

func Load(file string) error {
	config := confutil.Config{}
	return config.Load(file, &GlobalConf)
}
