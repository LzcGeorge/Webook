//go:build wireinject

package main

import (
	"Webook/webook/internal/repository"
	"Webook/webook/internal/repository/article"
	"Webook/webook/internal/repository/cache"
	rankCache "Webook/webook/internal/repository/cache/rank"
	"Webook/webook/internal/repository/dao"
	article2 "Webook/webook/internal/repository/dao/article"
	"Webook/webook/internal/service"
	"Webook/webook/internal/web"
	myjwt "Webook/webook/internal/web/jwt"
	"Webook/webook/ioc"

	"github.com/google/wire"
)

var rankingSvcSet = wire.NewSet(
	rankCache.NewRankingLocalCache,
	rankCache.NewRankingRedisCache,
	rankCache.NewCompositeRankingCache,
	repository.NewRankingRepository,
	service.NewRankingService,
)

func InitWebServer() *App {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		ioc.InitLogger,

		// Dao
		dao.NewUserDAO,
		article2.NewArticleDAO,
		// article2.NewGormArticleAuthorDAO,
		// article2.NewGormArticleReaderDAO,
		dao.NewInteractiveDAO,

		// Ranking Svc
		rankingSvcSet,
		ioc.InitRankingJob,
		ioc.InitJobs,

		// Cache
		cache.NewUserCache,
		cache.NewCodeCache,
		cache.NewRedisArticleCache,
		cache.NewInteractiveCache,

		// repository
		repository.NewUserRepository,
		repository.NewCodeRepository,
		article.NewArticleRepository,
		// article.NewArticleAuthorRepository,
		// article.NewArticleReaderRepository,
		repository.NewInteractiveRepository,

		// Service
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		// service.NewArticleServiceWithTwoRepo,
		service.NewInteractiveService,

		// Handler
		web.NewUserHandler,
		myjwt.NewRedisJWTHandler,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandler,
		web.NewArticleReaderHandler,
		ioc.InitGinMiddleware,
		ioc.InitWebServer,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
