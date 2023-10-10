package router

import (
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/bytebufferpool"
)

var (
	clientMgr *ClientMgr
)

//client
type GameClient struct {
	ClientId     string
	ServerConnId string
	RoomId       string
	UserId       string
}

type ClientMgr struct {
	clients     map[string]*GameClient
	gateConnMgr *engine.ConnManager
	gateClients map[string][]*GameClient
}

func (c *ClientMgr) Init(v interface{}) engine.IRouter {
	c.clients = make(map[string]*GameClient)
	clientMgr = c
	mgr := v.(*engine.ConnManager)
	c.gateConnMgr = mgr
	return c
}
func (c *ClientMgr) OnConnect(ctx *protocol.Context) {
	msg, err := engine.GetCtxInnerMsg(ctx)
	if err != nil {
		glog.Logger.Sugar().Errorf("OnConnect err:%s", err.Error())
		return
	}
	glog.Logger.Sugar().Infof("client:%s connect", msg.ClientIds[0])
	cId := msg.ClientIds[0]
	_, ok := c.clients[cId]
	if ok {
		glog.Logger.Sugar().Warnf("%s:already add ,delete !", cId)
	}
	c.clients[cId] = &GameClient{
		ClientId:     msg.ClientIds[0],
		ServerConnId: ctx.Conn.ID(),
	}

}
func (c *ClientMgr) OnDisconnect(ctx *protocol.Context) {
	msg, err := engine.GetCtxInnerMsg(ctx)
	if err != nil {
		glog.Logger.Sugar().Errorf("OnDisconnect err:%s", err.Error())
		return
	}
	glog.Logger.Sugar().Infof("client:%s disconnect", msg.ClientIds[0])
	delete(c.clients, msg.ClientIds[0])
}
func (c *ClientMgr) GetInfo(id string) (*GameClient, bool) {
	v, o := c.clients[id]
	return v, o
}

func (c *ClientMgr) OnSetGateIdx(ctx *protocol.Context) {
	msg, err := engine.GetCtxInnerMsg(ctx)
	if err != nil {
		glog.Logger.Sugar().Errorf("err:%s", err.Error())
		return
	}
	idx := msg.Properties["gateIdx"]
	if idx == "" {
		glog.Logger.Sugar().Errorf("invalid gate idx")
		return
	}
	ctx.Conn.SetId(idx)
	c.gateConnMgr.Add(ctx.Conn)

}
func (c *ClientMgr) Broadcast(ctx *protocol.Context, clientIds []string, object interface{}) {

	gateIds := map[string][]string{}
	if len(clientIds) == 0 {
		for _, info := range c.clients {
			if tmp, ok := gateIds[info.ServerConnId]; ok {
				gateIds[info.ServerConnId] = append(tmp, info.ClientId)
			} else {
				gateIds[info.ServerConnId] = []string{info.ClientId}
			}
		}
	} else {
		for cId, v := range c.clients {
			if tmp, ok := gateIds[v.ServerConnId]; ok {
				gateIds[v.ServerConnId] = append(tmp, cId)
			} else {
				gateIds[v.ServerConnId] = []string{cId}
			}
		}
	}
	msg, _ := engine.GetCtxInnerMsg(ctx)
	switch object.(type) {
	case string:
		msg.ClientMsg.Payload = util.StringToBytes(object.(string))
		msg.ClientCodeType = uint32(protocol.String)
	case []byte:
		msg.ClientMsg.Payload = object.([]byte)
	default:
		raw, err := protocol.GameCoder(protocol.CodeType(msg.ClientCodeType)).Marshal(object)
		if err != nil {
			panic(err)
		}
		msg.ClientMsg.Payload = raw
	}
	glog.Logger.Sugar().Infof("%+v", gateIds)
	for v, cIds := range gateIds {
		msg.ClientIds = cIds
		buffer := protocol.Encode(msg, protocol.ProtoBuffer, util.CallClient, ctx.Logic)
		_, conn := c.gateConnMgr.Get(v)
		err := conn.AsyncWrite(buffer.Bytes(), func(c gnet.Conn) error {
			bytebufferpool.Put(buffer)
			return nil
		})
		if err != nil {
			glog.Logger.Sugar().Errorf("Broadcast err:%s", err.Error())
		}

	}

}
