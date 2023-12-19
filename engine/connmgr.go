package engine

import (
	"github.com/cat3306/goworld/glog"
	"github.com/valyala/bytebufferpool"
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
	c.connections[conn.Fd()] = conn
}
func (c *ConnManager) Remove(id string) {
	c.locker.Lock()
	defer c.locker.Unlock()
	delete(c.connections, id)
}
func (c *ConnManager) Get(id string) (bool, gnet.Conn) {
	c.locker.Lock()
	defer c.locker.Unlock()
	con, ok := c.connections[id]
	return ok, con
}
func (c *ConnManager) Broadcast(buffer *bytebufferpool.ByteBuffer, bufferIsPut bool) {
	c.locker.RLock()
	defer func() {
		c.locker.RUnlock()
		if bufferIsPut {
			bytebufferpool.Put(buffer)
		}
	}()
	for _, v := range c.connections {
		_, err := v.Write(buffer.Bytes())
		if err != nil {
			glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
		}
	}
}
func (c *ConnManager) SendOne(buffer *bytebufferpool.ByteBuffer, id string) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	if conn, ok := c.connections[id]; ok {
		err := conn.AsyncWrite(buffer.Bytes(), func(c gnet.Conn) error {
			bytebufferpool.Put(buffer)
			return nil
		})
		if err != nil {
			glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
		}
	} else {
		glog.Logger.Sugar().Errorf("not found conn:%s", id)
	}
}
func (c *ConnManager) SendSomeone(buffer *bytebufferpool.ByteBuffer, ids []string, args string) {
	c.locker.RLock()
	defer func() {
		c.locker.RUnlock()
		bytebufferpool.Put(buffer)
	}()
	for _, id := range ids {
		if conn, ok := c.connections[id]; ok {
			_, err := conn.Write(buffer.Bytes())
			if err != nil {
				glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
			}
		} else {
			glog.Logger.Sugar().Errorf("%s not found conn:%s", args, id)
		}
	}

}
func (c *ConnManager) BroadcastExceptSelf(buffer *bytebufferpool.ByteBuffer, cid string) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	for _, v := range c.connections {
		if v.ID() == cid {
			continue
		}
		_, err := v.Write(buffer.Bytes())
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

func (c *ConnManager) SendRandOne(buffer *bytebufferpool.ByteBuffer) {
	c.locker.RLock()
	defer func() {
		c.locker.RUnlock()
		bytebufferpool.Put(buffer)
	}()
	for _, conn := range c.connections {
		_, err := conn.Write(buffer.Bytes())
		if err != nil {
			glog.Logger.Sugar().Errorf("Write err:%s", err.Error())
		}
		return
	}
}
