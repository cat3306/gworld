package engine

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/bytebufferpool"
)

type LogicConnMgr struct {
	logicConnMgr map[uint32]*ConnManager
}

func NewLogicConnMgr() *LogicConnMgr {
	return &LogicConnMgr{
		logicConnMgr: map[uint32]*ConnManager{},
	}
}
func (l *LogicConnMgr) Add(logic uint32, c gnet.Conn) {
	conMgr, ok := l.logicConnMgr[logic]
	if !ok {
		conMgr = NewConnManager()
		l.logicConnMgr[logic] = conMgr
	}
	conMgr.Add(c)
}

func (l *LogicConnMgr) Remove(id string) {
	for _, v := range l.logicConnMgr {
		ok, _ := v.Get(id)
		if ok {
			v.Remove(id)
			break
		}
	}
}
func (l *LogicConnMgr) GetByLogic(logic uint32) (*ConnManager, bool) {
	mgr, ok := l.logicConnMgr[logic]
	return mgr, ok
}

func (l *LogicConnMgr) Broadcast(buffer *bytebufferpool.ByteBuffer) {
	for _, v := range l.logicConnMgr {
		v.Broadcast(buffer)
	}
}
