package cache

import (
	"Webook/webook/config"
	"Webook/webook/internal/domain"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"

	"time"
)

type RankingCache interface {
	GetTop100(ctx context.Context) ([]domain.Article, error)
	SetTop100(ctx context.Context, articles []domain.Article) error
}
type RankingRedisCache struct {
	client     redis.Cmdable
	key        string
	expiration time.Duration
}

func NewRankingCache(client redis.Cmdable) RankingCache {
	return &RankingRedisCache{
		client:     client,
		key:        "ranking:top_100",
		expiration: config.DevRedisExpire.Top100,
	}
}

func (r *RankingRedisCache) GetTop100(ctx context.Context) ([]domain.Article, error) {
	val, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal(val, &res)
	return res, err
}

func (r *RankingRedisCache) SetTop100(ctx context.Context, arts []domain.Article) error {
	for i := range arts {
		arts[i].Content = arts[i].Abstract()
	}
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key, val, r.expiration).Err()
}
