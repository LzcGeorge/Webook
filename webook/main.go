package main

import (
	"Webook/webook/internal/repository"
	"Webook/webook/internal/repository/dao"
	"Webook/webook/internal/service"
	"Webook/webook/internal/web"
	"Webook/webook/internal/web/middleware"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

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

	// middleware: 跨域请求
	server.Use(cors.New(cors.Config{
		AllowHeaders: []string{"Content-Type", "Authorization"},
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

	// middleware：利用 session 插件，从 cookie 中获取 sessionID，校验登录状态
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.NewLoginMiddlewareBuilder().
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
