package dao

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 交互表：点赞、收藏、阅读，区分业务（biz, bizid）
type Interactive struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// <bizid, biz>
	BizId int64 `gorm:"uniqueIndex:biz_type_id"`
	// WHERE biz = ?
	Biz string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`

	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Ctime      int64
	Utime      int64
}

// 用户点赞业务表
type UserLikeBiz struct {
	Id     int64  `gorm:"primaryKey,autoIncrement"`
	Uid    int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizId  int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz    string `gorm:"type:varchar(128);uniqueIndex:uid_biz_type_id"`
	Status int    // 1 表示有效，0 表示已经软删除
	Ctime  int64
	Utime  int64
}

// 用户收藏业务表
type UserCollectBiz struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:uid_biz_type_id"`
	BizId int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Uid   int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Cid   int64  `gorm:"index"`
	Ctime int64
	Utime int64
}

type InteractiveDAO interface {
	IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error
	DeleteLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error
	InsertCollection(ctx context.Context, biz string, bizId int64, collectionId int64, userId int64) error
	GetInteractive(ctx context.Context, biz string, bizId int64) (Interactive, error)
	GetLiked(ctx context.Context, biz string, bizId int64, userId int64) (bool, error)
	GetCollected(ctx context.Context, biz string, bizId int64, userId int64) (bool, error)
}

type GormInteractiveDAO struct {
	db *gorm.DB
}

func NewInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GormInteractiveDAO{
		db: db,
	}
}

// 阅读计数: Upsert 语句
func (dao *GormInteractiveDAO) IncreaseReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"read_cnt": gorm.Expr("`read_cnt` + 1"),
			"utime":    now,
		}),
	}).Create(&Interactive{
		Biz:     biz,
		BizId:   bizId,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}

func (dao *GormInteractiveDAO) InsertLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		er := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"utime": now,
			}),
		}).Create(&UserLikeBiz{
			Uid:    userId,
			Biz:    biz,
			BizId:  bizId,
			Status: 1,
			Ctime:  now,
			Utime:  now,
		}).Error
		if er != nil {
			return er
		}

		// 维护 redis 中的 map 结构
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"like_cnt": gorm.Expr("`like_cnt` + 1"),
				"utime":    now,
			}),
		}).Create(&Interactive{
			Biz:     biz,
			BizId:   bizId,
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
		}).Error
	})
}

func (dao *GormInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId int64, userId int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).
			Where("uid=? AND biz_id = ? AND biz=?", userId, bizId, biz).
			Updates(map[string]interface{}{
				"utime":  now,
				"status": 0,
			}).Error
		if err != nil {
			return err
		}
		return tx.Model(&Interactive{}).
			Where("biz =? AND biz_id=?", biz, bizId).
			Updates(map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt` - 1"),
				"utime":    now,
			}).Error
	})
}

func (dao *GormInteractiveDAO) InsertCollection(ctx context.Context, biz string, bizId int64, collectionId int64, userId int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&UserCollectBiz{
			Uid:   userId,
			Biz:   biz,
			BizId: bizId,
			Cid:   collectionId,
			Ctime: now,
			Utime: now,
		}).Error
		if err != nil {
			return err
		}

		// 维护 redis 中的 map 结构
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"collect_cnt": gorm.Expr("`collect_cnt` + 1"),
				"utime":       now,
			}),
		}).Create(&Interactive{
			Biz:        biz,
			BizId:      bizId,
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
		}).Error
	})
}

func (dao *GormInteractiveDAO) GetInteractive(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	var res Interactive
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ?", biz, bizId).First(&res).Error
	return res, err
}

func (dao *GormInteractiveDAO) GetLiked(ctx context.Context, biz string, bizId int64, userId int64) (bool, error) {
	var like UserLikeBiz
	err := dao.db.WithContext(ctx).Where("uid = ? AND biz = ? AND biz_id = ?", userId, biz, bizId).First(&like).Error
	if err == gorm.ErrRecordNotFound {
		// 记录不存在，说明用户没有点赞，返回 false 且没有错误
		return false, nil
	}
	// 其他错误则返回错误
	if err != nil {
		return false, err
	}
	// 记录存在，根据 status 判断是否点赞
	return like.Status == 1, nil
}

// 获取用户是否收藏，如何用 sql 查询呢？
func (dao *GormInteractiveDAO) GetCollected(ctx context.Context, biz string, bizId int64, userId int64) (bool, error) {
	var collect UserCollectBiz
	// 查询是否存在,查询不到返回什么？
	err := dao.db.WithContext(ctx).Where("uid = ? AND biz = ? AND biz_id = ?", userId, biz, bizId).First(&collect).Error
	fmt.Printf("err : %v\n", err)
	fmt.Printf(" err != gorm.ErrRecordNotFound : %v\n", err != gorm.ErrRecordNotFound)
	return err != gorm.ErrRecordNotFound, err
}
