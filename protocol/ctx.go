package protocol

import (
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/util"
	"sync"

	"github.com/panjf2000/gnet/v2"
)

type Context struct {
	Payload  []byte
	CodeType CodeType
	Proto    uint32
	Conn     gnet.Conn
	property sync.Map
}

func (c *Context) SetProperty(k interface{}, v interface{}) {
	c.property.Store(k, v)
}
func (c *Context) GetProperty(k interface{}) (interface{}, bool) {
	return c.property.Load(k)
}
func (c *Context) DelProperty(k string) {
	c.property.Delete(k)
}
func (c *Context) Bind(v interface{}) error {

	return GameCoder(c.CodeType).Unmarshal(c.Payload, v)
}

func (c *Context) Send(v interface{}) {
	err := c.AsyncWrite(Encode(v, c.CodeType, c.Proto))
	if err != nil {
		glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
	}
}
func (c *Context) SendWithCodeType(v interface{}, codeType CodeType) {
	err := c.AsyncWrite(Encode(v, codeType, c.Proto))
	if err != nil {
		glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
	}
}
func (c *Context) SendWithParams(v interface{}, codeType CodeType, hash uint32) {
	err := c.AsyncWrite(Encode(v, codeType, hash))
	if err != nil {
		glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
	}
}

func (c *Context) AsyncWrite(raw []byte) error {
	return c.Conn.AsyncWrite(raw, func(c gnet.Conn) error {
		BUFFERPOOL.Put(raw)
		return nil
	})
}

func (c *Context) CheckClientAuth() bool {
	ok := c.Conn.GetProperty(util.ClientAuth)
	return ok != ""
}
