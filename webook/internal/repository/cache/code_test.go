package cache

import (
	"Webook/webook/internal/repository/cache/redismocks"
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
)

func TestCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable
		// 输入的参数
		ctx   context.Context
		biz   string
		phone string
		code  string

		// 预期的输出
		wantErr error
	}{
		// 设置成功
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				client := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				// 设置返回值
				res.SetVal(int64(0))
				// 设置期望值
				client.EXPECT().Eval(gomock.Any(), luaSetCode,
					// fmt.Sprintf("phone_code:%s:%s", biz, phone)
					[]string{"phone_code:login:1351234565768"},
					// code
					[]any{"190010"},
				).Return(res)
				return client
			},

			// 输入的参数
			ctx:   context.Background(),
			biz:   "login",
			phone: "1351234565768",
			code:  "190010",

			// 预期的输出
			wantErr: nil,
		},
		// redis 错误
		{
			name: "redis 错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				client := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				// 设置返回值
				res.SetErr(errors.New("redis error"))
				// 设置期望值
				client.EXPECT().Eval(gomock.Any(), luaSetCode,
					// fmt.Sprintf("phone_code:%s:%s", biz, phone)
					[]string{"phone_code:login:1351234565768"},
					// code
					[]any{"190010"},
				).Return(res)
				return client
			},

			// 输入的参数
			ctx:   context.Background(),
			biz:   "login",
			phone: "1351234565768",
			code:  "190010",

			// 预期的输出
			wantErr: errors.New("redis error"),
		},

		// 发送验证码太频繁
		{
			name: "发送验证码太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				client := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				// 设置返回值
				res.SetVal(int64(-1))
				// 设置期望值
				client.EXPECT().Eval(gomock.Any(), luaSetCode,
					// fmt.Sprintf("phone_code:%s:%s", biz, phone)
					[]string{"phone_code:login:1351234565768"},
					// code
					[]any{"190010"},
				).Return(res)
				return client
			},

			// 输入的参数
			ctx:   context.Background(),
			biz:   "login",
			phone: "1351234565768",
			code:  "190010",

			// 预期的输出
			wantErr: ErrCodeSetTooFrequent,
		},
		// 系统错误
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				client := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				// 设置返回值
				res.SetVal(int64(5))
				// 设置期望值
				client.EXPECT().Eval(gomock.Any(), luaSetCode,
					// fmt.Sprintf("phone_code:%s:%s", biz, phone)
					[]string{"phone_code:login:1351234565768"},
					// code
					[]any{"190010"},
				).Return(res)
				return client
			},

			// 输入的参数
			ctx:   context.Background(),
			biz:   "login",
			phone: "1351234565768",
			code:  "190010",

			// 预期的输出
			wantErr: errors.New("系统错误: 发送验证码"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			testRedis := tc.mock(ctrl)
			cache := NewCodeCache(testRedis)
			err := cache.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
