package main

import (
	"fmt"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/bytebufferpool"
	"time"
)

type GameServer struct {
	*engine.Server
}

func (g *GameServer) Run() {
	addr := fmt.Sprintf("tcp://:%d", g.Config.Port)
	f := func() {
		err := gnet.Run(g, addr,
			gnet.WithMulticore(true),
			gnet.WithSocketSendBuffer(g.Config.ConnWriteBuffer),
			gnet.WithSocketRecvBuffer(g.Config.ConnWriteBuffer),
			//gnet.WithTCPKeepAlive()
		)
		panic(err)
	}
	defer func() {
		g.HandlerMgr.GPool.Release()
	}()
	util.PanicRepeatRun(f, util.PanicRepeatRunArgs{
		Sleep: time.Second,
		Try:   20,
	})
}
func (g *GameServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	cId := util.GenConnId()
	c.SetId(cId)
	g.ConnMgr.Add(c)
	glog.Logger.Sugar().Infof("gate clinet conn cid:%s connect", c.ID())
	buffer := protocol.Encode(g.Config.Logic, protocol.String, util.MethodHash("SetLogic"), 0)
	out = buffer.Bytes()
	defer bytebufferpool.Put(buffer)
	return out, gnet.None
}
func (g *GameServer) OnTraffic(c gnet.Conn) gnet.Action {
	//defer func() {
	//	err := recover()
	//	if err != nil {
	//		glog.Logger.Sugar().Errorf("OnTraffic panic %v", err)
	//	}
	//}()
	//s.eng.CountConnections()
	context, err := protocol.Decode(c)
	if err != nil {
		glog.Logger.Sugar().Warnf("OnTraffic err:%s", err.Error())
		return gnet.None
	}
	if context == nil {
		panic("context nil")
	}
	innerMsg := &engine.InnerMsg{}
	err = context.Bind(innerMsg)
	if err != nil {
		glog.Logger.Sugar().Errorf("context.Bind err:%s", err.Error())
		return gnet.None
	}
	if len(innerMsg.ClientIds) == 0 {
		glog.Logger.Sugar().Errorf("client id none drop")
		return gnet.None
	}
	context.SetProperty(util.InnerMsgKey, innerMsg)
	g.ClientCtxChan <- context
	return gnet.None
}
