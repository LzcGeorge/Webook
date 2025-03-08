package dao

import (
	"Webook/webook/internal/repository/dao/article"

	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &article.Article{}, &article.PublishedArticle{}, &Interactive{}, &UserLikeBiz{}, &UserCollectBiz{})
}

func TruncateTable(db *gorm.DB, tableName string) error {
	return db.Exec("TRUNCATE TABLE " + tableName).Error
}
