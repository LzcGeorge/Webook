package middleware

import (
	"Webook/webook/internal/web"
	"net/http"

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

		// 解析 token 并验证 signature, 同时将 token 中的信息(userID) 解析到 claims 中
		tokenStr := web.ExtractToken(ctx)
		if tokenStr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
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

		// 将 claims 保存到 ctx 中
		ctx.Set("claims", claims)
	}
}
