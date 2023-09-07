package router

import (
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
)

type PlayerMgr struct {
	clients map[string]*GameClient
}

func (p *PlayerMgr) Init(v interface{}) engine.IRouter {
	return p
}
func (p *PlayerMgr) PlayerMove(ctx *protocol.Context) {
	msg := &engine.InnerMsg{}
	err := ctx.Bind(msg)
	if err != nil {
		glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
		return
	}
	glog.Logger.Sugar().Infof("%s", string(msg.ClientMsg.Payload))

	iMsg := &engine.InnerMsg{
		ClientIds: msg.ClientIds,
		ClientMsg: &engine.ClientMsg{
			Logic:    0,
			Payload:  msg.ClientMsg.Payload,
			Method:   msg.ClientMsg.Method,
			CodeType: uint32(protocol.Json),
		},
	}
	ctx.SendWithParams(iMsg, protocol.ProtoBuffer, util.CallClient)
}

func (p *PlayerMgr) CreatePlayer(ctx *protocol.Context) {
	msg := &engine.InnerMsg{}
	err := ctx.Bind(msg)
	if err != nil {
		glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
		return
	}
	glog.Logger.Sugar().Infof("%s", string(msg.ClientMsg.Payload))

	iMsg := &engine.InnerMsg{
		ClientIds: msg.ClientIds,
		ClientMsg: &engine.ClientMsg{
			Logic:    0,
			Payload:  []byte(util.GenId(8)),
			Method:   msg.ClientMsg.Method,
			CodeType: uint32(protocol.String),
		},
	}
	ctx.SendWithParams(iMsg, protocol.ProtoBuffer, util.CallClient)
}
