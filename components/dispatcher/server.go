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

type DispatcherServer struct {
	base             *engine.Server
	gameClientEvents *engine.ClientEvents
}

func (g *DispatcherServer) OnBoot(e gnet.Engine) (action gnet.Action) {
	g.base.OnBoot(e)
	return
}
func (g *DispatcherServer) OnTick() (delay time.Duration, action gnet.Action) {
	return
}
func (g *DispatcherServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	g.base.OnOpen(c)
	return
}
func (g *DispatcherServer) OnShutdown(e gnet.Engine) {
	g.base.OnShutdown(e)
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
	context.SetProperty(util.GameConnMgrKey, g.gameClientEvents)
	g.base.ClientCtxChan <- context
	return gnet.None
}

func (g *DispatcherServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	g.base.OnClose(c, err)
	return gnet.None
}
func (g *DispatcherServer) AddRouter(routers ...engine.IRouter) {
	g.base.AddRouter(routers...)
}
func (g *DispatcherServer) AddHandler(method string, f func(c *protocol.Context)) {
	g.base.AddHandler(method, f)
}
func (g *DispatcherServer) AddHandlerUint32(hash uint32, f func(c *protocol.Context)) {
	g.base.AddHandlerUint32(hash, f)
}

func (g *DispatcherServer) DispatcherInitialize() error {
	g.gameClientEvents = engine.NewClientEvents(util.ClusterTypeGame)
	g.gameClientEvents.AddRouter()
	cli, err := gnet.NewClient(g.gameClientEvents)
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
func (g *DispatcherServer) Run() {
	addr := fmt.Sprintf("tcp://:%d", g.base.Config.Port)
	f := func() {
		err := gnet.Run(g, addr,
			gnet.WithMulticore(true),
			gnet.WithSocketSendBuffer(g.base.Config.ConnWriteBuffer),
			gnet.WithSocketRecvBuffer(g.base.Config.ConnWriteBuffer),
			//gnet.WithTCPKeepAlive()
		)
		panic(err)
	}
	defer func() {
		g.base.HandlerMgr.GPool.Release()
	}()
	util.PanicRepeatRun(f, util.PanicRepeatRunArgs{
		Sleep: time.Second,
		Try:   20,
	})
}
