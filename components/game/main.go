package main

import (
	"flag"
	"fmt"
	"github.com/cat3306/goworld/components/game/router"
	"github.com/cat3306/goworld/components/game/server"
	"github.com/cat3306/goworld/conf"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/util"
	"os"
	"os/signal"
	"syscall"
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
	signalChan := make(chan os.Signal, 1)
	config := conf.GlobalConf.Select(util.ClusterTypeGame, idx)
	if config.Deploy.IsDaemon {
		ctx, err := engine.DaemonMode(&config.Deploy, fmt.Sprintf("game-%d", idx))
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

	//server := &GameServer{
	//	Server: engine.NewEngine(config, util.ClusterTypeGame),
	//}
	gameServer := server.NewGameServer(config, util.ClusterTypeGame)
	gameServer.AddRouter(
		new(engine.ClientMgr).Init(gameServer.ConnMgr),
		new(router.HeartBeat),
		new(router.RoomMgr),
		new(router.PlayerMgr).Init(),
	)
	setupSignals(signalChan, gameServer)
	gameServer.Run()
}

func setupSignals(ch chan os.Signal, server *server.GameServer) {
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt, syscall.SIGKILL)
	go func() {
		for sig := range ch {
			server.HandlerExit()
			glog.Logger.Sugar().Infof("%+v,game server exit graceful!", sig)
			os.Exit(0)
		}
	}()
}
