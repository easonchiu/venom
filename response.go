/*
 * @Author: zhaozhida zhaozhida@qiniu.com
 * @Date: 2023-07-26 10:32:09
 * @LastEditors: zhaozhida zhaozhida@qiniu.com
 * @LastEditTime: 2023-07-26 14:39:17
 * @Description:
 */
package venom

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, obj interface{}) {
	if config.SuccessFormat != nil {
		ctx.JSON(200, config.SuccessFormat(obj))
		// ctx.Abort()
		return
	}
	ctx.JSON(200, obj)
	// ctx.Abort()
}

func Fail(ctx *gin.Context, errCode interface{}, errMessage string, obj ...interface{}) {
	var data interface{} = nil
	if obj != nil && len(obj) > 0 {
		data = obj[0]
	}

	if config.FailFormat != nil {
		err := config.FailFormat(errCode, errMessage, data)
		ctx.JSON(200, err)
		_ = ctx.Error(fmt.Errorf("err_code: %v, data: %v", errCode, err))
		ctx.Abort()
		return
	}

	ctx.JSON(200, data)
	_ = ctx.Error(fmt.Errorf("err_code: %v, data: %v", errCode, data))
	ctx.Abort()
}
