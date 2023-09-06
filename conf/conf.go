package conf

import (
	"github.com/cat3306/gocommon/confutil"
	"github.com/cat3306/goworld/util"
)

var (
	GlobalConf ClusterConf
)

type AuthConfig struct {
	IsAuth         bool   `json:"is_auth"` //是否客户端验签
	PrivateKeyPath string `json:"private_key_path"`
}
type ClusterConf struct {
	AuthConfig AuthConfig                        `json:"auth_config"`
	Servers    map[util.ClusterType][]ServerConf `json:"servers"`
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
	Host            string                 `json:"host"`     //ip
	Port            int                    `json:"tcp_port"` //port
	MaxConn         int                    `json:"max_conn"` //最大连接数
	ConnWriteBuffer int                    `json:"conn_write_buffer"`
	ConnReadBuffer  int                    `json:"conn_read_buffer"`
	KV              map[string]interface{} `json:"kv"`
	OuterIp         string                 `json:"outer_ip"`
}

func Load(file string) error {
	config := confutil.Config{}
	return config.Load(file, &GlobalConf)
}
