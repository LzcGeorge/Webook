package ioc

import (
	"Webook/webook/config"

	"github.com/redis/go-redis/v9"
)

// InitRedis 初始化 Redis
func InitRedis() redis.Cmdable {
	var redisConfig = config.Config.Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisConfig.Addr,
	})
	return redisClient
}
