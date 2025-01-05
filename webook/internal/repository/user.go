package repository

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/cache"
	"Webook/webook/internal/repository/dao"
	"context"
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

func (repo *UserRepository) Create(ctx context.Context, user domain.User) error {
	// 调用 dao 层进行注册
	return repo.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (repo *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// cache 中的 user 是 domain.User
	user, err := repo.cache.Get(ctx, id)
	if err == nil {
		// 从缓存中获取到用户
		// println("从缓存中获取到用户,userID: ", id)
		return user, nil
	}

	// 缓存中没有这个数据, 从数据库中获取
	daoUser, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	// 将 dao.User 转换为 domain.User
	user = domain.User{
		Id:       daoUser.Id,
		Email:    daoUser.Email,
		Password: daoUser.Password,
	}

	if err = repo.cache.Set(ctx, user); err != nil {
		// 记录日志, 做监控，不返回错误

	}
	return user, nil
}
