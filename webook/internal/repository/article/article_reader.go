package article

import (
	"Webook/webook/internal/domain"
	"context"
)

type ArticleReaderRepository interface {
	// Save 读者只有保存写者创建或修改的文章，即只能被动更新
	Save(ctx context.Context, article domain.Article) (int64, error)
}
