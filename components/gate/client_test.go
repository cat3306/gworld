package main

import (
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
		time.Sleep(1000 * time.Millisecond)
	}
}
func TestHeartBeat(t *testing.T) {
	_, conn1 := Conn()
	//select {}
	//go heartBeat(conn, false)
	go heartBeat(conn1, "GlobalHeartBeat")
	select {}
}
