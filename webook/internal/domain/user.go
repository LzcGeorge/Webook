package domain

import "time"

// User 用户领域模型
type User struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
}
