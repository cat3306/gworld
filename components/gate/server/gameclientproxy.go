package server

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
	logicMgr      *engine.LogicConnMgr
	engine.BaseRouter
	tryConnectChan chan *engine.TryConnectMsg
	gClient        *gnet.Client
}

func NewGameClientProxy() *GameClientProxy {
	return &GameClientProxy{
		ConnMgr:        engine.NewConnManager(),
		handlerMgr:     engine.NewHandlerManager(),
		clientCtxChan:  make(chan *protocol.Context, util.ChanPacketSize),
		logicMgr:       engine.NewLogicConnMgr(),
		tryConnectChan: make(chan *engine.TryConnectMsg, 1024),
	}
}

func (ev *GameClientProxy) GetLogicMgr() *engine.LogicConnMgr {
	return ev.logicMgr
}
func (ev *GameClientProxy) SetGClient(c *gnet.Client) *GameClientProxy {
	ev.gClient = c
	return ev
}
func (ev *GameClientProxy) OnBoot(e gnet.Engine) (action gnet.Action) {
	go ev.mainRoutine()
	go ev.tryConnect()
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
	ev.logicMgr.Remove(c.ID())
	ev.tryConnectChan <- &engine.TryConnectMsg{
		NetWork: c.RemoteAddr().Network(),
		Addr:    c.RemoteAddr().String(),
	}
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
		ev.handlerMgr.RegisterRouter(v)
	}
	return ev
}
func (ev *GameClientProxy) tryConnect() {
	for msg := range ev.tryConnectChan {
		_, err := ev.gClient.Dial(msg.NetWork, msg.Addr)
		if err != nil {
			glog.Logger.Sugar().Errorf("tryConnect err:%s,msg:%+v", err.Error(), msg)
			time.Sleep(time.Second)
			ev.tryConnectChan <- msg
		}
	}
}
