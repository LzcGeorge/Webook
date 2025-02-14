package ioc

import (
	"Webook/webook/internal/repository/dao"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	// 利用 viper 读取配置文件
	type DBConfig struct {
		DSN string `yaml:"DSN"`
	}
	var dbConfig DBConfig = DBConfig{
		// 默认值
		DSN: "default value",
	}
	// 从配置文件中读取 db 配置，会覆盖默认值
	err := viper.UnmarshalKey("db", &dbConfig)
	if err != nil {
		panic(err)
	}

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
