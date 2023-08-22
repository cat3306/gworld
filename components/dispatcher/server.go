package main

import (
	"fmt"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
	"github.com/panjf2000/gnet/v2"
	"time"
)

type DispatcherServer struct {
	*engine.Server
	gateClients *engine.ConnManager
	gameClients *engine.ConnManager
}

func (g *DispatcherServer) OnTraffic(c gnet.Conn) gnet.Action {
	context, err := protocol.Decode(c)
	if err != nil {
		glog.Logger.Sugar().Warnf("OnTraffic err:%s", err.Error())
		return gnet.None
	}
	if context == nil {
		panic("context nil")
	}
	context.SetProperty(util.GameClientMgrKey, g.gameClients)
	context.SetProperty(util.GateClientMgrKey, g.gateClients)
	g.ClientCtxChan <- context
	return gnet.None
}

func (g *DispatcherServer) Run() {
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
