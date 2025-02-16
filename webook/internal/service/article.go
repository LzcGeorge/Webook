package service

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository"
	"context"
)

type ArticleService interface {
	Edit(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

// Edit 编辑文章： 返回文章 id
func (a *articleService) Edit(ctx context.Context, article domain.Article) (int64, error) {
	return a.repo.Create(ctx, article)
}
