
package config

var K8sConfig = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:30002)/webook?charset=utf8mb4",
	},
	Redis: RedisConfig{
		Addr: "localhost:30003",
	},
}
