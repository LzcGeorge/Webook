package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	server := InitWebServer()

	// 测试
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, Webook!")
	})

	_ = server.Run(":8080") // listen and serve on 8080
}
