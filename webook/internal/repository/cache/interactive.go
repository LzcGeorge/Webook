package cache

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed lua/interactive_incr_cnt.lua
var luaIncrCnt string

type InteractiveCache interface {
	IncreaseReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecreaseLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
}

type RedisInteractiveCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{
		client:     client,
		expiration: time.Minute * 10,
	}
}

func (r *RedisInteractiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}

func (r *RedisInteractiveCache) IncreaseReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)},
		"read_cnt",
		1,
	).Err()
}

func (r *RedisInteractiveCache) IncreaseLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)},
		"like_cnt",
		1,
	).Err()
}

func (r *RedisInteractiveCache) DecreaseLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)},
		"like_cnt",
		1,
	).Err()
}
