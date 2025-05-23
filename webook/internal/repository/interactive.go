package repository

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/cache"
	"Webook/webook/internal/repository/dao"
	"context"
)

type InteractiveRepository interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
	IncreaseLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error
	DecreaseLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error
	InsertCollection(ctx context.Context, biz string, bizId int64, collectionId int64, userId int64) error
	GetInteractive(ctx context.Context, biz string, bizId int64, userId int64) (domain.Interactive, error)
	GetInterMapByBizIds(ctx context.Context, biz string, BizIds []int64, userId int64) (map[int64]domain.Interactive, error)
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

func (r *interactiveRepository) InsertCollection(ctx context.Context, biz string, bizId int64, collectionId int64, userId int64) error {
	err := r.dao.InsertCollection(ctx, biz, bizId, collectionId, userId)
	if err != nil {
		return err
	}

	return r.cache.IncreaseCollectCntIfPresent(ctx, biz, bizId)
}

func (r *interactiveRepository) GetInteractive(ctx context.Context, biz string, bizId int64, userId int64) (domain.Interactive, error) {
	interactive, err := r.dao.GetInteractive(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	liked, err := r.dao.GetLiked(ctx, biz, bizId, userId)
	if err != nil {
		return domain.Interactive{}, err
	}
	collected, err := r.dao.GetCollected(ctx, biz, bizId, userId)

	return domain.Interactive{
		ReadCnt:    interactive.ReadCnt,
		LikeCnt:    interactive.LikeCnt,
		CollectCnt: interactive.CollectCnt,
		Liked:      liked,
		Collected:  collected,
	}, nil
}

func (r *interactiveRepository) GetInterMapByBizIds(ctx context.Context, biz string, BizIds []int64, userId int64) (map[int64]domain.Interactive, error) {
	inters, err := r.dao.GetByBizIds(ctx, biz, BizIds)
	if err != nil {
		return nil, err
	}
	likes, err := r.dao.GetLikedByBizIds(ctx, biz, BizIds, userId)
	if err != nil {
		return nil, err
	}
	collects, err := r.dao.GetCollectedByBizIds(ctx, biz, BizIds, userId)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.Interactive, len(inters))
	for _, inter := range inters {
		res[inter.BizId] = domain.Interactive{
			ReadCnt:    inter.ReadCnt,
			LikeCnt:    inter.LikeCnt,
			CollectCnt: inter.CollectCnt,
		}
	}
	const likeValid = 1
	const unCollected = 0
	for _, like := range likes {
		if like.Status == likeValid {
			if inter, ok := res[like.BizId]; ok {
				inter.Liked = true
				res[like.BizId] = inter
			}
		}
	}
	for _, collect := range collects {
		if collect.Cid != unCollected {
			if inter, ok := res[collect.BizId]; ok {
				inter.Collected = true
				res[collect.BizId] = inter
			}
		}
	}
	return res, nil
}
