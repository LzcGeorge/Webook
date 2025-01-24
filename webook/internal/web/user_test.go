package web

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/service"
	svcmocks "Webook/webook/internal/service/mocks"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func TestUserHandler_SignUp(t *testing.T) {
	// 测试用例定义
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		// 注册成功
		{
			name: "注册成功",
			// 模拟依赖：返回一个 mock 的 UserService
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				// 期待调用 SignUp 方法，传入任意 context 和匹配的 domain.User 对象，返回 nil
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "1234#qwe",
				}).Return(nil)
				return userSvc
			},
			// 请求参数
			reqBody: `{"email":"123@qq.com","password":"1234#qwe","confirmPassword":"1234#qwe"}`,
			// 期望响应
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		// Signup 调用之前：Bind 失败，signupReq 格式错误
		{
			name: "SignUpReq 格式错误，Bind 失败",
			// 模拟依赖：返回一个 mock 的 UserService
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			// 请求参数, 设置为无效的 json 格式
			reqBody: `{"email":"123@qq.com","password":}`,
			// 期望响应
			wantCode: http.StatusBadRequest,
		},
		// Signup 调用之前：邮箱格式错误
		{
			name: "邮箱格式错误",
			// 模拟依赖：返回一个 mock 的 UserService
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			// 请求参数
			reqBody: `{"email":"123","password":"1234#qwe","confirmPassword":"1234#qwe"}`,
			// 期望响应
			wantCode: http.StatusOK,
			wantBody: "你的邮箱格式不对",
		},
		// Signup 调用之前：密码格式错误
		{
			name: "密码格式错误",
			// 模拟依赖：返回一个 mock 的 UserService
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			// 请求参数
			reqBody: `{"email":"123@qq.com","password":"123","confirmPassword":"1234#qwe"}`,
			// 期望响应
			wantCode: http.StatusOK,
			wantBody: "密码必须包含至少一个字母、数字、特殊字符，长度至少8位",
		},
		// Signup 调用之前：确认密码和密码不一致
		{
			name: "确认密码和密码不一致",
			// 模拟依赖：返回一个 mock 的 UserService
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			// 请求参数
			reqBody: `{"email":"123@qq.com","password":"1234#qwe","confirmPassword":"111234#qwe"}`,
			// 期望响应
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		// Signup 调用之后：邮箱冲突
		{
			name: "邮箱冲突",
			// 模拟依赖：返回一个 mock 的 UserService
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "1234#qwe",
				}).Return(service.ErrUserDuplicate)
				return userSvc
			},
			// 请求参数
			reqBody: `{"email":"123@qq.com","password":"1234#qwe","confirmPassword":"1234#qwe"}`,
			// 期望响应
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		// Signup 调用之后：系统异常
		{
			name: "系统异常",
			// 模拟依赖：返回一个 mock 的 UserService
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "1234#qwe",
				}).Return(errors.New("mock error"))
				return userSvc
			},
			// 请求参数
			reqBody: `{"email":"123@qq.com","password":"1234#qwe","confirmPassword":"1234#qwe"}`,
			// 期望响应
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}

	//
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建 userHandler 及所需的依赖 userService
			server := gin.Default()
			userSvc := tc.mock(ctrl)
			userHandler := NewUserHandler(userSvc, nil)
			userHandler.RegisterRoutes(server.Group("/users"))

			// 创建请求
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			assert.Nil(t, err)

			// 执行请求
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			// 检查响应
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

func TestMock(t *testing.T) {
	// 创建一个控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userSvc := svcmocks.NewMockUserService(ctrl)

	userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
		Return(errors.New("mock error"))

	err := userSvc.SignUp(context.Background(), domain.User{
		Email: "123@qq.com",
	})

	t.Log(err)

}
