package venom

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, obj interface{}) {
	config := getConfig()

	if config.SuccessFormat != nil {
		ctx.JSON(200, config.SuccessFormat(obj))
		ctx.Abort()
		return
	}
	ctx.JSON(200, obj)
	ctx.Abort()
}

func Fail(ctx *gin.Context, errCode interface{}, errMessage string, obj ...interface{}) {
	config := getConfig()

	var data interface{} = nil
	if obj != nil && len(obj) > 0 {
		data = obj[0]
	}

	if config.FailFormat != nil {
		err := config.FailFormat(errCode, errMessage, data)
		ctx.JSON(200, err)
		_ = ctx.Error(errors.New(fmt.Sprintf("err_code: %v, data: %v", errCode, err)))
		ctx.Abort()
		return
	}

	ctx.JSON(200, data)
	_ = ctx.Error(errors.New(fmt.Sprintf("err_code: %v, data: %v", errCode, data)))
	ctx.Abort()
	return
}
