package article

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/cache"
	"Webook/webook/internal/repository/dao/article"
	"Webook/webook/pkg/logger"
	"context"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) (int64, error)
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, art domain.Article) (int64, error)
	List(ctx context.Context, userId int64, limit int, offset int) ([]domain.Article, error)
}

type CachedArticleRepository struct {
	dao    article.ArticleDAO
	cache  cache.ArticleCache
	logger logger.Logger
}

func NewArticleRepository(dao article.ArticleDAO, cache cache.ArticleCache, logger logger.Logger) ArticleRepository {
	return &CachedArticleRepository{
		dao:    dao,
		cache:  cache,
		logger: logger,
	}
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	// 数据修改后删除缓存
	defer func() {
		err := c.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			c.logger.Error("Create Article 后删除缓存失败",
				logger.Int64("userId", art.Author.Id),
				logger.Error(err),
			)
		}
	}()
	return c.dao.Insert(ctx, ToArticleEntity(art))
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) (int64, error) {
	// 数据修改后删除缓存
	defer func() {
		err := c.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			c.logger.Error("Update Article 后删除缓存失败",
				logger.Int64("userId", art.Author.Id),
				logger.Error(err),
			)
		}
	}()
	return c.dao.UpdateById(ctx, ToArticleEntity(art))
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	// 数据修改后删除缓存
	defer func() {
		err := c.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			c.logger.Error("Sync Article 后删除缓存失败",
				logger.Int64("userId", art.Author.Id),
				logger.Error(err),
			)
		}
	}()
	return c.dao.Upsert(ctx, ToArticleEntity(art))
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, art domain.Article) (int64, error) {
	// 数据修改后删除缓存
	defer func() {
		err := c.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			c.logger.Error("SyncStatus Article 后删除缓存失败",
				logger.Int64("userId", art.Author.Id),
				logger.Error(err),
			)
		}
	}()
	return c.dao.UpdateStatus(ctx, ToArticleEntity(art))
}

func (c *CachedArticleRepository) List(ctx context.Context, userId int64, limit int, offset int) ([]domain.Article, error) {
	// 如果是第一页，从缓存中获取
	if offset == 0 && limit <= 100 {
		cachedArts, err := c.cache.GetFirstPage(ctx, userId)
		if err == nil {
			// 缓存命中
			c.logger.Info("缓存命中",
				logger.Int64("userId", userId),
			)
			return cachedArts, nil
		}
	}

	// 缓存未命中，从数据库中获取
	arts, err := c.dao.GetByAuthorId(ctx, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Article, 0, len(arts))
	for _, art := range arts {
		result = append(result, ToArticleDomain(art))
	}

	// 异步缓存第一页的数据
	go func() {
		// 设置缓存超时时间
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// 缓存第一页的数据
		if offset == 0 && limit <= 100 {
			err := c.cache.SetFirstPage(ctx, userId, result)
			if err != nil {
				c.logger.Error("缓存第一页的数据失败",
					logger.Int64("userId", userId),
					logger.Int64("limit", int64(limit)),
					logger.Int64("offset", int64(offset)),
					logger.Error(err),
				)
			}
		}
	}()
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
