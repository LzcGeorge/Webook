package web

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtHandler struct {
	accessTokenKey  []byte
	refreshTokenKey []byte
}

func NewJWTHandler() jwtHandler {
	return jwtHandler{
		accessTokenKey:  []byte("sUwYXfLAdddhd1hyWJkWMd4gqQiFznp6"),
		refreshTokenKey: []byte("sUwYXfLAdddhd1hyWJkWMd4gqQiFznv2"),
	}
}

type UserClaims struct {
	UserId int64
	Ssid   string
	jwt.RegisteredClaims
	UserAgent string
}

func (h jwtHandler) setJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	// claims 中存储用户的信息
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置 token 的过期时间: 1 分钟（和 lua 代码中的过期时间一致）
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		UserId:    uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.accessTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

type RefreshTokenClaims struct {
	Uid  int64
	Ssid string
	jwt.RegisteredClaims
}

// setRefreshToken: 生成 refresh token
func (h jwtHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefreshTokenClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	refreshTokenStr, err := refreshToken.SignedString(h.refreshTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", refreshTokenStr)
	return nil
}

// ExtractToken: 从 Authorization 中提取 token
func ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	if tokenHeader == "" {
		// 没有登录
		return ""
	}

	// token 格式：Bearer <token>
	segs := strings.SplitN(tokenHeader, " ", 2)
	if len(segs) != 2 {
		// token 格式不对
		return ""
	}

	return segs[1]
}

func (h jwtHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.setJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}

	err = h.setRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}

	return nil
}
