package cache

import (
	"Webook/webook/internal/domain"
	"context"
	"errors"
	"sync"
	"time"
)

type RankingLocalCache struct {
	sync.RWMutex
	top100     []domain.Article
	ddl        time.Time
	expiration time.Duration
}

func NewRankingLocalCache() *RankingLocalCache {
	return &RankingLocalCache{
		expiration: time.Minute * 30, // 设置默认过期时间为30分钟
	}
}

func (l *RankingLocalCache) Set(ctx context.Context, arts []domain.Article) error {
	l.Lock()
	defer l.Unlock()
	// 处理文章内容，保持与 Redis 实现一致
	for i := range arts {
		arts[i].Content = arts[i].Abstract()
	}
	l.top100 = arts
	l.ddl = time.Now().Add(l.expiration)
	return nil
}

func (l *RankingLocalCache) Get(ctx context.Context) ([]domain.Article, error) {
	l.RLock()
	defer l.RUnlock()
	if time.Now().After(l.ddl) {
		return nil, errors.New("cache expired")
	}
	return l.top100, nil
}
