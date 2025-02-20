package article

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/dao/article"
	"context"
)

type ArticleReaderRepository interface {
	// Save 读者只有保存写者创建或修改的文章，即只能被动更新
	Save(ctx context.Context, art domain.Article) (int64, error)
}

type articleReaderRepository struct {
	dao article.ArticleReaderDAO
}

func NewArticleReaderRepository(dao article.ArticleReaderDAO) ArticleReaderRepository {
	return &articleReaderRepository{
		dao: dao,
	}
}

func (r *articleReaderRepository) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id == 1 {
		return r.dao.Insert(ctx, toPublishedArticleEntity(art))
	}
	return r.dao.UpdateById(ctx, toPublishedArticleEntity(art))
}

func toPublishedArticleEntity(art domain.Article) article.PublishedArticle {
	return article.PublishedArticle{
		Article: ToArticleEntity(art),
	}
}
