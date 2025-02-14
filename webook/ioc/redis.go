package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// InitRedis 初始化 Redis
func InitRedis() redis.Cmdable {
	// 利用 viper 读取配置文件
	type RedisConfig struct {
		Addr string `yaml:"Addr"`
	}
	var redisConfig RedisConfig
	err := viper.UnmarshalKey("redis", &redisConfig)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisConfig.Addr,
	})
	return redisClient
}
