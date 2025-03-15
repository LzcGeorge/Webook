package repository

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/cache"
	"Webook/webook/internal/repository/dao"
	"context"
	"database/sql"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByWechat(ctx context.Context, openId string) (domain.User, error)
	UpdateById(ctx context.Context, user domain.User) error
	GetNameMapByIds(ctx context.Context, ids []int64) (map[int64]string, error)
}
type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

func (repo *CachedUserRepository) Create(ctx context.Context, user domain.User) error {
	// 调用 dao 层进行注册
	return repo.dao.Insert(ctx, repo.domainToEntity(user))
}

func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(user), nil
}

func (repo *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// cache 中的 user 是 domain.User
	user, err := repo.cache.Get(ctx, id)
	if err == nil {
		// 从缓存中获取到用户
		println("从缓存中获取到用户,userID: ", id)
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

func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(user), nil
}

func (repo *CachedUserRepository) entityToDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Phone:    user.Phone.String,
		Password: user.Password,
		Ctime:    time.UnixMilli(user.Ctime),

		// 用户信息
		Nickname: user.Nickname,
		Birthday: time.UnixMilli(user.Birthday),
		AboutMe:  user.AboutMe,
	}
}

func (repo *CachedUserRepository) domainToEntity(user domain.User) dao.User {
	return dao.User{
		Id:       user.Id,
		Email:    sql.NullString{String: user.Email, Valid: user.Email != ""},
		Phone:    sql.NullString{String: user.Phone, Valid: user.Phone != ""},
		Password: user.Password,
		Ctime:    user.Ctime.UnixMilli(),

		// 用户信息
		Nickname: user.Nickname,
		Birthday: user.Birthday.UnixMilli(),
		AboutMe:  user.AboutMe,

		// 微信信息
		WechatOpenId:  sql.NullString{String: user.WechatInfo.OpenId, Valid: user.WechatInfo.OpenId != ""},
		WechatUnionId: sql.NullString{String: user.WechatInfo.UnionId, Valid: user.WechatInfo.UnionId != ""},
	}
}

func (repo *CachedUserRepository) UpdateById(ctx context.Context, user domain.User) error {
	_, err := repo.cache.Get(ctx, user.Id)
	if err == nil {
		// 缓存中存在这个数据, 删除缓存
		_ = repo.cache.Del(ctx, user.Id)
	}
	entity := repo.domainToEntity(user)
	return repo.dao.UpdateById(ctx, entity)
}

func (repo *CachedUserRepository) FindByWechat(ctx context.Context, openId string) (domain.User, error) {
	user, err := repo.dao.FindByWechat(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(user), nil
}

func (repo *CachedUserRepository) GetNameMapByIds(ctx context.Context, ids []int64) (map[int64]string, error) {
	users, err := repo.dao.FindByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]string, len(users))
	for _, user := range users {
		res[user.Id] = user.Nickname
	}
	return res, nil
}
