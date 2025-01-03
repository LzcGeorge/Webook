//go:build !k8s

// 使用 !k8s 这个编译标签
package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
