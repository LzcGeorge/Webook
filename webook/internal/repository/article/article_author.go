package article

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/dao/article"
	"context"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) (int64, error)
}

type articleAuthorRepository struct {
	dao article.ArticleAuthorDAO
}

func NewArticleAuthorRepository(dao article.ArticleAuthorDAO) ArticleAuthorRepository {
	return &articleAuthorRepository{
		dao: dao,
	}
}

func (r *articleAuthorRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return r.dao.Insert(ctx, ToArticleEntity(art))
}

func (r *articleAuthorRepository) Update(ctx context.Context, art domain.Article) (int64, error) {
	return r.dao.UpdateById(ctx, ToArticleEntity(art))
}
