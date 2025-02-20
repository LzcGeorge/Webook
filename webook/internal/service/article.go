package service

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/article"
	"Webook/webook/pkg/logger"
	"context"
	"errors"
	"time"
)

type ArticleService interface {
	Edit(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishWithTwoRepo(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	// 一个 Service 操作一个 Repo：读者写者共用一个库
	repo article.ArticleRepository

	// 一个 Service 操作两个 Repo：读者库，写者库
	authorRepo article.ArticleAuthorRepository
	readerRepo article.ArticleReaderRepository

	// logger
	logger logger.Logger
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func NewArticleServiceWithTwoRepo(authorRepo article.ArticleAuthorRepository, readerRepo article.ArticleReaderRepository, logger logger.Logger) ArticleService {
	return &articleService{
		authorRepo: authorRepo,
		readerRepo: readerRepo,
		logger:     logger,
	}
}

// Edit 编辑文章： 返回文章 id
func (a *articleService) Edit(ctx context.Context, article domain.Article) (int64, error) {
	//if article.Id > 0 {
	//	return a.repo.Update(ctx, article)
	//}
	//return a.repo.Create(ctx, article)

	return a.EditWithTwoRepo(ctx, article)
}

func (a *articleService) EditWithTwoRepo(ctx context.Context, art domain.Article) (int64, error) {
	id := art.Id
	var err error
	if id > 0 {
		id, err = a.authorRepo.Update(ctx, art)
	} else {
		id, err = a.authorRepo.Create(ctx, art)
	}

	if err != nil {
		err = errors.New("authorRepo create article failed")
	}

	// 线上库更新
	return a.readerRepo.Save(ctx, art)
}

// Publish 发布文章
func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {

	return a.PublishWithTwoRepo(ctx, art)
}

// PublishWithTwoRepo 采用读者库和写者库
func (a *articleService) PublishWithTwoRepo(ctx context.Context, art domain.Article) (int64, error) {
	// 写者库发表文章
	var id = art.Id
	var err error
	if art.Id > 0 {
		id, err = a.authorRepo.Update(ctx, art)
	} else {
		id, err = a.authorRepo.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	// 确保写者库和读者库的 id 一致
	art.Id = id

	// 读者库保存文章，如果失败，则重试, 重试至多 3 次
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * time.Duration(i))
		id, err = a.readerRepo.Save(ctx, art)
		if err == nil {
			break
		}
		a.logger.Error("save article to reader repo failed, try again",
			logger.Int64("article id: ", art.Id),
			logger.Int64("author id: ", art.Author.Id),
			logger.Error(err),
		)
	}

	if err != nil {
		// 重试 3 次仍然失败，则返回错误
		a.logger.Error("reader repo save art failed",
			logger.Int64("art id: ", art.Id),
			logger.Error(err),
		)
	}

	return id, nil
}
