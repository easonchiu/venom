package router

import (
	"github.com/easonchiu/venom"
	"github.com/easonchiu/venom/demo/internal/handler"
)

type UserRouter struct {
	handler *handler.UserHandler
}

func NewUserRouter(h *handler.UserHandler) *UserRouter {
	return &UserRouter{h}
}

func (r *UserRouter) Load(g *venom.Engine) {
	g.GET("/", r.handler.Demo)
}
