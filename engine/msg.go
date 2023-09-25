package engine

import (
	"errors"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
)

type TryConnectMsg struct {
	NetWork string
	Addr    string
}

func GameBind(v interface{}, ctx *protocol.Context) (*InnerMsg, error) {
	msg, err := GetCtxInnerMsg(ctx)
	if err != nil {
		return nil, err
	}
	return msg, protocol.GameCoder(protocol.CodeType(msg.ClientCodeType)).Unmarshal(msg.ClientMsg.Payload, v)
}

func GetCtxInnerMsg(ctx *protocol.Context) (*InnerMsg, error) {
	v, ok := ctx.GetProperty(util.InnerMsgKey)
	if !ok {
		return nil, errors.New("not found ContextInnerMsgKey")
	}
	msg, ok1 := v.(*InnerMsg)
	if !ok1 {
		return nil, errors.New("v not *InnerMsg")
	}
	return msg, nil
}
func GameSendSelf(ctx *protocol.Context, v interface{}) {
	msg, _ := GetCtxInnerMsg(ctx)
	switch v.(type) {
	case string:
		msg.ClientMsg.Payload = util.StringToBytes(v.(string))
		msg.ClientCodeType = uint32(protocol.String)
	case []byte:
		msg.ClientMsg.Payload = v.([]byte)
	default:
		raw, err := protocol.GameCoder(protocol.CodeType(msg.ClientCodeType)).Marshal(v)
		if err != nil {
			panic(err)
		}
		msg.ClientMsg.Payload = raw
	}
	ctx.SendWithParams(msg, protocol.ProtoBuffer, util.CallClient)
}
func GameSendSelfWithCodeType(ctx *protocol.Context, v interface{}, codeType protocol.CodeType) {
	msg, _ := GetCtxInnerMsg(ctx)
	switch v.(type) {
	case string:
		msg.ClientMsg.Payload = util.StringToBytes(v.(string))
		msg.ClientCodeType = uint32(protocol.String)
	case []byte:
		msg.ClientMsg.Payload = v.([]byte)
	default:
		raw, err := protocol.GameCoder(codeType).Marshal(v)
		if err != nil {
			panic(err)
		}
		msg.ClientMsg.Payload = raw
	}
	ctx.SendWithParams(msg, protocol.ProtoBuffer, util.CallClient)
}
func GameBroadcast(ctx *protocol.Context, v interface{}, clientIds []string) {
	msg, _ := GetCtxInnerMsg(ctx)
	switch v.(type) {
	case string:
		msg.ClientMsg.Payload = util.StringToBytes(v.(string))
		msg.ClientCodeType = uint32(protocol.String)
	case []byte:
		msg.ClientMsg.Payload = v.([]byte)
	default:
		raw, err := protocol.GameCoder(protocol.CodeType(msg.ClientCodeType)).Marshal(v)
		if err != nil {
			panic(err)
		}
		msg.ClientMsg.Payload = raw
	}
	msg.ClientIds = clientIds
	ctx.SendWithParams(msg, protocol.ProtoBuffer, util.CallClient)
}
