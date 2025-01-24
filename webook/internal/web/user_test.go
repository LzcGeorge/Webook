package web

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/service"
	svcmocks "Webook/webook/internal/service/mocks"
	"bytes"
	"context"
	"encoding/json"
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

func TestUserHandler_LoginSMSCodeVerify(t *testing.T) {
	testCases := []struct {
		name            string
		mock            func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBody         string
		wantHttpCode    int
		wantResultMsg   string
		wantResultId    int64
		wantResultPhone string
		wantResultCode  int
	}{
		// 验证成功
		{
			name: "验证成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				// 验证码校验成功
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "13812345678", "123456").
					Return(true, nil)

				// 查找或创建用户成功
				userSvc.EXPECT().FindOrCreate(gomock.Any(), "13812345678").
					Return(domain.User{
						Id:    123,
						Phone: "13812345678",
					}, nil)
				return userSvc, codeSvc
			},
			reqBody:         `{"phone":"13812345678","code":"123456"}`,
			wantHttpCode:    http.StatusOK,
			wantResultCode:  0,
			wantResultMsg:   "验证码校验通过",
			wantResultId:    123,
			wantResultPhone: "13812345678",
		},
		// 验证码错误
		{
			name: "验证码错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				// 验证码校验失败
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "13812345678", "123456").
					Return(false, nil)

				return userSvc, codeSvc
			},
			reqBody:        `{"phone":"13812345678","code":"123456"}`,
			wantHttpCode:   http.StatusOK,
			wantResultCode: 4,
			wantResultMsg:  "验证码错误，请重新输入",
		},
		// 验证码错误次数过多
		{
			name: "验证码错误次数过多",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				// 验证码错误次数过多
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "13812345678", "123456").
					Return(false, service.ErrCodeVerifyTooManyTimes)

				return userSvc, codeSvc
			},
			reqBody:        `{"phone":"13812345678","code":"123456"}`,
			wantHttpCode:   http.StatusOK,
			wantResultCode: 5,
			wantResultMsg:  "验证码错误次数过多，请稍后再试",
		},
		// 验证码验证中的系统错误
		{
			name: "验证码验证中的系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				// 验证码验证中的系统错误
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "13812345678", "123456").
					Return(false, errors.New("验证码系统中的其他错误"))

				return userSvc, codeSvc
			},
			reqBody:        `{"phone":"13812345678","code":"123456"}`,
			wantHttpCode:   http.StatusOK,
			wantResultCode: 5,
			wantResultMsg:  "系统错误",
		},
		// 查找或创建用户中的系统错误
		{
			name: "查找或创建用户中的系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				// 验证码校验成功
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "13812345678", "123456").
					Return(true, nil)

				// 查找或创建用户失败
				userSvc.EXPECT().FindOrCreate(gomock.Any(), "13812345678").
					Return(domain.User{
						Id:    123,
						Phone: "13812345678",
					}, errors.New("查找或创建用户中的系统错误"))
				return userSvc, codeSvc
			},
			reqBody:        `{"phone":"13812345678","code":"123456"}`,
			wantHttpCode:   http.StatusOK,
			wantResultCode: 5,
			wantResultMsg:  "系统错误",
		},
		// LoginSMSCodeVerifyReq 格式错误
		{
			name: "LoginSMSCodeVerifyReq 格式错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			// JSON 格式错误: code 缺失
			reqBody:        `{"phone":"13812345678","code"}`,
			wantHttpCode:   http.StatusBadRequest,
			wantResultCode: 5,
			wantResultMsg:  "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 创建 userHandler 及所需的依赖 userService
			server := gin.Default()
			userSvc, codeSvc := tc.mock(ctrl)
			userHandler := NewUserHandler(userSvc, codeSvc)
			userHandler.RegisterRoutes(server.Group("/users"))

			// 创建请求
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms", bytes.NewBuffer([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			assert.Nil(t, err)

			// 执行请求
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			// 检查响应
			assert.Equal(t, tc.wantHttpCode, resp.Code)
			if tc.wantHttpCode != http.StatusOK {
				return
			}

			// 解析响应 JSON
			var result Result
			err = json.Unmarshal(resp.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResultMsg, result.Msg)
			assert.Equal(t, tc.wantResultCode, result.Code)

			// Data 存在
			if result.Data != nil {
				// 先将 result.Data 转为 map
				dataMap := result.Data.(map[string]interface{})
				// 然后从 map 中获取值
				assert.Equal(t, tc.wantResultId, int64(dataMap["Id"].(float64)))
				assert.Equal(t, tc.wantResultPhone, dataMap["Phone"].(string))
			}
		})
	}
}
