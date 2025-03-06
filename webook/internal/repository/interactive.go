package repository

import (
	"Webook/webook/internal/repository/cache"
	"Webook/webook/internal/repository/dao"
	"context"
)

type InteractiveRepository interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error
	DecreaseLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error
}

type interactiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

func NewInteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache) InteractiveRepository {
	return &interactiveRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *interactiveRepository) IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error {
	if err := r.dao.IncreaseReadCnt(ctx, biz, bizId); err != nil {
		return err
	}

	// redis 中实现自增
	// 如果 dao 自增成功，数据库中的数据更新
	// 但是 redis 中更新失败（缓存过期 balabala）
	// 导致数据库和 redis 中的数据不一致
	//
	// 由于用户对阅读量不敏感，所以可以容忍这种不一致
	// 所以使用 redis 自增，后续有 Set 方法来回写 redis
	return r.cache.IncreaseReadCntIfPresent(ctx, biz, bizId)
}

func (r *interactiveRepository) IncreaseLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error {
	err := r.dao.InsertLikeInfo(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}

	// 缓存中增加点赞信息
	return r.cache.IncreaseLikeCntIfPresent(ctx, biz, bizId)
}

func (r *interactiveRepository) DecreaseLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error {
	err := r.dao.DeleteLikeInfo(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}

	// 缓存中减少点赞信息
	return r.cache.DecreaseLikeCntIfPresent(ctx, biz, bizId)
}
