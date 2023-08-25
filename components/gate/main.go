package main

import (
	"flag"
	"github.com/cat3306/goworld/conf"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/util"
)

func ParseFlag() (string, int) {
	var file string
	var idx int
	flag.StringVar(&file, "c", "", "use -c to bind conf file")
	flag.IntVar(&idx, "idx", 0, "set which conf use")
	flag.Parse()
	return file, idx
}
func main() {
	file, idx := ParseFlag()
	glog.Init()
	err := conf.Load(file)
	if err != nil {
		glog.Logger.Sugar().Errorf("conf.Load err:%s", err.Error())
		return
	}
	config := conf.GlobalConf.Select(util.ClusterTypeGate, idx)
	server := GateServer{
		Server: engine.NewEngine(config, util.ClusterTypeGate),
	}
	server.AddRouter(
		new(GateDispatcher).Init(&server),
	)
	//server.AddHandler("dispatcher", server.Dispatcher)
	err = server.GameInitialize()
	if err != nil {
		panic(err)
	}
	server.Run()
}
