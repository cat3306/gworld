package router

import (
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
)

type ServerInner struct {
	engine.BaseRouter
}

func (h *ServerInner) Init(v interface{}) engine.IRouter {
	return h
}
func (h *ServerInner) SetDispatcherType(ctx *protocol.Context) {
	args := uint32(0)
	if err := ctx.Bind(&args); err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err.Error())
		return
	}
	var connMgr *engine.ConnManager
	if args == uint32(util.ClusterTypeGate) {
		v, _ := ctx.GetProperty(util.GateClientMgrKey)
		connMgr = v.(*engine.ConnManager)
	} else {
		v, _ := ctx.GetProperty(util.GameClientMgrKey)
		connMgr = v.(*engine.ConnManager)
	}
	glog.Logger.Sugar().Infof("SetDispatcherType cid:%s,pro:%d", ctx.Conn.ID(), args)
	connMgr.Add(ctx.Conn)
}
