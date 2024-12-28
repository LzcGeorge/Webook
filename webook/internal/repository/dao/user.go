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
	ErrUserNotFound       = gorm.ErrRecordNotFound
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
			// 邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err

}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	// 按 email 查询，查询到的第一条数据绑定到 user 中
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}
