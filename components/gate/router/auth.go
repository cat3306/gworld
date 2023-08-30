package router

import (
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
)

type Auth struct {
	engine.BaseRouter
}

func (h *Auth) Init(v interface{}) engine.IRouter {
	return h
}
func (h *Auth) Auth(ctx *protocol.Context) {
	str := ""
	err := ctx.Bind(&str)
	glog.Logger.Sugar().Infof("HeartBeat:%s", str)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
	ctx.Conn.SetProperty(util.ClientAuth, "ok")
	ctx.Send("ok")
}
