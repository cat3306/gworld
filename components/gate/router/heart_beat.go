package router

import (
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
)

type HeartBeat struct {
	engine.BaseRouter
}

func (h *HeartBeat) Init() engine.IRouter {
	return h
}
func (h *HeartBeat) HeartBeat(ctx *protocol.Context) {
	str := ""
	err := ctx.Bind(&str)
	glog.Logger.Sugar().Infof("HeartBeat:%s", str)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
	ctx.Send("*")
}
func (h *HeartBeat) GlobalHeartBeat(ctx *protocol.Context) {
	str := ""
	err := ctx.Bind(&str)
	glog.Logger.Sugar().Infof("HeartBeat:%s", str)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
	v, _ := ctx.GetProperty(util.DispatcherConnMgrKey)
	connMgr := v.(*engine.ConnManager)
	connMgr.Broadcast(protocol.Encode(str, ctx.CodeType, ctx.Proto))
}

type DisHeartBeat struct {
	engine.BaseRouter
}

func (h *DisHeartBeat) Init() engine.IRouter {
	return h
}
func (h *DisHeartBeat) HeartBeat(ctx *protocol.Context) {
	str := ""
	err := ctx.Bind(&str)
	glog.Logger.Sugar().Infof("HeartBeat:%s", str)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
}
