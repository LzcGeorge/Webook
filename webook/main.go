package main

import (
	"Webook/webook/internal/web"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
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

	// 注册 User 路由
	u := web.NewUserHandler()
	u.RegisterRoutes(server.Group("/users"))
	_ = server.Run(":8080") // listen and serve on 8080
}
