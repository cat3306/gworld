package protocol

import (
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/util"
	"github.com/valyala/bytebufferpool"
	"sync"

	"github.com/panjf2000/gnet/v2"
)

type Context struct {
	Payload  *bytebufferpool.ByteBuffer
	CodeType CodeType
	Proto    uint32
	Conn     gnet.Conn
	Logic    uint32
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
	defer bytebufferpool.Put(c.Payload)
	return GameCoder(c.CodeType).Unmarshal(c.Payload.Bytes(), v)
}

func (c *Context) Send(v interface{}) {
	buffer := Encode(v, c.CodeType, c.Proto, c.Logic)
	err := c.Conn.AsyncWrite(buffer.Bytes(), func(c gnet.Conn) error {
		bytebufferpool.Put(buffer)
		return nil
	})
	if err != nil {
		glog.Logger.Sugar().Errorf("ctx Send err:%s", err.Error())
	}
}

func (c *Context) SendWithCodeType(v interface{}, codeType CodeType) {
	buffer := Encode(v, codeType, c.Proto, c.Logic)
	err := c.Conn.AsyncWrite(buffer.Bytes(), func(c gnet.Conn) error {
		bytebufferpool.Put(buffer)
		return nil
	})
	if err != nil {
		glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
	}
}
func (c *Context) SendWithParams(v interface{}, codeType CodeType, hash uint32) {
	buffer := Encode(v, codeType, hash, c.Logic)
	err := c.Conn.AsyncWrite(buffer.Bytes(), func(c gnet.Conn) error {
		bytebufferpool.Put(buffer)
		return nil
	})
	if err != nil {
		glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
	}
}

func (c *Context) CheckClientAuth() bool {
	ok := c.Conn.GetProperty(util.ClientAuth)
	return ok != ""
}
