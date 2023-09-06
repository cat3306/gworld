package conf

import (
	"github.com/cat3306/gocommon/confutil"
	"github.com/cat3306/goworld/util"
	"testing"
)

func TestSaveConf(t *testing.T) {
	c := confutil.Config{}
	c.Save("./cluster.json", ClusterConf{
		AuthConfig: AuthConfig{
			IsAuth:         true,
			PrivateKeyPath: "/Users/joker/code/go/src/github.com/cat3306/goworld/cert/private_key.pem",
		},
		Servers: map[util.ClusterType][]ServerConf{
			util.ClusterTypeGate: []ServerConf{
				{
					Host:            "0.0.0.0",
					Port:            8888,
					MaxConn:         1000,
					ConnWriteBuffer: 1048576,
					ConnReadBuffer:  1048576,
					KV:              map[string]interface{}{},
					Logic:           "gate",
					OuterIp:         "",
				},
			},
			util.ClusterTypeGame: []ServerConf{
				{
					Logic:           "base",
					Host:            "127.0.0.1",
					Port:            8890,
					MaxConn:         1000,
					ConnWriteBuffer: 1048576,
					ConnReadBuffer:  1048576,
					KV:              map[string]interface{}{},
					OuterIp:         "127.0.0.1",
				},
			},
		},
	})
	cc := ClusterConf{}
	err := c.Load("./cluster.json", &cc)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cc)
}
