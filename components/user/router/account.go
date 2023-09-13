package router

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"

	"github.com/cat3306/goworld/components/user/models"
	"github.com/cat3306/goworld/components/user/thirdmodule"
	"github.com/cat3306/goworld/engine"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"time"
)

type Account struct {
	engine.BaseRouter
	salt   string
	secret string
}

const (
	CodeTimeOut   = 3 * time.Minute
	GenSignUpCode = "GenSignUpCode"
	Login         = ""
)

func (a *Account) Init(v interface{}) engine.IRouter {
	//a.salt = conf.GameConfig.Account.Salt
	return a
}
func (a *Account) HashPwd(pwd string) string {
	h := sha256.New()
	h.Write([]byte(a.salt))
	h.Write([]byte("$"))
	h.Write([]byte(pwd))
	pre := "sha256"
	return fmt.Sprintf("%s:%s", pre, hex.EncodeToString(h.Sum(nil)))
}
func (a *Account) ComparePwd(reqPwd string, hashPwd string) bool {
	s := a.HashPwd(reqPwd)
	return s == hashPwd
}

type LoginReq struct {
	Email string `json:"Email"`
	Pwd   string `json:"Pwd"`
}
type LoginRsp struct {
	UserId string `json:"UserId"`
	Nick   string `json:"Nick"`
}

