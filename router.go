package venom

import (
  "github.com/gin-gonic/gin"
  "github.com/go-redis/redis/v8"
  "go.mongodb.org/mongo-driver/mongo"
  "net/http"
)

type Router struct {
  Config *Config
  Redis  *redis.Client
  Mongo  *mongo.Database
  Gin    *gin.Engine
}

type Handle func(ctx *Context) bool

func (e *Engine) Router() *Router {
  return &Router{
    Config: e.Config,
    Redis:  e.Redis,
    Mongo:  e.Mongo,
    Gin:    e.Gin,
  }
}

// 调用gin的handle, 并将gin的context包装一层，以便加入更多的功能
func (r *Router) handleGin(httpMethod, path string, handles ...Handle) gin.IRoutes {
  if handles == nil {
    return nil
  }

  funcs := make([]gin.HandlerFunc, 0, len(handles))
  for _, h := range handles {
    funcs = append(funcs, func(gctx *gin.Context) {
      h(&Context{Config: r.Config, Redis: r.Redis, Mongo: r.Mongo, GinContext: gctx})
    })
  }

  return r.Gin.Handle(httpMethod, path, funcs...)
}

func (r *Router) Handle(httpMethod, path string, handles ...Handle) {
  _ = r.handleGin(httpMethod, path, handles...)
}

func (r *Router) GET(path string, handles ...Handle) {
  _ = r.handleGin(http.MethodGet, path, handles...)
}

func (r *Router) POST(path string, handles ...Handle) {
  _ = r.handleGin(http.MethodPost, path, handles...)
}

func (r *Router) PUT(path string, handles ...Handle) {
  _ = r.handleGin(http.MethodPut, path, handles...)
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
