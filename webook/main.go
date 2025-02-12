package main

import (
	"Webook/webook/config"
	"Webook/webook/internal/repository"
	"Webook/webook/internal/repository/cache"
	"Webook/webook/internal/repository/dao"
	"Webook/webook/internal/service"
	"Webook/webook/internal/service/oauth2/wechat"
	"Webook/webook/internal/service/sms/memory"
	"Webook/webook/internal/web"
	"Webook/webook/internal/web/middleware"
	"os"

	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/redis/go-redis/v9"

	myjwt "Webook/webook/internal/web/jwt"

	"gorm.io/driver/mysql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {

	// 初始化数据库
	db := initDB()

	// 初始化 Redis
	var redisConfig = config.Config.Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisConfig.Addr,
	})

	// 初始化 jwt
	jwtHandler := myjwt.NewRedisJWTHandler(redisClient)

	// 初始化 user service
	user, userSvc := initUser(db, redisClient, jwtHandler)
	server := initWebServer(jwtHandler)
	user.RegisterRoutes(server.Group("/users"))

	// 配置微信扫码登录
	wechatSvc := InitWechatService()
	wechatHandler := web.NewOAuth2WechatHandler(wechatSvc, userSvc, jwtHandler)
	wechatHandler.RegisterRoutes(server.Group("/oauth2/wechat"))

	// 测试
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, Webook!")
	})

	_ = server.Run(":8080") // listen and serve on 8080
}

func initWebServer(jwtHandler myjwt.Handler) *gin.Engine {
	server := gin.Default()

	// TODO: 限流，抽象接口了，需要修改
	// server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	// middleware: 跨域请求
	server.Use(cors.New(newCORSConfig()))

	// middleware：获取 sessionID，校验登录状态
	// 使用 memstore 作为 session 的存储
	store := memstore.NewStore([]byte("sUwYXfLAdddhd1hyWJkWMd4gqQiFznp6"), []byte("JKK0iptdv10H1HnVP6mVCk2HDi8WjAKH"))
	server.Use(sessions.Sessions("mysession", store))

	// 忽略登录状态的请求

	server.Use(middleware.NewLoginJWTMiddlewareBuilder(jwtHandler).
		IgnorePaths("/users/login", "/users/signup").
		IgnorePaths("/users/login_sms/code/send", "/users/login_sms").
		IgnorePaths("/oauth2/wechat/authurl", "/oauth2/wechat/callback").
		IgnorePaths("/users/refresh_token").
		Build())

	return server
}

func initUser(db *gorm.DB, redisClient redis.Cmdable, jwtHandler myjwt.Handler) (*web.UserHandler, service.UserService) {
	// 用户基本操作：注册、登录、获取用户信息
	userDao := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(redisClient)
	userRepo := repository.NewUserRepository(userDao, userCache)
	userSvc := service.NewUserService(userRepo)

	// 验证码
	codeCache := cache.NewCodeCache(redisClient)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsSvc := memory.NewService() // 采用 memory 作为 sms 的实现
	codeSvc := service.NewCodeService(codeRepo, smsSvc)

	user := web.NewUserHandler(userSvc, codeSvc, jwtHandler)
	return user, userSvc
}

func initDB() *gorm.DB {
	var dbConfig = config.Config.DB // 也可以使用 k8s 的配置
	db, err := gorm.Open(mysql.Open(dbConfig.DSN))

	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func InitWechatService() wechat.Service {
	appID, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		// panic("找不到环境变量 WECHAT_APP_ID")
		appID = "wx6666666666666666"
	}
	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		// panic("找不到环境变量 WECHAT_APP_SECRET")
		appSecret = "66666666666666666666666666666666"
	}
	return wechat.NewService(appID, appSecret)
}

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
