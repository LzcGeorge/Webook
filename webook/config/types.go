package config

import "time"

type config struct {
	DB    DBConfig
	Redis RedisConfig
}

type DBConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}

type RedisExpire struct {
	ArticleDetail    time.Duration
	ArticleFirstPage time.Duration
	PublicArticle    time.Duration
	Top100           time.Duration
}
