package jwt

import "github.com/gin-gonic/gin"

type Handler interface {
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error
	CheckSession(ctx *gin.Context, ssid string) error
	ClearToken(ctx *gin.Context) error
	ExtractToken(ctx *gin.Context) string
	SetLoginToken(ctx *gin.Context, uid int64) error
}
