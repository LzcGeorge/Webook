package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"Webook/webook/internal/domain"

	"github.com/redis/go-redis/v9"
)

type ArticleCache interface {
	// 缓存第一页
	SetFirstPage(ctx context.Context, userId int64, arts []domain.Article) error
	GetFirstPage(ctx context.Context, userId int64) ([]domain.Article, error)
	DelFirstPage(ctx context.Context, userId int64) error

	// 缓存文章
	Set(ctx context.Context, art domain.Article) error
	Get(ctx context.Context, id int64) (domain.Article, error)
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func NewRedisArticleCache(client redis.Cmdable) ArticleCache {
	return &RedisArticleCache{
		client: client,
	}
}

func (c *RedisArticleCache) key(userId int64) string {
	return fmt.Sprintf("article:first_page:%d", userId)
}

// SetFirstPage 设置第一页缓存, 缓存 10 分钟
func (c *RedisArticleCache) SetFirstPage(ctx context.Context, userId int64, arts []domain.Article) error {
	// 列表中只需保存摘要
	for i := 0; i < len(arts); i++ {
		arts[i].Content = arts[i].Abstract()
	}

	// 序列化
	jsonData, err := json.Marshal(arts)
	if err != nil {
		return err
	}

	// 设置缓存
	return c.client.Set(ctx, c.key(userId), jsonData, time.Minute*10).Err()
}

func (c *RedisArticleCache) GetFirstPage(ctx context.Context, userId int64) ([]domain.Article, error) {
	// 获取缓存
	jsonData, err := c.client.Get(ctx, c.key(userId)).Result()
	if err != nil {
		return nil, err
	}

	// 反序列化
	var arts []domain.Article
	err = json.Unmarshal([]byte(jsonData), &arts)

	return arts, err
}

func (c *RedisArticleCache) DelFirstPage(ctx context.Context, userId int64) error {
	return c.client.Del(ctx, c.key(userId)).Err()
}

func (c *RedisArticleCache) detailKey(id int64) string {
	return fmt.Sprintf("article:detail:%d", id)
}

func (c *RedisArticleCache) Set(ctx context.Context, art domain.Article) error {
	jsonData, err := json.Marshal(art)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.detailKey(art.Id), jsonData, time.Minute).Err()
}

func (c *RedisArticleCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	jsonData, err := c.client.Get(ctx, c.detailKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}

	var art domain.Article
	err = json.Unmarshal(jsonData, &art)
	return art, err
}
