package router

import (
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
)

type RoomMgr struct {
	engine.BaseRouter
	rooms map[string]*Room
}

func (r *RoomMgr) Init(v interface{}) engine.IRouter {
	r.rooms = make(map[string]*Room)
	return r
}

type CreateRoomReq struct {
	Pwd       string `json:"Pwd"`
	MaxNum    int    `json:"MaxNum"`    //最大人数
	JoinState bool   `json:"JoinState"` //是否能加入
}
type CreateRoomRsp struct {
	Id string `json:"id"`
}

// APi CreateRoom
func (r *RoomMgr) CreateRoom(ctx *protocol.Context) {
	msg := &engine.InnerMsg{}
	err := ctx.Bind(msg)
	if err != nil {
		glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
		return
	}
	coder := protocol.GameCoder(protocol.CodeType(msg.ClientMsg.CodeType))
	req := CreateRoomReq{}
	err = coder.Unmarshal(msg.ClientMsg.Payload, &req)
	if err != nil {
		glog.Logger.Sugar().Errorf("coder.Unmarshal err:%s", err.Error())
		return
	}
	glog.Logger.Sugar().Infof("req:%+v", req)
	if req.MaxNum == 0 {
		req.MaxNum = 1
	}
	room := &Room{
		maxNum:    req.MaxNum,
		pwd:       req.Pwd,
		joinState: req.JoinState,
		gameState: false,
		scene:     0,
		Id:        util.GenId(7),
		clients:   map[string]*ClientInfo{},
	}
	room.clients[msg.ClientIds[0]] = &ClientInfo{}
	r.AddRoom(room)
	iMsg := &engine.InnerMsg{
		ClientIds: msg.ClientIds,
		ClientMsg: &engine.ClientMsg{
			Logic:    0,
			Payload:  []byte(room.Id),
			Method:   msg.ClientMsg.Method,
			CodeType: uint32(protocol.String),
		},
	}
	ctx.SendWithParams(iMsg, protocol.ProtoBuffer, util.CallClient)

}
func (r *RoomMgr) AddRoom(room *Room) {
	r.rooms[room.Id] = room
}
func (r *RoomMgr) DelRoom(id string) {
	delete(r.rooms, id)
}
func (r *RoomMgr) GetRoom(id string) (*Room, bool) {
	room, ok := r.rooms[id]
	return room, ok
}
func (r *RoomMgr) LeaveRoomByConnClose(roomId string, connId string) {

}

// API LeaveRoom
func (r *RoomMgr) LeaveRoom(ctx *protocol.Context) {
	////roomId := ctx.GetRoomId()
	//room, _ := r.GetRoom(roomId)
	//if room == nil {
	//	ctx.SendWithCodeType(JsonRspErr(fmt.Sprintf("room not found,id:%s", roomId)), protocol.Json)
	//	glog.Logger.Sugar().Errorf("room not found roomId:%s", roomId)
	//	return
	//}
	//
	//room.connMgr.Remove(ctx.Conn.ID())
	//if room.connMgr.Len() == 0 {
	//	r.DelRoom(roomId)
	//}
	//ctx.DelRoomId()
	//glog.Logger.Sugar().Infof("ok,roomId:%s,clientId:%s", roomId, ctx.Conn.ID())
	//
	//ctx.SendWithCodeType(JsonRspOK(""), protocol.Json)
}

type JoinRoomReq struct {
	RoomId string `json:"RoomId"`
	Pwd    string `json:"Pwd"`
}

//
func (r *RoomMgr) JoinRoom(ctx *protocol.Context) {
	//req := &JoinRoomReq{}
	//if err := ctx.Bind(req); err != nil {
	//	glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
	//	ctx.Send(JsonRspErr(err.Error()))
	//	return
	//}
	//room, err := r.joinRoom(req, ctx)
	//if err != nil {
	//	glog.Logger.Sugar().Errorf("JoinRoom err:%s,req:%+v", err.Error(), req)
	//	ctx.Send(JsonRspErr(err.Error()))
	//	return
	//}

	//room.Broadcast(JsonRspOK(""), ctx)
	//ctx.Send()
}
func (r *RoomMgr) joinRoom(req *JoinRoomReq, ctx *protocol.Context) (*Room, error) {
	//room, _ := r.GetRoom(req.RoomId)
	//if room == nil {
	//	return nil, errors.New("room not found")
	//}
	//if !room.joinState {
	//	return nil, errors.New("room not allow join")
	//}
	//if room.pwd != "" {
	//	if req.Pwd != room.pwd {
	//		return nil, errors.New("room pwd not correct")
	//	}
	//}
	//if room.connMgr.Len() >= room.maxNum {
	//	return nil, errors.New("room full")
	//}
	//room.connMgr.Add(ctx.Conn)
	//return room, nil

	return nil, nil
}
