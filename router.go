package venom

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Method      string
	Path        string
	Middlewares []string // 中间件
	Handle      func(*gin.Context)
	children    []Router
}

func NewRouter(method, path string, handle func(*gin.Context)) Router {
	return Router{Method: method, Path: path, Handle: handle}
}

func NewRouterWithMiddlewares(method, path, middlewares string, handle func(*gin.Context)) Router {
	return Router{Method: method, Path: path, Middlewares: strings.Split(middlewares, ","), Handle: handle}
}

func NewRouterGroup(path string, routers ...Router) Router {
	return Router{Path: path, children: routers}
}

func NewRouterGroupWithMiddlewares(path, middlewares string, routers ...Router) Router {
	return Router{Path: path, Middlewares: strings.Split(middlewares, ","), children: routers}
}

func (r *Router) IsGroup() bool {
	return r.Method == "" && len(r.children) > 0
}

func (r *Router) IsRouter() bool {
	return r.Method != "" && len(r.children) == 0
}
