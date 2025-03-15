package article

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository"
	"Webook/webook/internal/repository/cache"
	"Webook/webook/internal/repository/dao/article"
	"Webook/webook/pkg/logger"
	"context"
	"errors"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) (int64, error)
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, art domain.Article) (int64, error)
	List(ctx context.Context, userId int64, limit int, offset int) ([]domain.Article, error)
	FindById(ctx context.Context, id int64) (domain.Article, error)
	FindPublishedArticleById(ctx context.Context, id int64) (domain.Article, error)
	FindPublishedArticleList(ctx context.Context, end time.Time, offset int, limit int) ([]domain.Article, error)
}

type CachedArticleRepository struct {
	dao   article.ArticleDAO
	cache cache.ArticleCache
	// 查询用户信息
	userRepo repository.UserRepository
	logger   logger.Logger
}

func NewArticleRepository(dao article.ArticleDAO, cache cache.ArticleCache, userRepo repository.UserRepository, logger logger.Logger) ArticleRepository {
	return &CachedArticleRepository{
		dao:      dao,
		cache:    cache,
		userRepo: userRepo,
		logger:   logger,
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
		if err := c.cache.DelPublic(ctx, art.Id); err != nil {
			c.logger.Error("Create Article 后删除缓存 Public Article 失败",
				logger.Int64("articleId", art.Id),
				logger.Error(err),
			)
		}

		if err := c.cache.Del(ctx, art.Id); err != nil {
			c.logger.Error("Create Article 后删除缓存 article 失败",
				logger.Int64("articleId", art.Id),
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
		if err := c.cache.DelPublic(ctx, art.Id); err != nil {
			c.logger.Error("Update Article 后删除缓存 Public Article 失败",
				logger.Int64("articleId", art.Id),
				logger.Error(err),
			)
		}

		if err := c.cache.Del(ctx, art.Id); err != nil {
			c.logger.Error("Update Article 后删除缓存 article 失败",
				logger.Int64("articleId", art.Id),
				logger.Error(err),
			)
		}
	}()
	return c.dao.UpdateById(ctx, ToArticleEntity(art))
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	// 数据修改后删除缓存
	defer func() {
		if err := c.cache.DelFirstPage(ctx, art.Author.Id); err != nil {
			c.logger.Error("Sync Article 后删除缓存 FirstPage 失败",
				logger.Int64("userId", art.Author.Id),
				logger.Error(err),
			)
		}

		if err := c.cache.DelPublic(ctx, art.Id); err != nil {
			c.logger.Error("Sync Article 后删除缓存 Public Article  失败",
				logger.Int64("articleId", art.Id),
				logger.Error(err),
			)
		}

		if err := c.cache.Del(ctx, art.Id); err != nil {
			c.logger.Error("Sync Article 后删除缓存 article 失败",
				logger.Int64("articleId", art.Id),
				logger.Error(err),
			)
		}
	}()
	return c.dao.Upsert(ctx, ToArticleEntity(art))
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, art domain.Article) (int64, error) {
	// 数据修改后删除缓存
	defer func() {
		if err := c.cache.DelFirstPage(ctx, art.Author.Id); err != nil {
			c.logger.Error("SyncStatus Article 后删除缓存 FirstPage 失败",
				logger.Int64("userId", art.Author.Id),
				logger.Error(err),
			)
		}

		if err := c.cache.DelPublic(ctx, art.Id); err != nil {
			c.logger.Error("SyncStatus Article 后删除缓存 Public Article  失败",
				logger.Int64("articleId", art.Id),
				logger.Error(err),
			)
		}

		if err := c.cache.Del(ctx, art.Id); err != nil {
			c.logger.Error("SyncStatus Article 后删除缓存 article 失败",
				logger.Int64("articleId", art.Id),
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
			// 预缓存列表中的第一篇文章
			go func() {
				c.preCache(ctx, cachedArts)
			}()
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

	// 预缓存列表中的第一篇文章
	go func() {
		c.preCache(ctx, result)
	}()
	return result, nil
}

func (c *CachedArticleRepository) preCache(ctx context.Context, arts []domain.Article) {
	const detailExpire = time.Minute
	const contentSizeThreshold = 1024 * 1024 // 1MB

	if len(arts) > 0 && len(arts[0].Content) < contentSizeThreshold {
		art := arts[0]
		if err := c.cache.Set(ctx, art); err != nil {
			c.logger.Error("预缓存第一篇文章失败", logger.Error(err))
		}
	}
}
func (c *CachedArticleRepository) FindById(ctx context.Context, id int64) (domain.Article, error) {
	// 从缓存中获取
	res, err := c.cache.Get(ctx, id)
	if err == nil {
		return res, nil
	}

	// 缓存未命中，从数据库中获取
	art, err := c.dao.FindById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}

	// 缓存该文章
	domainArt := ToArticleDomain(art)
	defer func() {
		err := c.cache.Set(ctx, domainArt)
		if err != nil {
			c.logger.Error("缓存文章失败",
				logger.Int64("id", id),
				logger.Error(err),
			)
		}
	}()
	return domainArt, nil
}

func (c *CachedArticleRepository) FindPublishedArticleById(ctx context.Context, id int64) (domain.Article, error) {
	// 从缓存中获取
	res, err := c.cache.GetPublic(ctx, id)
	if err == nil {
		return res, nil
	}

	// 缓存未命中，从数据库中获取
	artPublic_published, err := c.dao.FindPublicById(ctx, id)
	if err == gorm.ErrRecordNotFound {
		return domain.Article{}, errors.New("文章不存在或未发表")
	}
	if err != nil {
		return domain.Article{}, err
	}

	// 获取作者信息
	artPublic := ToArticleDomain(artPublic_published.Article)
	author, err := c.userRepo.FindById(ctx, artPublic.Author.Id)
	if err != nil {
		return domain.Article{}, err
	}
	artPublic.Author.Name = author.Nickname

	// 缓存该文章
	go func() {
		err := c.cache.SetPublic(ctx, artPublic)
		if err != nil {
			c.logger.Error("缓存Public文章失败",
				logger.Int64("id", id),
				logger.Error(err),
			)
		}
	}()
	return artPublic, nil
}

func (c *CachedArticleRepository) FindPublishedArticleList(ctx context.Context, end time.Time, offset int, limit int) ([]domain.Article, error) {
	arts, err := c.dao.FindPublishedArticleList(ctx, end, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map[article.PublishedArticle, domain.Article](arts,
		func(idx int, src article.PublishedArticle) domain.Article {
			return ToArticleDomain(src.Article)
		}), nil
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
