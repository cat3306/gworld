package main

import (
	"fmt"
	"github.com/cat3306/goworld/conf"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
	"github.com/panjf2000/gnet/v2"
	"time"
)

type GameServer struct {
	base                   *engine.Server
	dispatcherClientEvents *engine.ClientEvents
}

func (g *GameServer) OnBoot(e gnet.Engine) (action gnet.Action) {
	g.base.OnBoot(e)
	return
}
func (g *GameServer) OnTick() (delay time.Duration, action gnet.Action) {
	return
}
func (g *GameServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	g.base.OnOpen(c)
	return
}
func (g *GameServer) OnShutdown(e gnet.Engine) {
	g.base.OnShutdown(e)
}

func (g *GameServer) OnTraffic(c gnet.Conn) gnet.Action {
	return g.base.OnTraffic(c)
}

func (g *GameServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	g.base.OnClose(c, err)
	return gnet.None
}
func (g *GameServer) AddRouter(routers ...engine.IRouter) {
	g.base.AddRouter(routers...)
}
func (g *GameServer) AddHandler(method string, f func(c *protocol.Context)) {
	g.base.AddHandler(method, f)
}
func (g *GameServer) AddHandlerUint32(hash uint32, f func(c *protocol.Context)) {
	g.base.AddHandlerUint32(hash, f)
}
func (g *GameServer) DispatcherInitialize() error {
	g.dispatcherClientEvents = engine.NewClientEvents(util.ClusterTypeGame)
	g.dispatcherClientEvents.AddRouter()
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
	return cli.Start()
}
func (g *GameServer) Run() {
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
