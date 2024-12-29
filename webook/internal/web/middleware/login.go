package middleware

import (
	"net/http"
	"time"

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
			return
		}

		// updateTime 是 session 中记录的更新时间
		updateTime := sess.Get("update_time")
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now().UnixMilli()
		if updateTime == nil {
			// 如果 updateTime 为空，则设置为当前时间
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
			return
		}

		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 如果 updateTime 已经超过 10 分钟，则重新设置 updateTime
		if now-updateTimeVal > 1000*30 { // 30 sec
			sess.Set("update_time", now)

			sess.Save()
		}
	}
}
