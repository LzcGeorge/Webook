package dao

import (
	"context"
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
	err := dao.db.WithContext(ctx).Model(&article).
		Where("id = ?", article.Id).
		Where("author_id = ?", article.AuthorId).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"utime":   now,
		}).Error
	return article.Id, err
}
