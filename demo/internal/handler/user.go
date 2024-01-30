package handler

import (
	"context"

	"github.com/easonchiu/venom/demo/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserSerivce
}

func NewUserHandler(s *service.UserSerivce) *UserHandler {
	return &UserHandler{s}
}

func (h *UserHandler) Demo(ctx *gin.Context) {
	user, _ := h.service.Get(context.TODO(), "123")
	ctx.JSON(200, user)
}
