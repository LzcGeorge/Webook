package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type App struct {
	server *gin.Engine
	cron   *cron.Cron
}
