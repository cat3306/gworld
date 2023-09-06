package engine

import (
	"github.com/panjf2000/gnet/v2"
	"sync"
)

type LogicConnMgr struct {
	logicConnMgr map[uint32]*ConnManager
	lock         sync.RWMutex
}

func NewLogicConnMgr() *LogicConnMgr {
	return &LogicConnMgr{
		logicConnMgr: map[uint32]*ConnManager{},
		lock:         sync.RWMutex{},
	}
}
func (l *LogicConnMgr) Add(logic uint32, c gnet.Conn) {
	l.lock.Lock()
	defer l.lock.Unlock()
	conMgr, ok := l.logicConnMgr[logic]
	if !ok {
		conMgr = NewConnManager()
		l.logicConnMgr[logic] = conMgr
	}
	conMgr.Add(c)
}

func (l *LogicConnMgr) Remove(id string) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	for _, v := range l.logicConnMgr {
		ok, _ := v.Get(id)
		if ok {
			v.Remove(id)
			break
		}
	}
}
func (l *LogicConnMgr) GetByLogic(logic uint32) (*ConnManager, bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	mgr, ok := l.logicConnMgr[logic]
	return mgr, ok
}

func (l *LogicConnMgr) Broadcast(raw []byte) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	for _, v := range l.logicConnMgr {
		v.Broadcast(raw)
	}
}
