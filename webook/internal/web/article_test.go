package web

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/service"
	svcmocks "Webook/webook/internal/service/mocks"
	myjwt "Webook/webook/internal/web/jwt"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.ArticleService
		reqBody  string
		wantCode int
		wantRes  Result
	}{
		// 发布文章
		{
			name: "发布文章",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "new article and publish",
					Content: "content",
					Author:  domain.Author{Id: 123},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `{
				"title":"new artice and publish",
				"content":"content"
			}`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Data: float64(1),
				Msg:  "发布成功",
			},
		},
		// 发布失败
		{
			name: "发布文章失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "new article and publish",
					Content: "content",
					Author:  domain.Author{Id: 123},
				}).Return(int64(0), errors.New("publish failed"))
				return svc
			},
			reqBody: `{
				"title":"new article and publish",
				"content":"content"
			}`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		// 修改已有文章，并发布
		{
			name: "修改已有文章，并发布",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      1, // 修改已有文章
					Title:   "edit article and publish",
					Content: "content",
					Author:  domain.Author{Id: 123},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `{
				"id":1,
				"title":"edit article and publish",
				"content":"content"
			}`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Data: float64(1),
				Msg:  "发布成功",
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 模拟登录态
			server := gin.Default()
			server.Use(func(c *gin.Context) {
				c.Set("claims", &myjwt.UserClaims{
					UserId: 123,
				})
			})

			articleHandler := NewArticleHandler(tc.mock(ctrl), nil)
			articleHandler.RegisterRoutes(server.Group("/articles"))

			// 创建请求
			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewBuffer([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			assert.Nil(t, err)

			// 执行请求
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			// 检查响应
			assert.Equal(t, tc.wantCode, resp.Code)
			var res Result
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.Nil(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
