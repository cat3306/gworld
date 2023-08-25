package main

import (
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
)

type GateDispatcher struct {
	gate *GateServer
	engine.BaseRouter
}

func (g *GateDispatcher) Init(v interface{}) engine.IRouter {
	g.gate = v.(*GateServer)
	return g
}

func (g *GateDispatcher) Dispatcher(ctx *protocol.Context) {
	req := engine.ClientMsg{}
	err := ctx.Bind(&req)
	glog.Logger.Sugar().Infof("GlobalHeartBeat:%s", req.Payload)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
	c, ok := g.gate.gameClientProxy.connMgrs[req.Logic]
	if !ok {
		glog.Logger.Sugar().Errorf("not found logic game server,logic:%d", req.Logic)
		return
	}
	s := &engine.InnerMsg{
		ClientId:  []string{ctx.Conn.ID()},
		ClientMsg: &req,
	}
	c.SendSomeOne(protocol.Encode(s, protocol.ProtoBuffer, req.Method))
}

func (g *GateDispatcher) HeartBeat(ctx *protocol.Context) {
	str := ""
	err := ctx.Bind(&str)
	glog.Logger.Sugar().Infof("HeartBeat:%s", str)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
	ctx.Send("❤️")
}
