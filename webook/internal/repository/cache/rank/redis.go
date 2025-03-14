package cache

import (
	"Webook/webook/config"
	"Webook/webook/internal/domain"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RankingRedisCache struct {
	client     redis.Cmdable
	key        string
	expiration time.Duration
}

func NewRankingRedisCache(client redis.Cmdable) *RankingRedisCache {
	return &RankingRedisCache{
		client:     client,
		key:        "ranking:top_100",
		expiration: config.DevRedisExpire.Top100,
	}
}

func (r *RankingRedisCache) Get(ctx context.Context) ([]domain.Article, error) {
	val, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal(val, &res)
	return res, err
}

func (r *RankingRedisCache) Set(ctx context.Context, arts []domain.Article) error {
	for i := range arts {
		arts[i].Content = arts[i].Abstract()
	}
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key, val, r.expiration).Err()
}
