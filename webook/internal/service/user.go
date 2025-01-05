package service

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository"
	"context"
	"errors"

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

var ErrInvalidUserOrPassword = errors.New("邮箱或密码不对")

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	// 根据邮箱查询用户是否存在
	user, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		// 用户不存在
		return domain.User{}, ErrInvalidUserOrPassword
	}

	// 查询用户失败：超时或者网络错误
	if err != nil {
		return domain.User{}, err
	}

	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return user, nil
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	user, err := svc.repo.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
