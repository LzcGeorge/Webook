package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type=BLOB"`

	// 作者 id, 在 author_id 上建立索引
	AuthorId int64 `gorm:"index"`

	// 创建和修改时间，毫秒时间戳
	Ctime int64
	Utime int64
}

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	Update(ctx context.Context, article Article) (int64, error)
}

type GormArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GormArticleDAO{
		db: db,
	}
}

func (dao *GormArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := dao.db.WithContext(ctx).Create(&article).Error
	return article.Id, err
}

// Update 更新文章的标题和内容
func (dao *GormArticleDAO) Update(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.Utime = now
	res := dao.db.WithContext(ctx).Model(&article).
		Where("id = ?", article.Id).
		Where("author_id = ?", article.AuthorId).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"utime":   now,
		})

	// 至少会有更新时间 Utime 会被更新，所以可以判断是否更新成功
	if res.RowsAffected == 0 {
		return article.Id, errors.New("可能是别人写的文章，或者已经删除了")
	}
	return article.Id, res.Error
}
