package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed slide_window.lua
var luaSlideWindow string

// RedisSlideWindowLimiter 基于 Redis 的滑动窗口限流器
type RedisSlideWindowLimiter struct {
	cmd redis.Cmdable
	// 窗口大小
	interval time.Duration
	// 阈值
	rate int
}

func (r *RedisSlideWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, luaSlideWindow, []string{key},
		r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
