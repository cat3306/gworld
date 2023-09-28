package engine

import (
	"github.com/cat3306/goworld/conf"
	"github.com/cat3306/goworld/util"
	"math"

	"github.com/cat3306/goworld/protocol"

	"github.com/cat3306/goworld/glog"

	"github.com/panjf2000/gnet/v2"
)

type Server struct {
	ConnMgr *ConnManager
	gnet.BuiltinEventEngine
	eng           gnet.Engine
	HandlerMgr    *HandlerManager
	Config        *conf.ServerConf
	ct            util.ClusterType
	ClientCtxChan chan *protocol.Context
}

func NewEngine(c *conf.ServerConf, ct util.ClusterType) *Server {
	return &Server{
		ConnMgr:       NewConnManager(),
		HandlerMgr:    NewHandlerManager(),
		Config:        c,
		ct:            ct,
		ClientCtxChan: make(chan *protocol.Context, util.ChanPacketSize),
	}
}
func (s *Server) OnBoot(e gnet.Engine) (action gnet.Action) {
	s.eng = e
	go s.MainRoutine()
	glog.Logger.Sugar().Infof("%s server is listening on:%d", s.ct, s.Config.Port)
	return
}
func (s *Server) OnTraffic(c gnet.Conn) gnet.Action {
	//defer func() {
	//	err := recover()
	//	if err != nil {
	//		glog.Logger.Sugar().Errorf("OnTraffic panic %v", err)
	//	}
	//}()
	//s.eng.CountConnections()
	context, err := protocol.Decode(c)
	if err != nil {
		glog.Logger.Sugar().Warnf("OnTraffic err:%s", err.Error())
		return gnet.None
	}
	if context == nil {
		panic("context nil")
	}
	s.ClientCtxChan <- context
	return gnet.None
}
func (s *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	reason := ""
	if err != nil {
		reason = err.Error()
	}
	glog.Logger.Sugar().Infof("cid:%s close,reason:%s", c.ID(), reason)
	s.ConnMgr.Remove(c.ID())
	return gnet.None
}
func (s *Server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	cId := util.GenConnId()
	c.SetId(cId)
	s.ConnMgr.Add(c)
	glog.Logger.Sugar().Infof("clinet conn cid:%s connect", c.ID())
	return out, gnet.None
}
func (s *Server) OnShutdown(e gnet.Engine) {

}

func (s *Server) AddRouter(routers ...IRouter) {
	for _, v := range routers {
		s.HandlerMgr.RegisterRouter(v)
	}
}
func (s *Server) AddHandler(method string, f func(c *protocol.Context)) {
	s.HandlerMgr.Register(util.MethodHash(method), f)
}
func (s *Server) AddHandlerUint32(hash uint32, f func(c *protocol.Context)) {
	s.HandlerMgr.Register(hash, f)
}

func (s *Server) MainRoutine() {
	f := func() {
		for {
			select {
			case ctx := <-s.ClientCtxChan:
				s.HandlerMgr.ExeHandler(ctx)
			}
		}
	}
	util.PanicRepeatRun(f, util.PanicRepeatRunArgs{
		Sleep: 0,
		Try:   math.MaxInt64,
	})

}
