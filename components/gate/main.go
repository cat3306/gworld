package main

import (
	"flag"
	"fmt"
	"github.com/cat3306/goworld/components/gate/router"
	"github.com/cat3306/goworld/components/gate/server"
	"github.com/cat3306/goworld/conf"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/util"
)

func ParseFlag() (string, string, int) {
	var file string
	var idx int
	var name string
	flag.StringVar(&file, "c", "", "use -c to bind conf file")
	flag.IntVar(&idx, "idx", 0, "set which conf use")
	flag.StringVar(&name, "name", "local", "specify namespace")
	flag.Parse()
	return file, name, idx
}
func main() {
	file, name, idx := ParseFlag()
	glog.Init()
	err := conf.LoadConf(file, name)
	if err != nil {
		glog.Logger.Sugar().Errorf("conf.Load err:%s", err.Error())
		return
	}
	config := conf.GlobalConf.Select(util.ClusterTypeGate, idx)
	if config.Deploy.IsDaemon {
		ctx, err := engine.DaemonMode(&config.Deploy, fmt.Sprintf("gate-%d", idx))
		if err != nil {
			glog.Logger.Sugar().Errorf("DaemonMode err:%s", err.Error())
			return
		}
		defer func() {
			err := ctx.Release()
			if err != nil {
				glog.Logger.Sugar().Errorf("ctx.Release() err:%s", err.Error())
			}
		}()
	}
	gate := server.NewGateServer(config, util.ClusterTypeGate)
	gate.AddRouter(
		new(router.GateDispatcher).Init(gate.GetGameClientProxy().GetLogicMgr()),
		new(router.Auth).Init(),
	)
	gate.GetGameClientProxy().AddRouter(
		new(router.GameClientProxyRouter).Init(gate.ConnMgr, gate.GetGameClientProxy().GetLogicMgr()),
	)
	//server.AddHandler("dispatcher", server.Dispatcher)
	err = gate.GameInitialize()
	if err != nil {
		panic(err)
	}
	gate.Run()
}
