package main

import (
	"fmt"
	"github.com/cat3306/goworld/components/gate/router"
	"github.com/cat3306/goworld/conf"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
	"github.com/panjf2000/gnet/v2"
	"time"
)

type GateServer struct {
	*engine.Server
	dispatcherClientEvents *engine.ClientEvents
}



func (g *GateServer) OnTraffic(c gnet.Conn) gnet.Action {
	context, err := protocol.Decode(c)
	if err != nil {
		glog.Logger.Sugar().Warnf("OnTraffic err:%s", err.Error())
		return gnet.None
	}
	if context == nil {
		panic("context nil")
	}
	context.SetProperty(util.DispatcherConnMgrKey, g.dispatcherClientEvents.ConnMgr)
	g.ClientCtxChan <- context

	return gnet.None
}

func (g *GateServer) DispatcherInitialize() error {
	g.dispatcherClientEvents = engine.NewClientEvents(util.ClusterTypeGate)
	g.dispatcherClientEvents.AddRouter(new(router.DisHeartBeat))
	cli, err := gnet.NewClient(g.dispatcherClientEvents)
	if err != nil {
		return err
	}
	list := conf.GlobalConf.ClusterList(util.ClusterTypeDispatcher)
	for _, v := range list {
		_, err = cli.Dial("tcp", fmt.Sprintf("%s:%d", v.Ip, v.Port))
		if err != nil {
			return err
		}
	}
	//go func() {
	//	for {
	//		raw := protocol.Encode("💓", protocol.String, util.MethodHash("HeartBeat"))
	//		g.dispatcherClientEvents.ConnMgr.Broadcast(raw)
	//		time.Sleep(time.Second * 1)
	//	}
	//
	//}()
	return cli.Start()
}
func (g *GateServer) Run() {
	addr := fmt.Sprintf("tcp://:%d", g.Config.Port)
	f := func() {
		err := gnet.Run(g, addr,
			gnet.WithMulticore(true),
			gnet.WithSocketSendBuffer(g.Config.ConnWriteBuffer),
			gnet.WithSocketRecvBuffer(g.Config.ConnWriteBuffer),
			//gnet.WithTCPKeepAlive()
		)
		panic(err)
	}
	defer func() {
		g.HandlerMgr.GPool.Release()
	}()
	util.PanicRepeatRun(f, util.PanicRepeatRunArgs{
		Sleep: time.Second,
		Try:   20,
	})
}
