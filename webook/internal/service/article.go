package service

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository"
	"context"
)

type ArticleService interface {
	Edit(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
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
	if article.Id > 0 {
		return a.repo.Update(ctx, article)
	}
	return a.repo.Create(ctx, article)
}

// Publish 发布文章
func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {

	return 0, nil
}
