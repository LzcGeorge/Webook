package article

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OSSArticleDAO struct {
	GormArticleDAO
	oss *s3.S3
}

func NewOSSArticleDAO(db *gorm.DB, oss *s3.S3) *OSSArticleDAO {
	return &OSSArticleDAO{
		GormArticleDAO: GormArticleDAO{db: db},
		oss:            oss,
	}
}

type PublishedArticleOSS struct {
	Id       int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title    string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	AuthorId int64  `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}

func (a *OSSArticleDAO) Upsert(ctx context.Context, art Article) (int64, error) {
	var id = art.Id

	// 使用事务，保证写者库和读者库的一致性
	err := a.db.WithContext(ctx).Transaction(func(txDb *gorm.DB) error {
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
		pubArt := PublishedArticleOSS{
			Id:       id,
			Title:    art.Title,
			AuthorId: art.AuthorId,
			Status:   art.Status,
			Ctime:    now,
			Utime:    now,
		}
		err = txDb.Clauses(clause.OnConflict{
			// id 冲突的时候执行 update，否则执行 insert
			Columns: []clause.Column{{Name: "id"}},
			// update 的时候，只更新 title 和 content, utime
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":  art.Title,
				"status": art.Status,
				"utime":  now,
			}),
		}).Create(&pubArt).Error
		return err
	})

	// 写入 OSS
	if err != nil {
		return 0, err
	}
	_, err = a.oss.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("webook-oss"),
		Key:         aws.String(strconv.FormatInt(id, 10)),
		Body:        bytes.NewReader([]byte(art.Content)),
		ContentType: aws.String("text/plain; charset=utf-8"),
	})
	return id, err
}

func (a *OSSArticleDAO) UpdateStatus(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()

	err := a.db.WithContext(ctx).Transaction(func(txDb *gorm.DB) error {
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
		pubArt := PublishedArticleOSS{
			Id:       art.Id,
			Title:    art.Title,
			AuthorId: art.AuthorId,
			Status:   art.Status,
			Ctime:    now,
			Utime:    now,
		}
		return txDb.Model(&pubArt).
			Where("id = ? and author_id = ?", pubArt.Id, pubArt.AuthorId).
			Updates(map[string]any{
				"status": art.Status,
				"utime":  now,
			}).Error
	})

	const statusPrivate = 3
	// 如果 status 是私密，则删除 OSS 上的文件
	// 其他状态，文章内容没有改变，所以不需要更新 OSS 上的文件
	if art.Status == statusPrivate {
		_, err = a.oss.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String("webook-oss"),
			Key:    aws.String(strconv.FormatInt(art.Id, 10)),
		})
	}
	return art.Id, err
}
