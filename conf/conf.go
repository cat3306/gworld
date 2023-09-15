package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cat3306/gocommon/confutil"
	"github.com/cat3306/goworld/util"
)

var (
	GlobalConf       *ClusterConf
	GlobalServerConf *ServerConf
)

type MysqlConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	User         string `json:"user"`
	Pwd          string `json:"pwd"`
	ConnPoolSize int    `json:"conn_pool_size"`
	SetLog       bool   `json:"set_log"`
}
type RedisConfig struct {
	Dbs      []int  `json:"dbs"`
	Addr     string `json:"addr"`
	Password string `json:"password"`
}
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
		for _, v := range list {
			if v.Idx == index {
				GlobalServerConf = &v
				return &v
			}
		}
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
	Name            string                 `json:"name"`
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
	Online          bool                   `json:"online"`
}

func MapToStruct(v interface{}, m map[string]interface{}) error {
	raw, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, v)
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
