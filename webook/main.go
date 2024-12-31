package main

import (
	"Webook/webook/internal/repository"
	"Webook/webook/internal/repository/dao"
	"Webook/webook/internal/service"
	"Webook/webook/internal/web"
	"Webook/webook/internal/web/middleware"
	"Webook/webook/pkg/ginx/middlewares/ratelimit"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/redis/go-redis/v9"

	"gorm.io/driver/mysql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {

	db := initDB()
	u := initUser(db)
	server := initWebServer(u)

	u.RegisterRoutes(server.Group("/users"))
	_ = server.Run(":8080") // listen and serve on 8080
}

func initWebServer(u *web.UserHandler) *gin.Engine {
	server := gin.Default()

	// 限流
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	// middleware: 跨域请求
	server.Use(cors.New(cors.Config{
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 暴露给前端，前端可以从 Header 中获取
		ExposeHeaders: []string{"x-jwt-token"},
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
	}))

	// middleware：获取 sessionID，校验登录状态
	// 使用 memstore 作为 session 的存储
	store := memstore.NewStore([]byte("sUwYXfLAdddhd1hyWJkWMd4gqQiFznp6"), []byte("JKK0iptdv10H1HnVP6mVCk2HDi8WjAKH"))

	// 使用 cookie 作为 session 的存储
	// store := cookie.NewStore([]byte("secret"))

	// 使用 redis 作为 session 的存储
	// store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	// 	// authentication 和 encryption 的密钥
	// 	[]byte("sUwYXfLAdddhd1hyWJkWMd4gqQiFznp6"), []byte("JKK0iptdv10H1HnVP6mVCk2HDi8WjAKH"))
	// if err != nil {
	// 	panic(err)
	// }

	server.Use(sessions.Sessions("mysession", store))
	// server.Use(middleware.NewLoginMiddlewareBuilder().
	// 	IgnorePaths("/users/login", "/users/signup").
	// 	Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/login", "/users/signup").
		Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
