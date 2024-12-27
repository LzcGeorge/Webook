package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
)

type UserDAO struct {
	db *gorm.DB
}

// dao.User 对应数据库中的 user 表
type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	// 创建和修改时间，毫秒时间戳
	Ctime int64
	Utime int64
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now

	// 调用 gorm 进行插入
	err := dao.db.WithContext(ctx).Create(&user).Error

	// 如果插入失败，返回错误
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrNo uint16 = 1062
		if mysqlErr.Number == uniqueIndexErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err

}
