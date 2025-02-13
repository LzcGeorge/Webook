package limiter

import "context"

type Limiter interface {
	// 返回是否被限流, key 是限流对象
	Limit(ctx context.Context, key string) (bool, error)
}
