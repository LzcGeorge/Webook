package repository

import (
	"Webook/webook/internal/domain"
	cache "Webook/webook/internal/repository/cache/rank"
	"context"
)

type RankingRepository interface {
	GetTop100(ctx context.Context) ([]domain.Article, error)
	ReplaceTop100(ctx context.Context, articles []domain.Article) error
}

type CachedRankingRepository struct {
	cache cache.RankingCache
}

func NewRankingRepository(cache cache.RankingCache) RankingRepository {
	return &CachedRankingRepository{
		cache: cache,
	}
}
func (repo *CachedRankingRepository) GetTop100(ctx context.Context) ([]domain.Article, error) {
	return repo.cache.GetTop100(ctx)
}

func (repo *CachedRankingRepository) ReplaceTop100(ctx context.Context, articles []domain.Article) error {
	return repo.cache.SetTop100(ctx, articles)
}
