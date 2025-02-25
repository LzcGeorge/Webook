package article

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/dao/article"
	"context"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) (int64, error)
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, art domain.Article) (int64, error)
	GetByAuthorId(ctx context.Context, userId int64, limit int, offset int) ([]domain.Article, error)
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

func (c *CachedArticleRepository) GetByAuthorId(ctx context.Context, userId int64, limit int, offset int) ([]domain.Article, error) {
	arts, err := c.dao.GetByAuthorId(ctx, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Article, 0, len(arts))
	for _, art := range arts {
		result = append(result, ToArticleDomain(art))
	}
	return result, nil
}

func ToArticleEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
		Ctime:    art.Ctime.UnixMilli(),
		Utime:    art.Utime.UnixMilli(),
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
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
	}
}
