package venom

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func setRouter(g *gin.Engine) {
	r := g.Group("/")
	r.GET("", func(ctx *gin.Context) {
		// qmgo := GetQmgoClient()
		// n, err := qmgo.C("celue").Find(ctx.Request.Context(), bson.M{}).Count()

		Fail(ctx, 999, "system error", "data")
		// Success(ctx, "lllll")
	})
}

func TestVenom(t *testing.T) {
	v := Init(&Config{
		Port: "5000",
		Mode: DevelopmentMode,
		SuccessFormat: func(obj interface{}) interface{} {
			return gin.H{
				"code":    0,
				"message": "ok",
				"data":    obj,
			}
		},
		FailFormat: func(errCode interface{}, errMessage string, obj interface{}) interface{} {
			return gin.H{
				"code":    errCode,
				"message": errMessage,
				"data":    obj,
			}
		},
		Logger: LoggerConfig{
			Filename:   "log",
			MaxSize:    500,
			MaxAge:     7,
			MaxBackups: 3,
			Level:      logrus.DebugLevel,
		},
		Qmgo: QmgoConfig{
			URI:      "mongodb://localhost:27017/admin",
			Database: "stock",
		},
	})

	setRouter(v.Engine())
	v.Start()
}
