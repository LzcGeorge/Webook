//go:build wireinject

package main

import (
	"Webook/webook/internal/repository"
	"Webook/webook/internal/repository/article"
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
		dao.NewArticleDAO,

		// Cache
		cache.NewUserCache,
		cache.NewCodeCache,

		// repository
		repository.NewUserRepository,
		repository.NewCodeRepository,
		article.NewArticleRepository,

		// Service
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		// Handler
		web.NewUserHandler,
		myjwt.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		ioc.InitGinMiddleware,
		ioc.InitWebServer,
	)
	return gin.Default()
}
