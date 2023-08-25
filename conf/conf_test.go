package conf

import (
	"github.com/cat3306/gocommon/confutil"
	"github.com/cat3306/goworld/util"
	"testing"
)

func TestSaveConf(t *testing.T) {
	c := confutil.Config{}
	c.Save("./cluster.json", ClusterConf{
		Servers: map[util.ClusterType][]ServerConf{
			util.ClusterTypeGate: []ServerConf{
				{
					Ip:              "0.0.0.0",
					Port:            8888,
					MaxConn:         1000,
					ConnWriteBuffer: 1048576,
					ConnReadBuffer:  1048576,
					KV:              map[string]interface{}{},
					Logic:           "gate",
				},
			},
			util.ClusterTypeGame: []ServerConf{
				{
					Logic:           "base",
					Ip:              "127.0.0.1",
					Port:            8890,
					MaxConn:         1000,
					ConnWriteBuffer: 1048576,
					ConnReadBuffer:  1048576,
					KV:              map[string]interface{}{},
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
