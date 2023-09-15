package conf

import "github.com/cat3306/gocommon/confutil"

var (
	Global ServerConf
)

type ServerConf struct {
	Port        string `json:"port"`
	ClusterPath string `json:"cluster_path"`
}

func Load(filePath string) error {
	config := confutil.Config{}

	return config.Load(filePath, &Global)
}
