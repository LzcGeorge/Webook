package article

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ArticleReaderDAO interface {
	Insert(ctx context.Context, art PublishedArticle) (int64, error)
	UpdateById(ctx context.Context, art PublishedArticle) (int64, error)
	FindById(ctx context.Context, id int64) (PublishedArticle, error)
}

type GormArticleReaderDAO struct {
	db *gorm.DB
}

func NewGormArticleReaderDAO(db *gorm.DB) ArticleReaderDAO {
	return &GormArticleReaderDAO{
		db: db,
	}
}

func (dao *GormArticleReaderDAO) Insert(ctx context.Context, art PublishedArticle) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

// Update 更新文章的标题和内容
func (dao *GormArticleReaderDAO) UpdateById(ctx context.Context, art PublishedArticle) (int64, error) {
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

func (dao *GormArticleReaderDAO) FindById(ctx context.Context, id int64) (PublishedArticle, error) {
	var art PublishedArticle
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&art).Error
	if err == gorm.ErrRecordNotFound {
		return PublishedArticle{
			Article: Article{
				Id: 0, // 表示不存在
			},
		}, nil
	}
	return art, err
}
