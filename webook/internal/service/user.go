package service

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(ctx context.Context, user domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Edit(ctx context.Context, user domain.User) error
}

type UserServiceStruct struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceStruct{
		repo: repo,
	}
}

var ErrUserDuplicate = repository.ErrUserDuplicate

func (svc *UserServiceStruct) SignUp(ctx context.Context, user domain.User) error {
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

func (svc *UserServiceStruct) Login(ctx context.Context, email, password string) (domain.User, error) {
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

func (svc *UserServiceStruct) Profile(ctx context.Context, id int64) (domain.User, error) {
	user, err := svc.repo.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (svc *UserServiceStruct) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	user, err := svc.repo.FindByPhone(ctx, phone)

	// 用户存在，直接返回
	if err != repository.ErrUserNotFound {
		return user, nil
	}

	// 用户不存在，创建用户
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && err != repository.ErrUserDuplicate {
		return domain.User{}, err
	}

	// 根据 phone 查询刚创建的用户
	// 这里会碰到主从延迟的问题，可能查询不到（
	user, err = svc.repo.FindByPhone(ctx, phone)
	return user, err
}

func (svc *UserServiceStruct) Edit(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateById(ctx, user)
}
