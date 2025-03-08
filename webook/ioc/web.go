package ioc

import (
	"Webook/webook/internal/web"
	myjwt "Webook/webook/internal/web/jwt"
	"Webook/webook/internal/web/middleware"
	logger2 "Webook/webook/pkg/ginx/middlewares/logger"
	"Webook/webook/pkg/ginx/middlewares/ratelimit"
	"Webook/webook/pkg/logger"
	"context"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"Webook/webook/pkg/limiter"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func newCORSConfig() cors.Config {
	return cors.Config{
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 暴露给前端，前端可以从 Header 中获取
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		// 允许跨域请求携带 cookie
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 本地开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}
}

// InitGinMiddleware 初始化 Gin 中间件
func InitGinMiddleware(redisClient redis.Cmdable, jwthandler myjwt.Handler, l logger.Logger) []gin.HandlerFunc {
	bd := logger2.NewBuilder(func(ctx context.Context, al *logger2.AccessLog) {
		l.Debug("HTTP请求", logger.Field{Key: "al", Value: al})
	}).AllowReqBody(true).AllowRespBody()
	viper.OnConfigChange(func(in fsnotify.Event) {
		ok := viper.GetBool("web.logreq")
		bd.AllowReqBody(ok)
	})
	return []gin.HandlerFunc{
		cors.New(newCORSConfig()),
		// 限流
		ratelimit.NewBuilder(limiter.NewRedisSlideWindowLimiter(redisClient, time.Second, 1000)).Build(),

		// 检查是否满足登录条件
		middleware.NewLoginJWTMiddlewareBuilder(jwthandler).
			IgnorePaths("/users/login", "/users/signup").
			IgnorePaths("/users/login_sms/code/send", "/users/login_sms").
			IgnorePaths("/oauth2/wechat/authurl", "/oauth2/wechat/callback").
			IgnorePaths("/users/refresh_token").
			Build(),
	}
}

// InitWebServer 初始化 Web 服务器
func InitWebServer(middlewares []gin.HandlerFunc,
	userHdl *web.UserHandler, wechatHdl *web.OAuth2WechatHandler,
	articleHdl *web.ArticleHandler, articleReaderHdl *web.ArticleReaderHandler,
) *gin.Engine {
	server := gin.Default()

	// 使用中间件
	server.Use(middlewares...)

	// 用户模块
	userHdl.RegisterRoutes(server.Group("/users"))
	wechatHdl.RegisterRoutes(server.Group("/oauth2/wechat"))

	// 文章模块
	articleHdl.RegisterRoutes(server.Group("/articles"))
	// 线上库文章
	articleReaderHdl.RegisterRoutes(server.Group("articles/pub"))
	return server
}
