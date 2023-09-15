package conf

import (
	"github.com/cat3306/gocommon/confutil"
	"testing"
)

func TestGenJson(t *testing.T) {
	c := confutil.Config{}
	c.Save("./server.json", GlobalConf{
		Port:        "8878",
		ClusterPath: "/Users/joker/.gamecluster.json",
	})
}
