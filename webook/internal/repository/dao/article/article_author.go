package article

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ArticleAuthorDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) (int64, error)
}

type GormArticleAuthorDAO struct {
	db *gorm.DB
}

func NewGormArticleAuthorDAO(db *gorm.DB) ArticleAuthorDAO {
	return &GormArticleAuthorDAO{
		db: db,
	}
}

func (dao *GormArticleAuthorDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (dao *GormArticleAuthorDAO) UpdateById(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := dao.db.WithContext(ctx).Model(&art).
		Where("id = ? and author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   now,
		})

	// 至少会有更新时间 Utime 会被更新，所以可以判断是否更新成功
	if res.RowsAffected == 0 {
		return art.Id, errors.New("可能是别人写的文章，或者已经删除了")
	}
	return art.Id, res.Error
}
