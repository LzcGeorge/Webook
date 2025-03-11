package service

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository"
	"context"
)

//go:generate mockgen -source=./interactive.go -package=svcmocks -destination=./mocks/interactive.mock.go
type InteractiveService interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
	IncreaseLike(ctx context.Context, biz string, bizId int64, userId int64) error
	DecreaseLike(ctx context.Context, biz string, bizId int64, userId int64) error
	Collect(ctx context.Context, biz string, bizId int64, collectionId int64, userId int64) error
	Get(ctx context.Context, biz string, bizId int64, userId int64) (domain.Interactive, error)
	GetInterMapByBizIds(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error)
}

type interactiveService struct {
	repo repository.InteractiveRepository
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo: repo,
	}
}

func (s *interactiveService) IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error {
	return s.repo.IncreaseReadCnt(ctx, biz, bizId)
}

func (s *interactiveService) IncreaseLike(ctx context.Context, biz string, bizId int64, userId int64) error {
	return s.repo.IncreaseLikeCnt(ctx, biz, bizId, userId)
}

func (s *interactiveService) DecreaseLike(ctx context.Context, biz string, bizId int64, userId int64) error {
	return s.repo.DecreaseLikeCnt(ctx, biz, bizId, userId)
}

func (s *interactiveService) Collect(ctx context.Context, biz string, bizId int64, collectionId int64, userId int64) error {
	return s.repo.InsertCollection(ctx, biz, bizId, collectionId, userId)
}

func (s *interactiveService) Get(ctx context.Context, biz string, bizId int64, userId int64) (domain.Interactive, error) {
	return s.repo.GetInteractive(ctx, biz, bizId, userId)
}

func (s *interactiveService) GetInterMapByBizIds(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error) {
	return s.repo.GetInterMapByBizIds(ctx, biz, bizIds)
}
