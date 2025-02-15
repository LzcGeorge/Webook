//go:build wireinject

package main

import (
	"Webook/webook/internal/repository"
	"Webook/webook/internal/repository/cache"
	"Webook/webook/internal/repository/dao"
	"Webook/webook/internal/service"
	"Webook/webook/internal/web"
	myjwt "Webook/webook/internal/web/jwt"
	"Webook/webook/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		ioc.InitLogger,

		// Dao
		dao.NewUserDAO,

		// Cache
		cache.NewUserCache,
		cache.NewCodeCache,

		// repository
		repository.NewUserRepository,
		repository.NewCodeRepository,

		// Service
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,

		// Handler
		web.NewUserHandler,
		myjwt.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,
		ioc.InitGinMiddleware,
		ioc.InitWebServer,
	)
	return gin.Default()
}
