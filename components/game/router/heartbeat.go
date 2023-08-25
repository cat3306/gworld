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

func (h *HeartBeat) Init(v interface{}) engine.IRouter {
	return h
}
func (h *HeartBeat) HeartBeat(ctx *protocol.Context) {
	str := ""
	err := ctx.Bind(&str)
	glog.Logger.Sugar().Infof("HeartBeat:%s", str)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
}

func (h *HeartBeat) GlobalHeartBeat(ctx *protocol.Context) {
	s := &engine.InnerMsg{}
	err := ctx.Bind(s)
	//glog.Logger.Sugar().Infof("HeartBeat:%s", str)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
	s.ClientMsg.Payload = []byte("️❤️")
	ctx.SendWithParams(s, protocol.ProtoBuffer, util.CallGate)

}
