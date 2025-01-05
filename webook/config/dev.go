package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
