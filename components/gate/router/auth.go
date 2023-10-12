package router

import (
	"github.com/cat3306/gocommon/cryptoutil"
	"github.com/cat3306/goworld/conf"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
)

type Auth struct {
	engine.BaseRouter
	rawPrivateKey []byte
}

func (h *Auth) Init(v ...interface{}) engine.IRouter {
	privateKeyRaw, err := cryptoutil.RawRSAKey(conf.GlobalConf.AuthConfig.PrivateKeyPath)
	if err != nil {
		panic(err)
	}
	h.rawPrivateKey = privateKeyRaw
	return h
}

type AuthReq struct {
	CipherText []byte `json:"CipherText"`
	Text       string `json:"Text"`
}

func (h *Auth) Auth(ctx *protocol.Context) {
	req := AuthReq{}
	err := ctx.Bind(&req)
	if err != nil {
		glog.Logger.Sugar().Errorf("param err:%s", err.Error())
		return
	}
	glog.Logger.Sugar().Infof("cid:%s auth", ctx.Conn.ID())
	text := cryptoutil.RsaDecrypt(req.CipherText, h.rawPrivateKey)
	if string(text) != req.Text {
		glog.Logger.Sugar().Errorf("认证失败!")
		ctx.SendWithParams("认证失败", protocol.String, ctx.Proto)
	} else {
		ctx.Conn.SetProperty(util.ClientAuth, "ok")
		ctx.SendWithParams("ok", protocol.String, ctx.Proto)
	}
}
