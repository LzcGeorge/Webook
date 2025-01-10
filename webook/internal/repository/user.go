package repository

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/cache"
	"Webook/webook/internal/repository/dao"
	"context"
	"database/sql"
	"time"
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
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

func (repo *UserRepository) Create(ctx context.Context, user domain.User) error {
	// 调用 dao 层进行注册
	return repo.dao.Insert(ctx, repo.domainToEntity(user))
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(user), nil
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
	user = repo.entityToDomain(daoUser)

	if err = repo.cache.Set(ctx, user); err != nil {
		// 记录日志, 做监控，不返回错误

	}
	return user, nil
}

func (repo *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(user), nil
}

func (repo *UserRepository) entityToDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Phone:    user.Phone.String,
		Password: user.Password,
		Ctime:    time.UnixMilli(user.Ctime),
	}
}

func (repo *UserRepository) domainToEntity(user domain.User) dao.User {
	return dao.User{
		Id:       user.Id,
		Email:    sql.NullString{String: user.Email, Valid: user.Email != ""},
		Phone:    sql.NullString{String: user.Phone, Valid: user.Phone != ""},
		Password: user.Password,
		Ctime:    user.Ctime.UnixMilli(),
	}
}
