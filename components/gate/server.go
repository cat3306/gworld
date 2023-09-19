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
	gameClientProxy          *GameClientProxy
	innerGameServerBroadcast uint32
}

func (g *GateServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	reason := ""
	if err != nil {
		reason = err.Error()
	}
	glog.Logger.Sugar().Infof("cid:%s close,reason:%s", c.ID(), reason)
	g.ConnMgr.Remove(c.ID())
	if g.innerGameServerBroadcast == 0 {
		g.innerGameServerBroadcast = util.MethodHash("InnerOnBroadcast")
	}
	ctx := &protocol.Context{
		Proto: g.innerGameServerBroadcast,
		Conn:  c,
	}
	ctx.SetProperty("Proto", util.MethodHash("OnDisconnect"))
	ctx.SetProperty("cid", c.ID())
	g.ClientCtxChan <- ctx
	return gnet.None
}

func (g *GateServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	cId := util.GenConnId()
	c.SetId(cId)
	g.ConnMgr.Add(c)
	glog.Logger.Sugar().Infof("clinet conn cid:%s connect", c.ID())
	if g.innerGameServerBroadcast == 0 {
		g.innerGameServerBroadcast = util.MethodHash("InnerOnBroadcast")
	}
	ctx := &protocol.Context{
		Proto: g.innerGameServerBroadcast,
		Conn:  c,
	}
	ctx.SetProperty("cid", c.ID())
	ctx.SetProperty("Proto", util.MethodHash("OnConnect"))
	g.ClientCtxChan <- ctx
	return out, gnet.None
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
	g.ClientCtxChan <- context

	return gnet.None
}

func (g *GateServer) GameInitialize() error {
	g.gameClientProxy = NewGameClientProxy().SetServer(g)
	g.gameClientProxy.AddRouter(new(GameClientProxyRouter).Init(g))
	cli, err := gnet.NewClient(g.gameClientProxy)
	if err != nil {
		return err
	}
	g.gameClientProxy.SetGClient(cli)
	list := conf.GlobalConf.ClusterList(util.ClusterTypeGame)
	for _, v := range list {
		if !v.Online {
			continue
		}
		g.gameClientProxy.tryConnectChan <- &engine.TryConnectMsg{
			NetWork: "tcp",
			Addr:    fmt.Sprintf("%s:%d", v.OuterIp, v.Port),
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
