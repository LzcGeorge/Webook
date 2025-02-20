package article

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// 作者库：author 进行写入和更新，删除。。
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

// 线上库：reader 进行被动更新
type PublishedArticle struct {
	Article
}

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) (int64, error)
}

type GormArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GormArticleDAO{
		db: db,
	}
}

func (dao *GormArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

// Update 更新文章的标题和内容
func (dao *GormArticleDAO) UpdateById(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := dao.db.WithContext(ctx).Model(&art).
		Where("id = ?", art.Id).
		Where("author_id = ?", art.AuthorId).
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
