package engine

import (
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
)

func PreSetInnerMsgMsg() Handler {
	return func(ctx *protocol.Context) {
		innerMsg := &InnerMsg{}
		err := ctx.Bind(innerMsg)
		if err != nil {
			glog.Logger.Sugar().Errorf("context.Bind err:%s", err.Error())
			return
		}
		ctx.SetProperty(util.InnerMsgKey, innerMsg)
	}
}
