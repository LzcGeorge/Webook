package repository

import (
	"Webook/webook/internal/repository/cache"
	"context"
)

var (
	ErrCodeSendTooFrequent    = cache.ErrCodeSetTooFrequent
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CachedCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CachedCodeRepository{
		cache: c,
	}
}

func (repo *CachedCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CachedCodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}
