package venom

import (
  "github.com/gin-gonic/gin"
  "net/http"
)

type IRouter interface {
  handleGin(httpMethod, path string, handles ...Handle) gin.IRoutes
  Handle(httpMethod, path string, handles ...Handle)
  CONNECT(path string, handles ...Handle)
  PUT(path string, handles ...Handle)
  POST(path string, handles ...Handle)
  DELETE(path string, handles ...Handle)
  GET(path string, handles ...Handle)
  HEAD(path string, handles ...Handle)
  OPTIONS(path string, handles ...Handle)
  PATCH(path string, handles ...Handle)
  TRACH(path string, handles ...Handle)
  Group(path string) *Router
}

type Router struct {
  Config         *Config
  Redis          *RedisClient
  Mongo          *MongoClient
  Qmgo           *QmgoClient
  Gin            *gin.Engine
  GinRouterGroup *gin.RouterGroup
}

type Handle func(ctx *Context) bool

func (e *Engine) Router() *Router {
  return &Router{
    Config:         e.Config,
    Redis:          e.Redis,
    Mongo:          e.Mongo,
    Qmgo:           e.Qmgo,
    Gin:            e.Gin,
    GinRouterGroup: e.Gin.Group(""),
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
      h(&Context{r.Config, r.Redis, r.Mongo, r.Qmgo, gctx})
    })
  }

  return r.GinRouterGroup.Handle(httpMethod, path, funcs...)
}

func (r *Router) Handle(httpMethod, path string, handles ...Handle) {
  _ = r.handleGin(httpMethod, path, handles...)
}

func (r *Router) CONNECT(path string, handles ...Handle) {
  r.Handle(http.MethodConnect, path, handles...)
}

func (r *Router) PUT(path string, handles ...Handle) {
  r.Handle(http.MethodPut, path, handles...)
}

func (r *Router) POST(path string, handles ...Handle) {
  r.Handle(http.MethodPost, path, handles...)
}

func (r *Router) DELETE(path string, handles ...Handle) {
  r.Handle(http.MethodDelete, path, handles...)
}

func (r *Router) GET(path string, handles ...Handle) {
  r.Handle(http.MethodGet, path, handles...)
}

func (r *Router) HEAD(path string, handles ...Handle) {
  r.Handle(http.MethodHead, path, handles...)
}

func (r *Router) OPTIONS(path string, handles ...Handle) {
  r.Handle(http.MethodOptions, path, handles...)
}

func (r *Router) PATCH(path string, handles ...Handle) {
  r.Handle(http.MethodPatch, path, handles...)
}

func (r *Router) TRACH(path string, handles ...Handle) {
  r.Handle(http.MethodTrace, path, handles...)
}

func (r *Router) Group(path string) *Router {
  return &Router{
    Config:         r.Config,
    Redis:          r.Redis,
    Mongo:          r.Mongo,
    Qmgo:           r.Qmgo,
    Gin:            r.Gin,
    GinRouterGroup: r.GinRouterGroup.Group(path),
  }
}

var (
  _ IRouter = &Router{}
)
