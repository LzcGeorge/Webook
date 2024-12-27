package repository

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/dao"
	"context"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail

func (repo *UserRepository) Create(ctx context.Context, user domain.User) error {
	// 调用 dao 层进行注册
	return repo.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}
