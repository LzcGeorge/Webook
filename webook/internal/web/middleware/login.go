package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	ignorePaths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path ...string) *LoginMiddlewareBuilder {
	l.ignorePaths = append(l.ignorePaths, path...)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// V1：登录和注册页面不校验登录状态
		// if ctx.Request.URL.Path == "/users/login" ||
		// 	ctx.Request.URL.Path == "/users/signup" {
		// 	return
		// }

		// V2：ignorePaths 中的路径不校验登录状态
		for _, path := range l.ignorePaths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			// 用户未登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
