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
	base                   *engine.Server
	dispatcherClientEvents *engine.ClientEvents
}

func (g *GateServer) OnBoot(e gnet.Engine) (action gnet.Action) {
	g.base.OnBoot(e)
	return
}
func (g *GateServer) OnTick() (delay time.Duration, action gnet.Action) {
	return
}
func (g *GateServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	g.base.OnOpen(c)
	return
}
func (g *GateServer) OnShutdown(e gnet.Engine) {
	g.base.OnShutdown(e)
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
	g.base.ClientCtxChan <- context

	return gnet.None
}

func (g *GateServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	g.base.OnClose(c, err)
	return gnet.None
}
func (g *GateServer) AddRouter(routers ...engine.IRouter) {
	g.base.AddRouter(routers...)
}
func (g *GateServer) AddHandler(method string, f func(c *protocol.Context)) {
	g.base.AddHandler(method, f)
}
func (g *GateServer) AddHandlerUint32(hash uint32, f func(c *protocol.Context)) {
	g.base.AddHandlerUint32(hash, f)
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
	go g.base.MainRoutine()
	util.PanicRepeatRun(f, util.PanicRepeatRunArgs{
		Sleep: time.Second,
		Try:   20,
	})
}
