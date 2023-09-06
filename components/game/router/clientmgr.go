package router

import (
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
)

var (
	clientMgr *ClientMgr
)

//client
type GameClient struct {
	ConnId string
	RoomId string
}

type ClientMgr struct {
	clients map[string]*GameClient
}

func (c *ClientMgr) Init(v interface{}) engine.IRouter {
	c.clients = make(map[string]*GameClient)
	clientMgr = c
	return c
}
func (c *ClientMgr) OnConnect(ctx *protocol.Context) {
	msg := &engine.InnerMsg{}
	err := ctx.Bind(msg)
	if err != nil {
		glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
		return
	}
	if !engine.CheckInnerMsg(msg) {
		glog.Logger.Sugar().Errorf("msg check faied")
		return
	}
	glog.Logger.Sugar().Infof("client:%s connect", msg.ClientIds[0])
	cId := msg.ClientIds[0]
	_, ok := c.clients[cId]
	if ok {
		glog.Logger.Sugar().Warnf("%s:already add ,delete !", cId)
	}
	c.clients[cId] = &GameClient{
		ConnId: msg.ClientIds[0],
	}

}
func (c *ClientMgr) OnDisconnect(ctx *protocol.Context) {
	msg := &engine.InnerMsg{}
	err := ctx.Bind(msg)
	if err != nil {
		glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
		return
	}
	if !engine.CheckInnerMsg(msg) {
		glog.Logger.Sugar().Errorf("msg check faied")
		return
	}
	glog.Logger.Sugar().Infof("client:%s disconnect", msg.ClientIds[0])
	delete(c.clients, msg.ClientIds[0])
}
func (c *ClientMgr) GetInfo(id string) (*GameClient, bool) {
	v, o := c.clients[id]
	return v, o
}

func (c *ClientMgr) ObjectMove(ctx *protocol.Context) {
	msg := &engine.InnerMsg{}
	err := ctx.Bind(msg)
	if err != nil {
		glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
		return
	}
	glog.Logger.Sugar().Infof("%s",string(msg.ClientMsg.Payload))
}
