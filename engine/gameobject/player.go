package gameobject

import (
	"encoding/json"
	"github.com/cat3306/goworld/glog"
)

type Player struct {
	GameObject
}

func (p *Player) OnSaveData() []byte {
	raw, err := json.Marshal(p)
	if err != nil {
		glog.Logger.Sugar().Errorf("OnSaveData err:%s", err.Error())
	}
	return raw
}
