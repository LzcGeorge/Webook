package service

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail

func (svc *UserService) SignUp(ctx context.Context, user domain.User) error {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 调用 repository 层进行注册
	return svc.repo.Create(ctx, user)
}
