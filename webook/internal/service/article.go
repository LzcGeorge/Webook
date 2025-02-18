package service

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/article"
	"context"
)

type ArticleService interface {
	Edit(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	PublishWithTwoRepo(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	// 一个 Service 操作一个 Repo：读者写者共用一个库
	repo article.ArticleRepository

	// 一个 Service 操作两个 Repo：读者库，写者库
	authorRepo article.ArticleAuthorRepository
	readerRepo article.ArticleReaderRepository
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func NewArticleServiceWithTwoRepo(authorRepo article.ArticleAuthorRepository, readerRepo article.ArticleReaderRepository) ArticleService {
	return &articleService{
		authorRepo: authorRepo,
		readerRepo: readerRepo,
	}
}

// Edit 编辑文章： 返回文章 id
func (a *articleService) Edit(ctx context.Context, article domain.Article) (int64, error) {
	if article.Id > 0 {
		return a.repo.Update(ctx, article)
	}
	return a.repo.Create(ctx, article)
}

// Publish 发布文章
func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {

	return 0, nil
}

// PublishWithTwoRepo 采用读者库和写者库
func (a *articleService) PublishWithTwoRepo(ctx context.Context, article domain.Article) (int64, error) {
	// 写者库发表文章
	var id = article.Id
	var err error
	if article.Id > 0 {
		id, err = a.authorRepo.Update(ctx, article)
	} else {
		id, err = a.authorRepo.Create(ctx, article)
	}
	if err != nil {
		return 0, err
	}
	// 确保写者库和读者库的 id 一致
	article.Id = id
	return a.readerRepo.Save(ctx, article)
}
