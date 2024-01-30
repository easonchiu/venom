package router

import (
	"github.com/easonchiu/venom"
	"github.com/easonchiu/venom/demo/internal/handler/ping"
)

type MyRouter struct {
	pingHandler *ping.Handler
}

func NewMyRouter() *MyRouter {
	return &MyRouter{
		pingHandler: ping.NewHandler(),
	}
}

func (r *MyRouter) Load(g *venom.Engine) {
	g.GET("/", r.pingHandler.Ping)
}
