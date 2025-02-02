package repository

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/cache"
	cachemocks "Webook/webook/internal/repository/cache/mocks"
	"Webook/webook/internal/repository/dao"
	daomocks "Webook/webook/internal/repository/dao/mocks"
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	now := time.Now()
	// 你要去掉毫秒以外的部分
	// 111ms
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx context.Context
		id  int64

		wantUser domain.User
		wantErr  error
	}{
		// 缓存未命中，查询成功
		{
			name: "缓存未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 从缓存中未找到数据
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotFound)

				// 从数据库中查询, 返回的事 dao.User
				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{
						Id:       123,
						Email:    "123@qq.com",
						Password: "this is password",
						Phone: sql.NullString{
							String: "13512345678",
							Valid:  true,
						},
						Ctime: now.UnixMilli(),
						Utime: now.UnixMilli(),
					}, nil)

				// 写入缓存
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "this is password",
					Phone:    "13512345678",
					Ctime:    now,
				}).Return(nil)

				return d, c
			},

			// 输入的参数
			ctx: context.Background(),
			id:  123,

			// 期望的输出
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "this is password",
				Phone:    "13512345678",
				Ctime:    now,
			},
			wantErr: nil,
		},

		// 缓存命中，查询成功
		{
			name: "缓存命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 从缓存中找到数据
				c := cachemocks.NewMockUserCache(ctrl)
				d := daomocks.NewMockUserDAO(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{
						Id:       123,
						Email:    "123@qq.com",
						Password: "this is password",
						Phone:    "13512345678",
						Ctime:    now,
					}, nil)

				return d, c
			},

			// 输入的参数
			ctx: context.Background(),
			id:  123,

			// 期望的输出
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "this is password",
				Phone:    "13512345678",
				Ctime:    now,
			},
			wantErr: nil,
		},

		// 缓存未命中，查询数据库失败
		{
			name: "缓存未命中，查询数据库失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 从缓存中未找到数据
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotFound)

				// 从数据库中查询, 返回的事 dao.User
				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{}, errors.New("db query error"))

				return d, c
			},

			// 输入的参数
			ctx: context.Background(),
			id:  123,

			// 期望的输出
			wantUser: domain.User{},
			wantErr:  errors.New("db query error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dao, cache := tc.mock(ctrl)
			repo := NewUserRepository(dao, cache)
			user, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantUser, user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
