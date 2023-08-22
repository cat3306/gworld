package engine

import (
	"github.com/cat3306/goworld/glog"
	"sync"

	"github.com/panjf2000/gnet/v2"
)

type ConnManager struct {
	connections map[string]gnet.Conn
	locker      sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[string]gnet.Conn),
	}
}
func (c *ConnManager) Add(conn gnet.Conn) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.connections[conn.ID()] = conn
}
func (c *ConnManager) Remove(id string) {
	c.locker.Lock()
	defer c.locker.Unlock()
	delete(c.connections, id)
}
func (c *ConnManager) Broadcast(raw []byte) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	for _, v := range c.connections {
		err := v.AsyncWrite(raw, nil)
		if err != nil {
			glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
		}
	}
}
func (c *ConnManager) BroadcastExceptSelf(raw []byte, cid string) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	for _, v := range c.connections {
		if v.ID() == cid {
			continue
		}
		err := v.AsyncWrite(raw, nil)
		if err != nil {
			glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
		}
	}
}
func (c *ConnManager) Len() int {
	c.locker.RLock()
	defer c.locker.RUnlock()
	return len(c.connections)
}
