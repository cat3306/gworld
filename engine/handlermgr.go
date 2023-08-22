package engine

import (
	"errors"
	"fmt"
	"github.com/cat3306/goworld/glog"
	"github.com/cat3306/goworld/protocol"
	"github.com/cat3306/goworld/util"
	"reflect"

	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
)

var (
	ErrHandlerNotFound = errors.New("handler not found")
)

type Handler func(c *protocol.Context)
type GoHandler func(c *protocol.Context, none struct{})

func NewHandlerManager() *HandlerManager {
	return &HandlerManager{
		handlers:  make(map[uint32]Handler),
		goHandler: make(map[uint32]GoHandler),
		GPool:     goroutine.Default(),
	}
}

type HandlerManager struct {
	handlers  map[uint32]Handler
	goHandler map[uint32]GoHandler
	GPool     *goroutine.Pool
}

func (h *HandlerManager) Register(hashCode uint32, handler Handler) {
	if _, ok := h.handlers[hashCode]; ok {
		panic(fmt.Sprintf("Register repeated method:%d", hashCode))
	}
	h.handlers[hashCode] = handler
}
func (h *HandlerManager) GoRegister(hashCode uint32, handler GoHandler) {
	if _, ok := h.goHandler[hashCode]; ok {
		panic(fmt.Sprintf("Register repeated method:%d", hashCode))
	}
	h.goHandler[hashCode] = handler
}
func (h *HandlerManager) RegisterRouter(iG IRouter) {
	t := reflect.TypeOf(iG)
	tName := t.String()
	vl := reflect.ValueOf(iG)
	for i := 0; i < t.NumMethod(); i++ {
		name := t.Method(i).Name
		v, ok := vl.Method(i).Interface().(func(ctx *protocol.Context))
		if ok {
			if checkoutMethod(name) {
				hashId := util.MethodHash(name)
				h.Register(hashId, v)
				glog.Logger.Sugar().Infof("[%s.%s] hashId:%d", tName, name, hashId)
			}
		}
		v1, ok1 := vl.Method(i).Interface().(func(c *protocol.Context, none struct{}))
		if ok1 {
			if checkoutMethod(name) {
				hashId := util.MethodHash(name)
				h.GoRegister(hashId, v1)
				glog.Logger.Sugar().Infof("[%s.go_%s] hashId:%d", tName, name, hashId)
			}
		}

	}
}

//函数签名首字母大写才会被注入
func checkoutMethod(m string) bool {
	if len(m) == 0 {
		return false
	}
	if m[0] >= 'A' && m[0] <= 'W' {
		return true
	}
	return false
}
func (h *HandlerManager) GetHandler(proto uint32) Handler {
	f := h.handlers[proto]
	return f
}
func (h *HandlerManager) GetGoHandler(proto uint32) GoHandler {
	f := h.goHandler[proto]
	return f
}

//同步handler
func (h *HandlerManager) exeSyncHandler(ctx *protocol.Context) error {
	f := h.GetHandler(ctx.Proto)
	if f != nil {
		f(ctx)
		return nil
	}
	return ErrHandlerNotFound
}

//异步handler
func (h *HandlerManager) exeAsyncHandler(ctx *protocol.Context) error {
	f := h.GetGoHandler(ctx.Proto)
	if f != nil {
		newBuffer := protocol.BUFFERPOOL.Get(uint32(len(ctx.Payload)))
		copy(*newBuffer, ctx.Payload)
		ctx.Payload = *newBuffer
		err := h.GPool.Submit(func() {
			f(ctx, struct{}{})
		})
		if err != nil {
			glog.Logger.Sugar().Errorf("exeGoHandler err:%s", err.Error())
			return err
		}
		return nil
	}
	return ErrHandlerNotFound
}

func (h *HandlerManager) ExeHandler(ctx *protocol.Context) {
	err := h.exeSyncHandler(ctx)
	if !errors.Is(err, ErrHandlerNotFound) {
		return
	}
	err = h.exeAsyncHandler(ctx)
	if err != nil {
		glog.Logger.Sugar().Errorf("ExeHandler err:%s,pro:%d", err.Error(), ctx.Proto)
	}
}
