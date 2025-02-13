package ioc

import (
	"Webook/webook/config"
	"Webook/webook/internal/repository/dao"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	var dbConfig = config.Config.DB // 也可以使用 k8s 的配置
	db, err := gorm.Open(mysql.Open(dbConfig.DSN))

	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
