package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type LocalCodeCache struct {
	cache *lru.Cache
	mu    sync.Mutex
	exp   time.Duration
}

func NewLocalCodeCache(cache *lru.Cache, exp time.Duration) *LocalCodeCache {
	return &LocalCodeCache{
		cache: cache,
		exp:   exp,
	}
}

func (c *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.key(biz, phone)
	now := time.Now()

	val, ok := c.cache.Get(key)
	const maxCnt = 3
	if !ok {
		// 不存在，创建一个
		c.cache.Add(key, codeItem{
			code:   code,
			cnt:    maxCnt,
			expire: now.Add(c.exp),
		})
		return nil
	}

	itm, _ := val.(codeItem)
	if itm.expire.Sub(now) > time.Minute*9 {
		return ErrCodeSetTooFrequent
	}
	c.cache.Add(key, codeItem{
		code:   code,
		cnt:    maxCnt,
		expire: now.Add(c.exp),
	})
	return nil
}

func (c *LocalCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.key(biz, phone)
	val, ok := c.cache.Get(key)
	if !ok {
		return false, nil
	}
	itm, _ := val.(codeItem)
	if itm.cnt <= 0 {
		return false, ErrCodeVerifyTooManyTimes
	}
	itm.cnt--
	c.cache.Add(key, itm)
	return itm.code == inputCode, nil
}

func (l *LocalCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

type codeItem struct {
	code string
	// 可验证次数
	cnt int
	// 过期时间
	expire time.Time
}
