package article

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/dao/article"
	"context"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) (int64, error)
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, art domain.Article) (int64, error)
}

type CachedArticleRepository struct {
	dao article.ArticleDAO
}

func NewArticleRepository(dao article.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, ToArticleEntity(art))
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.UpdateById(ctx, ToArticleEntity(art))
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Upsert(ctx, ToArticleEntity(art))
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.UpdateStatus(ctx, ToArticleEntity(art))
}

func ToArticleEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}

func ToArticleDomain(art article.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Status: domain.ArticleStatus(art.Status),
	}
}
