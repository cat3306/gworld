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
	req := &gameobject.PlayerPos{}
	msg, err := engine.GameBind(req, ctx)
	if err != nil {
		glog.Logger.Sugar().Errorf("GameBind err:%s", err.Error())
		return
	}
	//glog.Logger.Sugar().Infof("%s", util.BytesToString(msg.ClientMsg.Payload))
	obj := p.Players.Get(req.NetObjId)
	if obj == nil {
		glog.Logger.Sugar().Errorf("not found player id:%s", req.NetObjId)
		return
	}
	obj.OnMove(req.Vector3, gameobject.Vector3{X: req.CX, Y: req.Yaw})
	engine.GameBroadcast(ctx, msg.ClientMsg.Payload, msg.ClientIds)
	//iMsg := &engine.InnerMsg{
	//	ClientIds:      msg.ClientIds,
	//	ClientCodeType: uint32(ctx.CodeType),
	//	ClientMsg: &engine.ClientMsg{
	//		Payload: msg.ClientMsg.Payload,
	//		Method:  msg.ClientMsg.Method,
	//	},
	//}
	//ctx.SendWithParams(iMsg, protocol.ProtoBuffer, util.CallClient)
}

func (p *PlayerMgr) CreatePlayer(ctx *protocol.Context) {
	msg, err := engine.GetCtxInnerMsg(ctx)
	if err != nil {
		glog.Logger.Sugar().Errorf("GetCtxInnerMsg err:%s", err.Error())
		return
	}
	glog.Logger.Info("haha")
	playerId := util.GenId(8)
	player := &gameobject.Player{
	}
	player.OnCreated(playerId)
	p.Players.Add(player)
	//iMsg := &engine.InnerMsg{
	//	ClientIds:      msg.ClientIds,
	//	ClientCodeType: uint32(protocol.String),
	//	ClientMsg: &engine.ClientMsg{
	//		Payload: []byte(playerId),
	//		Method:  msg.ClientMsg.Method,
	//	},
	//}
	engine.GameBroadcast(ctx, playerId, msg.ClientIds)
	//ctx.SendWithParams(iMsg, protocol.ProtoBuffer, util.CallClient)
}
