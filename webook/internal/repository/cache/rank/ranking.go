package cache

import (
	"Webook/webook/internal/domain"
	"context"
)

type RankingCache interface {
	GetTop100(ctx context.Context) ([]domain.Article, error)
	SetTop100(ctx context.Context, articles []domain.Article) error
}

type CompositeRankingCache struct {
	local *RankingLocalCache
	redis *RankingRedisCache
}

func NewCompositeRankingCache(local *RankingLocalCache, redis *RankingRedisCache) RankingCache {
	return &CompositeRankingCache{
		local: local,
		redis: redis,
	}
}

func (c *CompositeRankingCache) GetTop100(ctx context.Context) ([]domain.Article, error) {
	// 先尝试从本地缓存获取
	arts, err := c.local.Get(ctx)
	if err == nil {
		return arts, nil
	}

	// 本地缓存失效，从Redis获取
	arts, err = c.redis.Get(ctx)
	if err != nil {
		return nil, err
	}

	// 设置到本地缓存
	_ = c.local.Set(ctx, arts)
	return arts, nil
}

func (c *CompositeRankingCache) SetTop100(ctx context.Context, arts []domain.Article) error {
	// 先设置本地缓存
	if err := c.local.Set(ctx, arts); err != nil {
		return err
	}

	// 再设置Redis缓存
	return c.redis.Set(ctx, arts)
}
