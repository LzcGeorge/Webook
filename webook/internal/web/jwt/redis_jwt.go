package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisJWTHandler struct {
	cmd                redis.Cmdable
	signingMethod      jwt.SigningMethod
	refreshTokenExpire time.Duration
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd:                cmd,
		signingMethod:      jwt.SigningMethodHS512,
		refreshTokenExpire: time.Hour * 24 * 7,
	}
}

type UserClaims struct {
	UserId int64
	Ssid   string
	jwt.RegisteredClaims
	UserAgent string
}

var AccessTokenKey = []byte("sUwYXfLAdddhd1hyWJkWMd4gqQiFznp6")
var RefreshTokenKey = []byte("sUwYXfLAdddhd1hyWJkWMd4gqQiFznv2")

// SetJWTToken 生成 JWT token
func (h *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	// claims 中存储用户的信息
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置 token 的过期时间: 1 分钟（和 lua 代码中的过期时间一致）
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserId:    uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(h.signingMethod, claims)
	tokenStr, err := token.SignedString(AccessTokenKey)
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

// SetRefreshToken 重新生成 refresh token
func (h *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefreshTokenClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.refreshTokenExpire)),
		},
	}

	refreshToken := jwt.NewWithClaims(h.signingMethod, claims)
	refreshTokenStr, err := refreshToken.SignedString(RefreshTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", refreshTokenStr)
	return nil
}

// CheckSession 检查 Redis 中是否存在 ssid，存在说明已经退出登录
func (h *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	cnt, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	if err != nil || cnt > 0 {
		return fmt.Errorf("您已退出登录")
	}
	return nil
}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	// 设置 JWT token 和 refresh token 为空
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	// 获取 JWT 中的用户信息
	claims, ok := ctx.Get("claims")
	if !ok {
		return fmt.Errorf("系统错误")
	}
	userClaims := claims.(*UserClaims)

	// 设置 Ssid 为有效，表示退出登录
	return h.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", userClaims.Ssid), "", time.Hour*24*7).Err()
}

// ExtractToken: 从 Authorization 中提取 token
func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
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

// SetLoginToken: 生成 JWT token 和 refresh token
func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}

	err = h.SetRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}

	return nil
}
