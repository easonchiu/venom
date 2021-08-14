package controller

import (
  "venom"
  "venom/example/service"
)

type UserController struct {
  UserService *service.UserService
}

func (c *UserController) Get(ctx *venom.Context) bool {
  res, err := c.UserService.Get(ctx)

  if err != nil {
    return ctx.Error200("999999", err)
  }

  return ctx.Success200(res)
}
