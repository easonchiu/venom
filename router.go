package venom

import (
  "github.com/gin-gonic/gin"
  "net/http"
)

type Router struct {
  Config *Config
  Client *Client
  Gin    *gin.Engine
}

type Handle func(ctx *Ctx) bool

func (e *Engine) Router() *Router {
  return &Router{
    Config: e.Config,
    Client: e.Client,
    Gin:    e.Gin,
  }
}

// 调用gin的handle, 并将gin的context包装一层，以便加入更多的功能
func (r *Router) handleGin(httpMethod, path string, handle Handle) gin.IRoutes {
  return r.Gin.Handle(httpMethod, path, func(gctx *gin.Context) {
    handle(&Ctx{Config: r.Config, Client: r.Client, GinContext: gctx})
  })
}

func (r *Router) Handle(httpMethod, path string, handle Handle) {
  _ = r.handleGin(httpMethod, path, handle)
}

func (r *Router) GET(path string, handle Handle) {
  _ = r.handleGin(http.MethodGet, path, handle)
}

func (r *Router) POST(path string, handle Handle) {
  _ = r.handleGin(http.MethodPost, path, handle)
}

func (r *Router) PUT(path string, handle Handle) {
  _ = r.handleGin(http.MethodPut, path, handle)
}

type RouterGroup struct {
  Router         *Router
  ginRouterGroup *gin.RouterGroup
}

func (r *Router) Group(path string) *RouterGroup {
  return &RouterGroup{
    Router:         r,
    ginRouterGroup: r.Gin.Group(path),
  }
}

func (r *RouterGroup) GET(path string, handle Handle) {
  r.Router.GET(r.ginRouterGroup.BasePath()+path, handle)
}

func (r *RouterGroup) Group(path string) *RouterGroup {
  return &RouterGroup{
    Router:         r.Router,
    ginRouterGroup: r.ginRouterGroup.Group(path),
  }
}
