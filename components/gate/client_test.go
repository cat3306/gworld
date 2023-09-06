package main

import (
	"encoding/json"
	"fmt"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
	"github.com/panjf2000/gnet/v2"
	"os"
	"testing"
	"time"
)

type MsgModel struct {
	engine.BaseRouter
}

func (h *MsgModel) Init(v interface{}) engine.IRouter {
	return h
}
func (h *MsgModel) HeartBeat(ctx *protocol.Context) {
	str := ""
	err := ctx.Bind(&str)
	glog.Logger.Sugar().Infof("HeartBeat:%s,id:%s", str, ctx.Conn.ID())
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
}

func (h *MsgModel) GlobalHeartBeat(ctx *protocol.Context) {
	s := ""
	err := ctx.Bind(&s)
	if err != nil {
		glog.Logger.Sugar().Errorf("GlobalHeartBeat:err:%s", err.Error())
	}
	glog.Logger.Sugar().Infof("GlobalHeartBeat:%s", s)

}
func (h *MsgModel) Auth(ctx *protocol.Context) {
	str := ""
	err := ctx.Bind(&str)
	glog.Logger.Sugar().Infof("Auth:%s,id:%s", str, ctx.Conn.ID())
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
}
func (h *MsgModel) CreateRoom(ctx *protocol.Context) {
	str := ""
	err := ctx.Bind(&str)
	glog.Logger.Sugar().Infof("Auth:%s,id:%s", str, ctx.Conn.ID())
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
}
func Conn() (gnet.Conn, gnet.Conn) {
	ev := engine.NewClientEvents(util.ClusterTypeGate)
	ev.AddRouter(new(MsgModel))
	cli, err := gnet.NewClient(ev)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	//con, err := cli.Dial("tcp", "127.0.0.1:8888")
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(0)
	//}
	err = cli.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	con1, err := cli.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return nil, con1
}
func init() {
	glog.Init()
}

func heartBeat(conn gnet.Conn, m string) {
	msg := &engine.ClientMsg{
		Logic:    util.MethodHash("base"),
		Payload:  []byte("💓"),
		Method:   util.MethodHash(m),
		CodeType: uint32(protocol.String),
	}
	raw := protocol.Encode(msg, protocol.ProtoBuffer, util.MethodHash("Dispatcher"))
	for {
		_, err := conn.Write(raw)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
}
func TestHeartBeat(t *testing.T) {
	_, conn1 := Conn()
	auth(conn1)
	//select {}
	//go heartBeat(conn, false)
	go heartBeat(conn1, "GlobalHeartBeat")
	select {}
}

func TestAuth(t *testing.T) {
	_, conn1 := Conn()
	auth(conn1)
	select {}

}
func auth(c gnet.Conn) {
	raw := protocol.Encode("1", protocol.String, util.MethodHash("Auth"))
	c.Write(raw)
}
func TestCreateRoom(t *testing.T) {
	_, conn1 := Conn()
	auth(conn1)
	time.Sleep(time.Second)
	createRoom(conn1)
	select {}
}

type CreateRoomReq struct {
	Pwd       string `json:"Pwd"`
	MaxNum    int    `json:"MaxNum"`    //最大人数
	JoinState bool   `json:"JoinState"` //是否能加入
}

func createRoom(conn gnet.Conn) {
	req := CreateRoomReq{
		Pwd:       "123",
		MaxNum:    10,
		JoinState: true,
	}
	reqR, _ := json.Marshal(req)
	msg := &engine.ClientMsg{
		Logic:    util.MethodHash("base"),
		Payload:  reqR,
		Method:   util.MethodHash("CreateRoom"),
		CodeType: uint32(protocol.Json),
	}
	raw := protocol.Encode(msg, protocol.ProtoBuffer, util.MethodHash("Dispatcher"))
	fmt.Println(conn.Write(raw))
}
