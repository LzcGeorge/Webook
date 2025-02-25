package domain

import "time"

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
	Ctime   time.Time
	Utime   time.Time
}

type Author struct {
	Id   int64
	Name string
}

type ArticleStatus uint8

const (
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
	ArticleStatusArchived // 已删除
)

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

func (a Article) Abstract() string {
	content := []rune(a.Content)
	if len(content) < 100 {
		return string(content)
	}
	return string(content[:100])
}
