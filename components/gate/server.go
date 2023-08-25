package main

import (
	"fmt"
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
	gameClientProxy *GameClientProxy
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
	context.SetProperty(util.GameClientProxyMgrKey, g.gameClientProxy.ConnMgr)
	g.ClientCtxChan <- context

	return gnet.None
}

func (g *GateServer) SetLogic(ctx *protocol.Context) {
	var logic string

	err := ctx.Bind(&logic)
	if err != nil {
		glog.Logger.Sugar().Errorf("SetLogic err:%s", err.Error())
		return
	}
	logicHash := util.MethodHash(logic)
	if _, ok := g.gameClientProxy.connMgrs[logicHash]; ok {
		glog.Logger.Sugar().Errorf("register repeated logic")
		return
	}
	mgr := engine.NewConnManager()
	mgr.Add(ctx.Conn)
	glog.Logger.Sugar().Infof("set logic from game,logic:%s", logic)
	g.gameClientProxy.connMgrs[logicHash] = mgr
}

func (g *GateServer) GameInitialize() error {
	g.gameClientProxy = NewGameClientProxy().SetServer(g)
	//AddHandler("GlobalHeartBeat", router.GlobalHeartBeat2).
	g.gameClientProxy.
		AddHandler("SetLogic", g.SetLogic).
		AddHandlerUint32(util.CallGate, g.gameClientProxy.HandleGame)
	cli, err := gnet.NewClient(g.gameClientProxy)
	if err != nil {
		return err
	}
	list := conf.GlobalConf.ClusterList(util.ClusterTypeGame)
	for _, v := range list {
		_, err = cli.Dial("tcp", fmt.Sprintf("%s:%d", v.Ip, v.Port))
		if err != nil {
			return err
		}
	}
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
