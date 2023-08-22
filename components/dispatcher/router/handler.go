package router

import (
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
)

func SetDispatcherType(ctx *protocol.Context) {
	args := uint32(0)
	if err := ctx.Bind(&args); err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err.Error())
		return
	}
	if args == uint32(util.ClusterTypeGate) {
		v, _ := ctx.GetProperty(util.GateClientMgrKey)
		connMgr := v.(*engine.ConnManager)
		connMgr.Add(ctx.Conn)
		glog.Logger.Sugar().Infof("SetDispatcherType cid:%s", ctx.Conn.ID())
	} else {

	}
}
