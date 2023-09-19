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
	if !ctx.CheckClientAuth() {
		glog.Logger.Sugar().Warnf("authentication required:%s", ctx.Conn.ID())
		ctx.Conn.Close()
		return
	}
	req := engine.ClientMsg{}
	err := ctx.Bind(&req)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}

	c, ok := g.gate.gameClientProxy.logicMgr.GetByLogic(ctx.Logic)
	if !ok {
		glog.Logger.Sugar().Errorf("not found logic game server,logic:%d", ctx.Logic)
		return
	}
	if c.Len() == 0 {
		glog.Logger.Sugar().Errorf("not found game server,logic:%d", ctx.Logic)
		return
	}
	s := &engine.InnerMsg{
		ClientIds:      []string{ctx.Conn.ID()},
		ClientMsg:      &req,
		ClientCodeType: uint32(ctx.CodeType),
		Properties:     map[string]string{},
	}
	//glog.Logger.Info("haha")
	c.SendSomeOne(protocol.Encode(s, protocol.ProtoBuffer, req.Method, ctx.Logic))
}

func (g *GateDispatcher) HeartBeat(ctx *protocol.Context) {
	str := ""
	err := ctx.Bind(&str)
	//glog.Logger.Sugar().Infof("HeartBeat:%s", str)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
	//glog.Logger.Sugar().Infof("%v", ctx.CodeType)
	ctx.Send("❤️")
}
func (g *GateDispatcher) InnerOnBroadcast(ctx *protocol.Context) {
	v, ok := ctx.GetProperty("Proto")
	if !ok {
		glog.Logger.Sugar().Errorf("not found Proto")
		return
	}
	pro := v.(uint32)

	cidV, ok := ctx.GetProperty("cid")
	if !ok {
		glog.Logger.Sugar().Errorf("not found cid")
		return
	}
	cid := cidV.(string)
	var payload []byte
	if ctx.Payload != nil {
		payload = ctx.Payload.Bytes()
	}
	buffer := protocol.Encode(&engine.InnerMsg{
		ClientIds: []string{cid},
		ClientMsg: &engine.ClientMsg{
			Payload: payload,
		},
	}, protocol.ProtoBuffer, pro, 0)
	g.gate.gameClientProxy.logicMgr.Broadcast(buffer)
}
