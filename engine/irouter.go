package engine

type BaseRouter struct {
}

func (b *BaseRouter) Init() IRouter {
	return b
}

type IRouter interface {
	Init() IRouter
}
