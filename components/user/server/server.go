package server

import (
	"fmt"
	"github.com/cat3306/goworld/components/game/router"
	"github.com/cat3306/goworld/conf"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/bytebufferpool"
	"math"
	"time"
)

type UserServer struct {
	ConnMgr *engine.ConnManager
	gnet.BuiltinEventEngine
	eng           gnet.Engine
	HandlerMgr    *engine.HandlerManager
	Config        *conf.ServerConf
	ct            util.ClusterType
	ClientCtxChan chan *protocol.Context
}

func (g *UserServer) OnShutdown(e gnet.Engine) {

}
func (g *UserServer) Run() {
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
		Try:   math.MaxInt64,
	})
}
func (g *UserServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	//cId := util.GenConnId()
	//c.SetId(cId)
	//g.ConnMgr.Add(c)
	glog.Logger.Sugar().Infof("gate clinet conn cid:%s connect", c.ID())
	buffer := protocol.Encode(g.Config.Logic, protocol.String, util.MethodHash("SetLogic"), 0)

	copyOut := make([]byte, buffer.Len())
	copy(copyOut, buffer.Bytes())
	out = copyOut
	bytebufferpool.Put(buffer)
	return out, gnet.None
}
func (g *UserServer) OnBoot(e gnet.Engine) (action gnet.Action) {
	g.eng = e
	go g.MainRoutine()
	glog.Logger.Sugar().Infof("%s server is listening on:%d", g.ct, g.Config.Port)
	return
}
func (g *UserServer) MainRoutine() {
	f := func() {
		for {
			select {
			case ctx := <-g.ClientCtxChan:
				g.HandlerMgr.ExeHandler(ctx)
			}
		}
	}
	util.PanicRepeatRun(f, util.PanicRepeatRunArgs{
		Sleep: 0,
		Try:   math.MaxInt64,
	})
}
func (g *UserServer) OnTraffic(c gnet.Conn) gnet.Action {
	ctx, err := protocol.Decode(c)
	if err != nil {
		glog.Logger.Sugar().Warnf("OnTraffic err:%s", err.Error())
		return gnet.None
	}
	if ctx == nil {
		panic("context nil")
	}

	g.ClientCtxChan <- ctx
	return gnet.None
}
func (g *UserServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	reason := ""
	if err != nil {
		reason = err.Error()
	}
	glog.Logger.Sugar().Infof("cid:%s close,reason:%s", c.ID(), reason)
	g.ConnMgr.Remove(c.ID())
	return gnet.None
}

func (g *UserServer) HandlerExit() {
	router.SaveData()
}

func (g *UserServer) AddRouter(routers ...engine.IRouter) {
	g.HandlerMgr.SetPreHandlers([]engine.Handler{
		engine.PreSetInnerMsgMsg(),
	}...)
	for _, v := range routers {
		g.HandlerMgr.RegisterRouter(v)
	}
}

func NewUserServer(c *conf.ServerConf, ct util.ClusterType) *UserServer {
	return &UserServer{
		ConnMgr:       engine.NewConnManager(),
		HandlerMgr:    engine.NewHandlerManager(),
		Config:        c,
		ct:            ct,
		ClientCtxChan: make(chan *protocol.Context, util.ChanPacketSize),
	}
}
