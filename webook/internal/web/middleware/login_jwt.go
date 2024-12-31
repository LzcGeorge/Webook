package middleware

import (
	"Webook/webook/internal/web"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type LoginJWTMiddlewareBuilder struct {
	ignorePaths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path ...string) *LoginJWTMiddlewareBuilder {
	l.ignorePaths = append(l.ignorePaths, path...)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		for _, path := range l.ignorePaths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// token 格式：Bearer <token>
		segs := strings.SplitN(tokenHeader, " ", 2)
		if len(segs) != 2 {
			// token 格式不对
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 解析 token 并验证 signature, 同时将 token 中的信息(userID) 解析到 claims 中
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("sUwYXfLAdddhd1hyWJkWMd4gqQiFznp6"), nil
		})
		if err != nil {
			// 解析失败
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.UserId == 0 {
			// 解析成功，但是 token 无效
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 检查 userAgent 是否一致
		if claims.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// token refresh: 每十秒刷新一次 token
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("sUwYXfLAdddhd1hyWJkWMd4gqQiFznp6"))
			if err != nil {
				// 记录日志
				println(err.Error())
			}

			ctx.Header("x-jwt-token", tokenStr)
		}

		// 将 claims 保存到 ctx 中
		ctx.Set("claims", claims)
	}
}
