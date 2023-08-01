package venom

import (
	"github.com/gin-gonic/gin"
)

type IMiddleware interface {
	Name() string
	OnStart(*Config)
	OnDestroy(*Config)
	GetGinMiddleware(*Config) gin.HandlerFunc
}

func (engine *Engine) GetMiddlewares(names []string) []IMiddleware {
	result := make([]IMiddleware, 0)

	if len(names) == 0 || engine.__middlewares == nil {
		return result
	}

	for _, name := range names {
		if name == "" {
			continue
		}
		for _, mw := range engine.__middlewares {
			if mw.Name() == name {
				result = append(result, mw)
			}
		}
	}

	return result
}
