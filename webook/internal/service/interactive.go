package service

import (
	"Webook/webook/internal/repository"
	"context"
)

type InteractiveService interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
	IncreaseLike(ctx context.Context, biz string, bizId int64, userId int64) error
	DecreaseLike(ctx context.Context, biz string, bizId int64, userId int64) error
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
