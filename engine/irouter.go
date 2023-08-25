package engine

type BaseRouter struct {
}

func (b *BaseRouter) Init(v interface{}) IRouter {
	return b
}

type IRouter interface {
	Init(interface{}) IRouter
}
