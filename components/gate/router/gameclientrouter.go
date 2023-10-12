package router

import (
	"github.com/cat3306/goworld/conf"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
	"strconv"
)

type GameClientProxyRouter struct {
	connMgr      *engine.ConnManager
	logicConnMgr *engine.LogicConnMgr
}

func (g *GameClientProxyRouter) Init(v ...interface{}) engine.IRouter {
	if len(v) != 0 {
		g.connMgr = v[0].(*engine.ConnManager)
		g.logicConnMgr = v[1].(*engine.LogicConnMgr)
	}

	return g
}
func (g *GameClientProxyRouter) SetClientProperty(ctx *protocol.Context) {
	req := engine.InnerMsg{}
	err := ctx.Bind(&req)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
		return
	}
	for _, v := range req.ClientIds {
		ok, conn := g.connMgr.Get(v)
		if ok {
			for k, val := range req.Properties {
				conn.SetProperty(k, val)
			}
			glog.Logger.Info(conn.GetProperty(util.RoomId))
		} else {
			glog.Logger.Sugar().Errorf("not found clinet id %s", v)
		}
	}
}
func (g *GameClientProxyRouter) SetLogic(ctx *protocol.Context) {
	var logic string

	err := ctx.Bind(&logic)
	if err != nil {
		glog.Logger.Sugar().Errorf("SetLogic err:%s", err.Error())
		return
	}
	logicHash := util.MethodHash(logic)
	g.logicConnMgr.Add(logicHash, ctx.Conn)
	msg := &engine.InnerMsg{
		Properties: map[string]string{
			"gateIdx": strconv.Itoa(conf.GlobalServerConf.Idx),
		},
	}
	//buffer := protocol.Encode(g.server.Config.Idx, protocol.String, util.MethodHash("OnSetGateIdx"), util.LogicNone)
	glog.Logger.Sugar().Infof("set logic from game,logic:%s", logic)
	ctx.SendWithParams(msg, protocol.ProtoBuffer, util.MethodHash("OnSetGateIdx"))
}

func (g *GameClientProxyRouter) CallClient(ctx *protocol.Context) {
	msg := &engine.InnerMsg{}
	err := ctx.Bind(msg)
	if err != nil {
		glog.Logger.Sugar().Errorf("HandleGame err:%s", err.Error())
		return
	}
	buffer := protocol.Encode(msg.ClientMsg.Payload, protocol.CodeType(msg.ClientCodeType), msg.ClientMsg.Method, 0)
	//glog.Logger.Sugar().Infof("%s,%d", buffer.Bytes()[14:], protocol.CodeType(msg.ClientCodeType))
	g.connMgr.SendSomeone(buffer, msg.ClientIds, "CallClient")
}
