package handler

import (
	itconf "github.com/cat3306/goworld/components/gateweb/conf"
	"github.com/cat3306/goworld/conf"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/util"
	"github.com/gin-gonic/gin"
)

type GateInfoReq struct {
	Name string `json:"name"`
}

func GateInfo(c *gin.Context) {
	req := GateInfoReq{}
	if err := c.BindJSON(&req); err != nil {
		glog.Logger.Sugar().Errorf("GateInfo params invalid,err:%s", err.Error())
		RspError(c, "params invalid")
		return
	}
	rsp, err := gateInfo(&req)
	if err != nil {
		glog.Logger.Sugar().Errorf("GateInfo err:%s,req:%+v", err.Error(), req)
		RspError(c, err.Error())
		return
	}
	RspOk(c, rsp)
}

type gateInfoData struct {
	Port int    `json:"port"`
	Ip   string `json:"ip"`
	Name string `json:"name"`
}

func gateInfo(req *GateInfoReq) (interface{}, error) {
	err := conf.LoadConf(itconf.Global.ClusterPath, req.Name)
	if err != nil {
		return nil, err
	}
	rsp := make([]gateInfoData, 0)
	for k, servers := range conf.GlobalConf.Servers {
		if k == util.ClusterTypeGate {
			for _, vv := range servers {
				if !vv.Online {
					continue
				}
				rsp = append(rsp, gateInfoData{
					Port: vv.Port,
					Ip:   vv.OuterIp,
					Name: vv.Name,
				})
			}

		}
	}
	glog.Logger.Sugar().Infof("%+v", rsp)
	return rsp, nil
}
