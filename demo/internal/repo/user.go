package repo

import (
	"context"

	"github.com/easonchiu/venom/demo/internal/model"
	"github.com/qiniu/qmgo"
)

type UserRepo struct {
	db *qmgo.Collection
}

type IUserRepo interface {
	Get(ctx context.Context, name string) (*model.UserModel, error)
}

var _ IUserRepo = (*UserRepo)(nil)

func NewUserRepo(db *qmgo.Collection) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) Get(ctx context.Context, name string) (*model.UserModel, error) {
	var (
		user *model.UserModel
	)

	err := r.db.Find(ctx, nil).One(user)
	return user, err
}
