package service

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository"
	repomocks "Webook/webook/internal/repository/mocks"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func Test_UserServiceStruct_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// Login 的参数
		ctx      context.Context
		email    string
		password string

		// Login 的返回值
		wantUser domain.User
		wantErr  error
	}{
		// 登录成功
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 模拟调用 FindByEmail 方法, 参数为邮箱 123@qq.com，返回值为 domain.User 和 nil
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email: "123@qq.com",
						// 密码是 123456#qwer 的加密结果
						Password: "$2a$10$teTdyp4lF/nxYQT506m.cu7z9XylX61m6Sg0zLoWdhcBIa0cGY0em",
						Phone:    "13512345678",
						Ctime:    now,
					}, nil)

				return repo
			},

			// 输入的参数
			email:    "123@qq.com",
			password: "123456#qwer",

			// 期望的返回值
			wantUser: domain.User{
				Email: "123@qq.com",
				// 密码是 123456#qwer 的加密结果
				Password: "$2a$10$teTdyp4lF/nxYQT506m.cu7z9XylX61m6Sg0zLoWdhcBIa0cGY0em",
				Phone:    "13512345678",
				Ctime:    now,
			},

			wantErr: nil,
		},
		// 用户不存在
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 模拟调用 FindByEmail 方法, 参数为邮箱 123@qq.com，返回值为 domain.User 和 nil
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},

			// 输入的参数
			email:    "123@qq.com",
			password: "123456#qwer",

			// 期望的返回值
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		// 查询用户失败
		{
			name: "查询用户失败",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 模拟调用 FindByEmail 方法, 参数为邮箱 123@qq.com，返回值为 domain.User 和 nil
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("db query error"))
				return repo
			},

			// 输入的参数
			email:    "123@qq.com",
			password: "123456#qwer",

			// 期望的返回值
			wantUser: domain.User{},
			wantErr:  errors.New("db query error"),
		},
		// 密码错误
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 模拟调用 FindByEmail 方法, 参数为邮箱 123@qq.com，返回值为 domain.User 和 nil
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email: "123@qq.com",
						// 密码是 123456#qwer 的加密结果
						Password: "$2a$10$teTdyp4lF/nxYQT506m.cu7z9XylX61m6Sg0zLoWdhcBIa0cGY0em",
						Phone:    "13512345678",
						Ctime:    now,
					}, nil)

				return repo
			},

			// 输入的参数
			email:    "123@qq.com",
			password: "123231456#qwer",

			// 期望的返回值
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.mock(ctrl)
			svc := NewUserService(repo, nil)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}

func TestXXX(t *testing.T) {
	password := "123456#qwer"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Println(string(hashedPassword))
}
