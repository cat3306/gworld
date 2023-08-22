package engine

import (
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
	"github.com/panjf2000/gnet/v2"
	"time"
)

type ClientEvents struct {
	*gnet.BuiltinEventEngine
	handlerMgr    *HandlerManager
	packetLen     int
	clientCtxChan chan *protocol.Context
	ConnMgr       *ConnManager
	clusterType   util.ClusterType
}

func NewClientEvents(ct util.ClusterType) *ClientEvents {
	return &ClientEvents{
		ConnMgr:       NewConnManager(),
		handlerMgr:    NewHandlerManager(),
		clientCtxChan: make(chan *protocol.Context, util.ChanPacketSize),
		clusterType:   ct,
	}
}
func (ev *ClientEvents) OnBoot(e gnet.Engine) (action gnet.Action) {
	go ev.mainRoutine()
	return
}
func (ev *ClientEvents) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetId(util.GenId(9))
	ev.ConnMgr.Add(c)
	glog.Logger.Sugar().Info("")
	return nil, gnet.None
}

func (ev *ClientEvents) OnClose(c gnet.Conn, err error) gnet.Action {
	reason := ""
	if err != nil {
		reason = err.Error()
	}
	glog.Logger.Sugar().Infof("cid:%s close,reason:%s", c.ID(), reason)
	ev.ConnMgr.Remove(c.ID())
	return gnet.None
}

func (ev *ClientEvents) OnTraffic(c gnet.Conn) (action gnet.Action) {
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

func (ev *ClientEvents) OnTick() (delay time.Duration, action gnet.Action) {

	return
}
func (ev *ClientEvents) mainRoutine() {
	for {
		select {
		case ctx := <-ev.clientCtxChan:
			ev.handlerMgr.ExeHandler(ctx)
		}
	}
}

func (ev *ClientEvents) AddRouter(routers ...IRouter) {
	for _, v := range routers {
		ev.handlerMgr.RegisterRouter(v.Init())
	}
}
func (ev *ClientEvents) AddHandler(method string, f func(c *protocol.Context)) {
	ev.handlerMgr.Register(util.MethodHash(method), f)
}
func (ev *ClientEvents) AddHandlerUint32(hash uint32, f func(c *protocol.Context)) {
	ev.handlerMgr.Register(hash, f)
}
