package router

import (
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
)

type HeartBeat struct {
	engine.BaseRouter
}

func (h *HeartBeat) Init(v ...interface{}) engine.IRouter {
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
func (h *HeartBeat) Health(ctx *protocol.Context) {
	msg := ""
	err := ctx.Bind(&msg)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
	ctx.Send("ok")
}
