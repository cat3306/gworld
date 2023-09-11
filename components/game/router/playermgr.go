package router

import (
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/engine/gameobject"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
)

type PlayerMgr struct {
	Players gameobject.GameObjectSet
}

func (p *PlayerMgr) Init(v interface{}) engine.IRouter {
	p.Players = gameobject.GameObjectSet{}
	return p
}

func (p *PlayerMgr) PlayerMove(ctx *protocol.Context) {
	msg := &engine.InnerMsg{}
	err := ctx.Bind(msg)
	if err != nil {
		glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
		return
	}
	req := &gameobject.PosInfo{}
	err = msg.ClientMsg.Bind(req)
	if err != nil {
		glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
		return
	}
	glog.Logger.Sugar().Infof("%s", util.BytesToString(msg.ClientMsg.Payload))
	obj := p.Players.Get(req.NetObjId)
	obj.OnMove(req.Vector3, gameobject.Vector3{})
	iMsg := &engine.InnerMsg{
		ClientIds: msg.ClientIds,
		ClientMsg: &engine.ClientMsg{
			Logic:    0,
			Payload:  msg.ClientMsg.Payload,
			Method:   msg.ClientMsg.Method,
			CodeType: msg.ClientMsg.CodeType,
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
	glog.Logger.Info("haha")
	playerId := util.GenId(8)
	player := &gameobject.Player{
	}
	player.OnCreated(playerId)
	p.Players.Add(player)
	iMsg := &engine.InnerMsg{
		ClientIds: msg.ClientIds,
		ClientMsg: &engine.ClientMsg{
			Logic:    0,
			Payload:  []byte(playerId),
			Method:   msg.ClientMsg.Method,
			CodeType: uint32(protocol.String),
		},
	}
	ctx.SendWithParams(iMsg, protocol.ProtoBuffer, util.CallClient)
}
