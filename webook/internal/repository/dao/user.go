package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicate = errors.New("邮箱或手机号冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByWechat(ctx context.Context, openId string) (User, error)
	UpdateById(ctx context.Context, user User) error
}

type GormUserDAO struct {
	db *gorm.DB
}

// dao.User 对应数据库中的 user 表
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`

	// 唯一索引允许有多个空值，但不能有多个 ""
	Email sql.NullString `gorm:"unique"`
	Phone sql.NullString `gorm:"unique"`

	Password string
	// 创建和修改时间，毫秒时间戳
	Ctime int64
	Utime int64

	// 用户信息
	Nickname string
	Birthday time.Time `gorm:"default:1992-02-03 00:00:00.000"`
	AboutMe  string

	// 微信信息
	WechatOpenId  sql.NullString `gorm:"unique"`
	WechatUnionId sql.NullString
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GormUserDAO{
		db: db,
	}
}

func (dao *GormUserDAO) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now

	// 调用 gorm 进行插入
	err := dao.db.WithContext(ctx).Create(&user).Error

	// 如果插入失败，返回错误
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrNo uint16 = 1062
		if mysqlErr.Number == uniqueIndexErrNo {
			// 邮箱冲突 or 手机号冲突
			return ErrUserDuplicate
		}
	}
	return err

}

func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	// 按 email 查询，查询到的第一条数据绑定到 user 中
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (dao *GormUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

func (dao *GormUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}

func (dao *GormUserDAO) UpdateById(ctx context.Context, user User) error {
	user.Utime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&user).Where("id = ?", user.Id).Updates(user).Error
}

func (dao *GormUserDAO) FindByWechat(ctx context.Context, openId string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("wechat_open_id = ?", openId).First(&user).Error
	return user, err
}
