package conf

import (
	"errors"
	"fmt"
	"github.com/cat3306/gocommon/confutil"
	"github.com/cat3306/goworld/util"
)

var (
	GlobalConf *ClusterConf
)

type AuthConfig struct {
	IsAuth         bool   `json:"is_auth"` //是否客户端验签
	PrivateKeyPath string `json:"private_key_path"`
}
type LogConf struct {
	Level string `json:"level"`
	Path  string `json:"path"`
}
type DeployConf struct {
	IsDaemon bool   `json:"is_daemon"`
	BinPath  string `json:"bin_path"`
	PidFile  string `json:"pid_file_name"`
	LogFile  string `json:"log_file_name"`
}
type ClusterConf struct {
	Name       string                            `json:"name"`
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
	Idx             int                    `json:"idx"`
	Logic           string                 `json:"logic"`
	Host            string                 `json:"host"`     //ip
	Port            int                    `json:"tcp_port"` //port
	MaxConn         int                    `json:"max_conn"` //最大连接数
	ConnWriteBuffer int                    `json:"conn_write_buffer"`
	ConnReadBuffer  int                    `json:"conn_read_buffer"`
	KV              map[string]interface{} `json:"kv"`
	OuterIp         string                 `json:"outer_ip"`
	Deploy          DeployConf             `json:"deploy"`
}

func Load(file string) error {
	config := confutil.Config{}
	return config.Load(file, &GlobalConf)
}
func LoadConf(filePath string, name string) error {
	config := confutil.Config{}
	var clusterList []ClusterConf
	err := config.Load(filePath, &clusterList)
	if err != nil {
		return err
	}
	var (
		find bool
	)
	for _, v := range clusterList {
		if v.Name == name {
			find = true
			GlobalConf = &v
			break
		}
	}
	if !find || GlobalConf == nil {
		return errors.New(fmt.Sprintf("not found %s config", name))
	}
	return err
}
