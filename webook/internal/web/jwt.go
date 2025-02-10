package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type jwtHandler struct {
}

type UserClaims struct {
	UserId int64
	jwt.RegisteredClaims
	UserAgent string
}

func (u jwtHandler) setJWTToken(ctx *gin.Context, uid int64) {
	// claims 中存储用户的信息
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置 token 的过期时间: 1 分钟（和 lua 代码中的过期时间一致）
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		UserId:    uid,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("sUwYXfLAdddhd1hyWJkWMd4gqQiFznp6"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	//ctx.String(http.StatusOK, "JWT 登录成功，token: %s", tokenStr)
}
