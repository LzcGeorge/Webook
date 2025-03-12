package article

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 作者库：author 进行写入和更新，删除。。
type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title   string `gorm:"type=varchar(1024)" bson:"title,omitempty"`
	Content string `gorm:"type=BLOB" bson:"content,omitempty"`

	// 作者 id, 在 author_id 上建立索引
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`

	// 状态
	Status uint8 `bson:"status,omitempty"`

	// 创建和修改时间，毫秒时间戳
	Ctime int64 `bson:"ctime,omitempty"`
	Utime int64 `bson:"utime,omitempty"`
}

// 线上库：reader 进行被动更新
type PublishedArticle struct {
	Article
}

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) (int64, error)
	Upsert(ctx context.Context, art Article) (int64, error)
	UpdateStatus(ctx context.Context, art Article) (int64, error)
	GetByAuthorId(ctx context.Context, userId int64, limit int, offset int) ([]Article, error)
	FindById(ctx context.Context, id int64) (Article, error)
	FindPublicById(ctx context.Context, id int64) (PublishedArticle, error)
	FindPublishedArticleList(ctx context.Context, end time.Time, offset int, limit int) ([]PublishedArticle, error)
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
			"status":  art.Status,
			"utime":   now,
		})

	// 至少会有更新时间 Utime 会被更新，所以可以判断是否更新成功
	if res.RowsAffected == 0 {
		return art.Id, errors.New("可能是别人写的文章，或者已经删除了")
	}
	return art.Id, res.Error
}

func (dao *GormArticleDAO) Upsert(ctx context.Context, art Article) (int64, error) {
	var id = art.Id

	// 使用事务，保证写者库和读者库的一致性
	err := dao.db.WithContext(ctx).Transaction(func(txDb *gorm.DB) error {
		var err error
		now := time.Now().UnixMilli()
		// 写者库,
		dao := NewArticleDAO(txDb)
		if id > 0 {
			id, err = dao.UpdateById(ctx, art)
		} else {
			id, err = dao.Insert(ctx, art)
		}

		if err != nil {
			return err
		}

		// 读者库：Upsert 即 update or insert
		art.Id = id
		pubArt := PublishedArticle{
			Article: art,
		}
		pubArt.Ctime = now
		pubArt.Utime = now
		err = txDb.Clauses(clause.OnConflict{
			// id 冲突的时候执行 update，否则执行 insert
			Columns: []clause.Column{{Name: "id"}},
			// update 的时候，只更新 title 和 content, utime
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   art.Title,
				"content": art.Content,
				"status":  art.Status,
				"utime":   now,
			}),
		}).Create(&pubArt).Error
		return err
	})
	return id, err
}

func (dao *GormArticleDAO) UpdateStatus(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()

	err := dao.db.WithContext(ctx).Transaction(func(txDb *gorm.DB) error {
		res := txDb.Model(&art).
			Where("id = ? and author_id = ?", art.Id, art.AuthorId).
			Updates(map[string]any{
				"status": art.Status,
				"utime":  now,
			})

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return errors.New("可能是别人写的文章")
		}

		// 读者库：更新 status
		pubArt := PublishedArticle{
			Article: art,
		}
		pubArt.Utime = now
		return txDb.Model(&pubArt).
			Where("id = ? and author_id = ?", pubArt.Id, pubArt.AuthorId).
			Updates(map[string]any{
				"status": art.Status,
				"utime":  now,
			}).Error
	})
	return art.Id, err
}

func (dao *GormArticleDAO) GetByAuthorId(ctx context.Context, userId int64, limit int, offset int) ([]Article, error) {
	var arts []Article
	err := dao.db.WithContext(ctx).Where("author_id = ?", userId).
		Order("ctime desc").
		Limit(limit).Offset(offset).Find(&arts).Error
	return arts, err
}

func (dao *GormArticleDAO) FindById(ctx context.Context, id int64) (Article, error) {
	var art Article
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&art).Error
	return art, err
}

// FindPublicById 从线上库获取文章，只能查询状态为 ArticleStatusPublished：2 的文章
func (dao *GormArticleDAO) FindPublicById(ctx context.Context, id int64) (PublishedArticle, error) {
	var art PublishedArticle
	err := dao.db.WithContext(ctx).Where("id = ? and status = ? ", id, 2).First(&art).Error
	return art, err
}

// FindPublishedArticleList 获取线上库文章列表
func (dao *GormArticleDAO) FindPublishedArticleList(ctx context.Context, end time.Time, offset int, limit int) ([]PublishedArticle, error) {
	var arts []PublishedArticle
	err := dao.db.WithContext(ctx).Where("status = ? and utime < ?", 2, end.UnixMilli()).
		Order("utime desc").
		Limit(limit).Offset(offset).Find(&arts).Error
	return arts, err
}
