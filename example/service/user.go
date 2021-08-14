package service

import (
  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/bson"
  "venom"
)

type UserService struct{}

func (*UserService) Get(ctx *venom.Context) (gin.H, error) {
  res := make(gin.H)
  err := ctx.Mongo.Collection("users").FindOne(nil, bson.M{}).Decode(&res)
  return res, err
}
