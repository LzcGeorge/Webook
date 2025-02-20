package article

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/dao/article"
	"context"
)

type ArticleReaderRepository interface {
	// Save 读者只有保存写者创建或修改的文章，即只能被动更新
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) (int64, error)
	FindById(ctx context.Context, id int64) (domain.Article, error)
}

type articleReaderRepository struct {
	dao article.ArticleReaderDAO
}

func NewArticleReaderRepository(dao article.ArticleReaderDAO) ArticleReaderRepository {
	return &articleReaderRepository{
		dao: dao,
	}
}

func (r *articleReaderRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return r.dao.Insert(ctx, toPublishedArticleEntity(art))
}

func (r *articleReaderRepository) Update(ctx context.Context, art domain.Article) (int64, error) {
	return r.dao.UpdateById(ctx, toPublishedArticleEntity(art))
}

func toPublishedArticleEntity(art domain.Article) article.PublishedArticle {
	return article.PublishedArticle{
		Article: ToArticleEntity(art),
	}
}

func (r *articleReaderRepository) FindById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return ToArticleDomain(art.Article), nil
}
