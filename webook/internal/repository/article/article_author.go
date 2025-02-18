package article

import (
	"Webook/webook/internal/domain"
	"context"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) (int64, error)
}
