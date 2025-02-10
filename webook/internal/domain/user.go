package domain

import "time"

// User 用户领域模型
type User struct {
	Id       int64
	Email    string
	Phone    string
	Password string
	Ctime    time.Time

	// 用户信息
	Nickname string
	Birthday time.Time
	AboutMe  string

	// 微信信息
	WechatInfo WechatInfo
}
