package venom

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func getEngine() *Engine {
	return Init(&Config{
		Port: "5000",
		Mode: DevelopmentMode,
		SuccessFormat: func(code int, obj interface{}) interface{} {
			return gin.H{
				"code":    0,
				"message": "ok",
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
		Mongo: MongoConfig{
			URI:      "mongodb://localhost:27017/admin",
			Database: "stock",
		},
		Qmgo: QmgoConfig{
			URI:      "mongodb://localhost:27017/admin",
			Database: "stock",
		},
	})
}

func setRouter(v *Engine) {
	r := v.Router()
	r.GET("", func(ctx *Context) bool {
		i, e := ctx.Mongo.C("day").CountDocuments(ctx.Bg(), bson.M{})

		n, err := ctx.Qmgo.C("day").Find(ctx.Bg(), bson.M{}).Count()

		return ctx.Success200(gin.H{
			"count_mongo": i,
			"err_mongo":   e,
			"count_qmgo":  n,
			"err_qmgo":    err,
		})
	})
}

func TestVenom(t *testing.T) {
	v := getEngine()

	setRouter(v)

	_ = v.Start()
}
