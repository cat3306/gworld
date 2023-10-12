package router

import (
	"encoding/json"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/engine/erouter"
	"github.com/cat3306/goworld/engine/gameobject"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
)

var (
	PlayerManager *PlayerMgr
)

type PlayerMgr struct {
	Players gameobject.GameObjectSet `json:"players"`
}

func (p *PlayerMgr) Init(v ...interface{}) engine.IRouter {
	p.Players = make(gameobject.GameObjectSet)
	PlayerManager = p
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
	//engine.GameBroadcast(ctx, msg.ClientMsg.Payload, msg.ClientIds)
	erouter.GameClientMgr.Broadcast(ctx, nil, msg.ClientMsg.Payload)
}

func (p *PlayerMgr) CreatePlayer(ctx *protocol.Context) {
	msg, err := engine.GetCtxInnerMsg(ctx)
	if err != nil {
		glog.Logger.Sugar().Errorf("GetCtxInnerMsg err:%s", err.Error())
		return
	}
	glog.Logger.Info("haha")
	playerId := util.GenGameObjectId()
	player := &gameobject.Player{}
	player.OnCreated(playerId)
	p.Players.Add(player)
	erouter.GameClientMgr.Broadcast(ctx, msg.ClientIds, playerId)
}
func (p *PlayerMgr) SaveData() {
	if len(p.Players) == 0 {
		return
	}
	raw, err := json.Marshal(p)
	if err != nil {
		glog.Logger.Sugar().Infof("SaveData err:%s", err.Error())
		return
	}
	glog.Logger.Sugar().Infof("%s", string(raw))

}
