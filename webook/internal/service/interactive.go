package service

import (
	"Webook/webook/internal/repository"
	"context"
)

type InteractiveService interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
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
