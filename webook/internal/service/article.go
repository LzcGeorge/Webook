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
	Save(ctx context.Context, art domain.Article) (int64, error)
	SaveWithTwoRepo(ctx context.Context, art domain.Article) (int64, error)
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

// Save 保存到线上库： 返回文章 id
func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	//if article.Id > 0 {
	//	return a.repo.Update(ctx, article)
	//}
	//return a.repo.Create(ctx, article)

	return a.SaveWithTwoRepo(ctx, article)
}

func (a *articleService) SaveWithTwoRepo(ctx context.Context, art domain.Article) (int64, error) {
	id := art.Id
	var err error

	// 写者库更新
	if id > 0 {
		id, err = a.authorRepo.Update(ctx, art)
	} else {
		id, err = a.authorRepo.Create(ctx, art)
	}

	if err != nil {
		a.logger.Error("authorRepo create article failed",
			logger.Int64("article id: ", art.Id),
			logger.Int64("author id: ", art.Author.Id),
			logger.Error(err),
		)
		return 0, errors.New("authorRepo create article failed, " + err.Error())
	}

	// 线上库更新
	// 类似于 FindOrCreate 中的实现，先查询线上库是否存在，不存在则创建，存在则更新
	res, err := a.readerRepo.FindById(ctx, art.Id)
	if err != nil {
		a.logger.Error("find article by id failed",
			logger.Int64("article id: ", art.Id),
			logger.Error(err),
		)
		return 0, err
	}

	// 线上库的最小 id 是 1，则说明文章不存在，创建文章
	if res.Id < 1 {
		return a.readerRepo.Create(ctx, art)
	}

	// 线上库存在，则更新
	return a.readerRepo.Update(ctx, art)
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
		id, err = a.saveArticle(ctx, art)
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

func (a *articleService) saveArticle(ctx context.Context, art domain.Article) (int64, error) {
	res, err := a.readerRepo.FindById(ctx, art.Id)
	if err != nil {
		return 0, err
	}

	if res.Id < 1 {
		return a.readerRepo.Create(ctx, art)
	}

	return a.readerRepo.Update(ctx, art)
}