func (a *Account) Login(ctx *protocol.Context) {
	req := &LoginReq{}
	innerMsg, err := engine.GameBind(req, ctx)
	if err != nil {
		glog.Logger.Sugar().Errorf("Login err:%s,innerMsg:%+v", err.Error(), innerMsg)
	}
	rsp, err := a.login(req)
	if err != nil {
		glog.Logger.Sugar().Errorf("Login err:%s,req:%+v", err.Error(), req)
		engine.GameSendSelf(ctx, engine.JsonRspErr(err.Error()))
		return
	}
	engine.GameSendSelf(ctx, engine.JsonRspOK(rsp))
}
func (a *Account) login(req *LoginReq) (interface{}, error) {
	if req.Email == "" || req.Pwd == "" {
		return nil, errors.New("invalid args")
	}
	userProfile := models.UserProfile{}
	err := userProfile.GetByEmail(thirdmodule.MysqlDb, req.Email)
	if err != nil {
		return nil, err
	}
	if userProfile.Id == 0 {
		return nil, errors.New("你还没注册，亲亲")
	}
	if !a.ComparePwd(req.Pwd, userProfile.Pwd) {
		return nil, errors.New("密码错误")
	}
	var okInt int64
	okInt, err = thirdmodule.CacheSelect(0).Exists(context.Background(), userProfile.UserId).Result()
	if err != nil {
		return nil, err
	}
	if okInt == 1 {
		glog.Logger.Sugar().Errorf("already login ")
		return nil, errors.New("已登录")
	}
	//	ctx.DelUserId()
	rsp := LoginRsp{
		UserId: userProfile.UserId,
		Nick:   userProfile.NickName,
	}
	data := map[string]string{
		"email": userProfile.Email,
		"nick":  userProfile.NickName,
	}
	info, _ := json.Marshal(data)
	err = thirdmodule.CacheSelect(0).Set(context.Background(), rsp.UserId, string(info), -1).Err()
	if err != nil {
		glog.Logger.Sugar().Errorf("Cache.Set err:%s,req:%+v", err.Error(), req)
		return nil, err
	}
	return rsp, nil
}
func (a *Account) Logout(ctx *protocol.Context) {
	//err := thirdmodule.Cache.Del(ctx.GetUserId()).Err()
	//if err != nil {
	//	glog.Logger.Sugar().Errorf("Cache.Del err:%s", err.Error())
	//	ctx.SendWithCodeType(gnet.JsonRspErr("注销错误"), gface.Json)
	//}
	//ctx.DelUserId()
	ctx.Send("OK")
}
func (a *Account) Register(ctx *protocol.Context, tag struct{}) {
	type Req struct {
		Email string `json:"Email"`
		Pwd   string `json:"Pwd"`
		Nick  string `json:"Nick"`
		Code  string `json:"Code"`
	}
	req := &Req{}
	if err := ctx.Bind(req); err != nil {
		glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
		//ctx.Send(gnet.JsonRspErr(err.Error()))
		return
	}
	stashCode, err := thirdmodule.CacheSelect(0).Get(context.Background(),GenSignUpCode + req.Email).Result()
	if err != nil {
		glog.Logger.Sugar().Errorf("thirdmodule.Cache.GetString failed,err:%s", err.Error())
		//ctx.Send(gnet.JsonRspErr("获取验证码失败!"))
		return
	}
	if stashCode != req.Code {
		glog.Logger.Sugar().Errorf("code wrong !,email:%s", req.Email)
		//ctx.Send(gnet.JsonRspErr("验证码错误"))
		return
	}
	user := models.UserProfile{}
	err = user.GetByEmail(thirdmodule.MysqlDb, req.Email)
	if err != nil {
		glog.Logger.Sugar().Errorf("user.GetByEmail failed to %s", err.Error())
		//ctx.Send(gnet.JsonRspErr(err.Error()))
		return
	}
	if user.Id != 0 {
		glog.Logger.Sugar().Errorf("已注册,email:%s", req.Email)
		//ctx.Send(gnet.JsonRspErr("已注册"))
		return
	}
	uid := uuid.New()
	user = models.UserProfile{
		NickName:   req.Nick,
		Pwd:        a.HashPwd(req.Pwd),
		Email:      req.Email,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		UserId:     uid.String(),
	}
	if err = user.Create(thirdmodule.MysqlDb); err != nil {
		glog.Logger.Sugar().Errorf(" user.Create failed to,err:%s", err.Error())
		//ctx.Send(gnet.JsonRspErr(err.Error()))
		return
	}
	//ctx.Send(gnet.JsonRspOK("注册成功!"))
}
func (a *Account) EmailCode(ctx *protocol.Context, tag struct{}) {
	type Req struct {
		Email string `json:"Email"`
	}
	req := &Req{}
	if err := ctx.Bind(req); err != nil {
		glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
		//ctx.Send(gnet.JsonRspErr(err.Error()))
		return
	}
	user := models.UserProfile{}
	if err := user.GetByEmail(thirdmodule.MysqlDb, req.Email); err != nil {
		glog.Logger.Sugar().Errorf("user.GetByEmail err:%s,req:%+v", err.Error(), req)
		//ctx.Send(gnet.JsonRspErr(err.Error()))
		return
	}
	if user.Id != 0 {
		//ctx.Send(gnet.JsonRspErr("该邮箱已注册"))
		glog.Logger.Sugar().Errorf("email:%s,already register", req.Email)
		return
	}
	code := GenEmailCode(6)
	err := thirdmodule.CacheSelect(0).Set(context.Background(),GenSignUpCode+req.Email, code, CodeTimeOut).Err()
	if err != nil {
		glog.Logger.Sugar().Errorf("configure.Cache.Set failed err:%s", err.Error())
		//ctx.Send(gnet.JsonRspErr(err.Error()))
		return
	}
	to := []string{fmt.Sprintf("%v", req.Email)}
	title := "【恶魔射手】"
	content := fmt.Sprintf("亲爱的用户,您的验证码为%v,有效时间%d分钟,祝您游戏愉快！", code, CodeTimeOut/60)
	from := "恶魔射手官网"
	err = EmailSend(to, title, content, from)
	if err != nil {
		glog.Logger.Sugar().Errorf(" email.SendEmail,err:%s", err.Error())
		//ctx.Send(gnet.JsonRspErr(err.Error()))
		return
	}
	//ctx.Send(gnet.JsonRspOK("发送成功，注意查收"))
}
