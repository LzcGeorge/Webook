package main

import (
	"net/http"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	InitViperWithFlags()
	InitLogger()
	server := InitWebServer()

	// 测试
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, Webook!")
	})

	_ = server.Run(":8080") // listen and serve on 8080
}

func InitViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func InitViperWithFlags() {
	// 设置默认配置文件: config/config.yaml
	// 若运行时传入 --config 参数，则使用传入的配置文件
	cfile := pflag.String("config", "config/config.yaml", "viper config file")
	pflag.Parse()

	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)

	// 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		println("Viper Config file changed:", e.Name)
	})

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func InitViperRemote() {
	// 设置远程配置文件: etcd
	viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "/webook")

	// 设置配置文件类型
	viper.SetConfigType("yaml")

	// 读取配置文件
	err := viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
