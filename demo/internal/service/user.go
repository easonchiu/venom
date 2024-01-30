package service

import (
	"context"

	"github.com/easonchiu/venom/demo/internal/model"
	"github.com/easonchiu/venom/demo/internal/repo"
)

type UserSerivce struct {
	repo *repo.UserRepo
}

type IUserService interface {
	Get(ctx context.Context, name string) (*model.UserModel, error)
}

var _ IUserService = (*UserSerivce)(nil)

func NewUserService(r *repo.UserRepo) *UserSerivce {
	return &UserSerivce{r}
}

func (u *UserSerivce) Get(ctx context.Context, name string) (*model.UserModel, error) {
	return u.repo.Get(ctx, name)
}
