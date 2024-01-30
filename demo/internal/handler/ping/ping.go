package ping

import "github.com/gin-gonic/gin"

type Handler struct {
	mySvr string
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Ping(ctx *gin.Context) {
	ctx.JSON(200, "??")
}
