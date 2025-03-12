package config

import "time"

var DevRedisExpire = RedisExpire{
	ArticleDetail:    time.Minute,
	ArticleFirstPage: time.Minute * 10,
	PublicArticle:    time.Minute * 10,
	Top100:           time.Minute * 3,
}
