package ioc

import (
	"Webook/webook/internal/repository/dao"
	"Webook/webook/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

// InitDB 初始化数据库
func InitDB(l logger.Logger) *gorm.DB {
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

	db, err := gorm.Open(mysql.Open(dbConfig.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询阈值，只有执行时间超过这个阈值，才会使用
			// 50ms， 100ms
			// SQL 查询必然要求命中索引，最好就是走一次磁盘 IO
			// 一次磁盘 IO 是不到 10ms
			SlowThreshold:             time.Millisecond * 10,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  glogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}
