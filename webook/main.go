package main

import (
	"Webook/webook/internal/web"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	// 注册 User 路由
	u := web.NewUserHandler()
	u.RegisterRoutes(server.Group("/users"))
	_ = server.Run(":8080") // listen and serve on 8080
}
