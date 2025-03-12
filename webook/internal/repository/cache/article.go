package cache

import (
	"Webook/webook/config"
	"Webook/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"

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
	Del(ctx context.Context, id int64) error

	// 缓存Public文章
	SetPublic(ctx context.Context, art domain.Article) error
	GetPublic(ctx context.Context, id int64) (domain.Article, error)
	DelPublic(ctx context.Context, id int64) error
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
	dur := config.DevRedisExpire.ArticleFirstPage
	return c.client.Set(ctx, c.key(userId), jsonData, dur).Err()
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
	if err == nil {
		c.client.Expire(ctx, c.key(userId), config.DevRedisExpire.ArticleFirstPage)
	}
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

	dur := config.DevRedisExpire.ArticleDetail
	return c.client.Set(ctx, c.detailKey(art.Id), jsonData, dur).Err()
}

func (c *RedisArticleCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	jsonData, err := c.client.Get(ctx, c.detailKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var art domain.Article
	err = json.Unmarshal(jsonData, &art)
	if err == nil {
		c.client.Expire(ctx, c.detailKey(id), config.DevRedisExpire.ArticleDetail)
	}
	return art, err
}

func (c *RedisArticleCache) Del(ctx context.Context, id int64) error {
	return c.client.Del(ctx, c.detailKey(id)).Err()
}
func (c *RedisArticleCache) PublicKey(id int64) string {
	return fmt.Sprintf("article:public:%d", id)
}

func (c *RedisArticleCache) SetPublic(ctx context.Context, art domain.Article) error {
	jsonData, err := json.Marshal(art)
	if err != nil {
		return err
	}

	dur := config.DevRedisExpire.PublicArticle
	return c.client.Set(ctx, c.PublicKey(art.Id), jsonData, dur).Err()
}

func (c *RedisArticleCache) GetPublic(ctx context.Context, id int64) (domain.Article, error) {
	jsonData, err := c.client.Get(ctx, c.PublicKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}

	var art domain.Article
	err = json.Unmarshal(jsonData, &art)
	if err == nil {
		c.client.Expire(ctx, c.PublicKey(id), config.DevRedisExpire.PublicArticle)
	}

	return art, err
}

func (c *RedisArticleCache) DelPublic(ctx context.Context, id int64) error {
	return c.client.Del(ctx, c.PublicKey(id)).Err()
}
