package main

import (
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
	"github.com/panjf2000/gnet/v2"
	"time"
)

type GameClientProxy struct {
	*gnet.BuiltinEventEngine
	handlerMgr    *engine.HandlerManager
	clientCtxChan chan *protocol.Context
	ConnMgr       *engine.ConnManager
	server        *GateServer
	connMgrs      map[uint32]*engine.ConnManager
	engine.BaseRouter
}

func (ev *GameClientProxy) Init(interface{}) engine.IRouter {
	return ev
}
func NewGameClientProxy() *GameClientProxy {
	return &GameClientProxy{
		ConnMgr:       engine.NewConnManager(),
		handlerMgr:    engine.NewHandlerManager(),
		clientCtxChan: make(chan *protocol.Context, util.ChanPacketSize),
		connMgrs:      map[uint32]*engine.ConnManager{},
	}
}
func (ev *GameClientProxy) SetServer(g *GateServer) *GameClientProxy {
	ev.server = g
	return ev
}
func (ev *GameClientProxy) OnBoot(e gnet.Engine) (action gnet.Action) {
	go ev.mainRoutine()
	return
}
func (ev *GameClientProxy) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetId(util.GenConnId())
	ev.ConnMgr.Add(c)
	glog.Logger.Sugar().Infof("game client proxy conn:%s", c.ID())
	return nil, gnet.None
}

func (ev *GameClientProxy) OnClose(c gnet.Conn, err error) gnet.Action {
	reason := ""
	if err != nil {
		reason = err.Error()
	}
	glog.Logger.Sugar().Infof("cid:%s close,reason:%s", c.ID(), reason)
	ev.ConnMgr.Remove(c.ID())
	return gnet.None
}

func (ev *GameClientProxy) OnTraffic(c gnet.Conn) (action gnet.Action) {
	ctx, err := protocol.Decode(c)
	if err != nil {
		glog.Logger.Sugar().Errorf("OnTraffic err:%s", err.Error())
		return
	}
	if ctx == nil {
		panic("context nil")
	}
	ctx.SetProperty(util.GateClientMgrKey, ev.server.ConnMgr)
	ev.clientCtxChan <- ctx
	return gnet.None
}

func (ev *GameClientProxy) OnTick() (delay time.Duration, action gnet.Action) {

	return
}
func (ev *GameClientProxy) mainRoutine() {
	for {
		select {
		case ctx := <-ev.clientCtxChan:
			ev.handlerMgr.ExeHandler(ctx)
		}
	}
}

func (ev *GameClientProxy) AddRouter(routers ...engine.IRouter) *GameClientProxy {
	for _, v := range routers {
		ev.handlerMgr.RegisterRouter(v.Init(nil))
	}
	return ev
}
func (ev *GameClientProxy) AddHandler(method string, f func(c *protocol.Context)) *GameClientProxy {
	ev.handlerMgr.Register(util.MethodHash(method), f)
	return ev
}
func (ev *GameClientProxy) AddHandlerUint32(hash uint32, f func(c *protocol.Context)) *GameClientProxy {
	ev.handlerMgr.Register(hash, f)
	return ev
}

func (ev *GameClientProxy) HandleGame(ctx *protocol.Context) {
	s := engine.ServerInnerMsg{}
	err := ctx.Bind(&s)
	if err != nil {
		glog.Logger.Sugar().Errorf("HandleGame err:%s", err.Error())
		return
	}
	ev.server.ConnMgr.SendByOne(protocol.Encode(s.Payload, s.ClientCodeType, s.ClientMethod), s.ClientId)
}
