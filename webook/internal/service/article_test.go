package service

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository/article"
	repomocks "Webook/webook/internal/repository/article/mocks"
	"Webook/webook/pkg/logger"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestArticleService_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository)

		// service 的参数
		article domain.Article

		// service 的期待返回值
		wantId  int64
		wantErr error
	}{
		{
			name: "创建文章，并发布成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)

				// 模拟写者库创建文章的过程，要求入参为 Id 为 0 的 Article
				// 返回 1，nil
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Id:      0, // 默认是 0，不写这行也行
					Title:   "create article and publish",
					Content: "this is content",
					Author: domain.Author{
						Id: 666,
					},
				}).Return(int64(1), nil)

				// 模拟查询读者库，返回 空文章，nil ，表示文章不存在
				readerRepo.EXPECT().FindById(gomock.Any(), int64(1)).Return(domain.Article{}, nil)

				// 模拟读者库创建文章的过程，要求入参为 Id 为 1 的 Article
				// 返回 1，nil
				readerRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Id:      1, // 用写者库的 id
					Title:   "create article and publish",
					Content: "this is content",
					Author: domain.Author{
						Id: 666,
					},
				}).Return(int64(1), nil)

				return authorRepo, readerRepo
			},
			article: domain.Article{
				Title:   "create article and publish",
				Content: "this is content",
				Author: domain.Author{
					Id: 666,
				},
			},
			wantId:  1,
			wantErr: nil,
		},

		{
			name: "修改文章，并发布成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)

				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      12, // Id > 0 表示是已经创建的文章
					Title:   "edit article and publish",
					Content: "fix: this is content",
					Author: domain.Author{
						Id: 666,
					},
				}).Return(int64(12), nil)

				// 模拟查询读者库，返回 文章 id，nil ，表示文章存在
				readerRepo.EXPECT().FindById(gomock.Any(), int64(12)).Return(domain.Article{
					Id: 12,
				}, nil)

				readerRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      12, // 用写者库的 id
					Title:   "edit article and publish",
					Content: "fix: this is content",
					Author: domain.Author{
						Id: 666,
					},
				}).Return(int64(12), nil)

				return authorRepo, readerRepo
			},
			article: domain.Article{
				Id:      12,
				Title:   "edit article and publish",
				Content: "fix: this is content",
				Author: domain.Author{
					Id: 666,
				},
			},
			wantId:  12,
			wantErr: nil,
		},
		{
			name: "写者库修改成功，读者库保存失败, 重试保存成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)

				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      12, // Id > 0 表示是已经创建的文章
					Title:   "edit article and publish",
					Content: "fix: this is content",
					Author: domain.Author{
						Id: 666,
					},
				}).Return(int64(12), nil)

				// 模拟查询读者库，返回 文章 id，nil ，表示文章存在
				readerRepo.EXPECT().FindById(gomock.Any(), int64(12)).Return(domain.Article{
					Id: 12,
				}, nil)

				readerRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      12, // 用写者库的 id
					Title:   "edit article and publish",
					Content: "fix: this is content",
					Author: domain.Author{
						Id: 666,
					},
				}).Return(int64(12), errors.New("update failed"))

				// 重试
				readerRepo.EXPECT().FindById(gomock.Any(), int64(12)).Return(domain.Article{
					Id: 12,
				}, nil)
				readerRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      12, // 用写者库的 id
					Title:   "edit article and publish",
					Content: "fix: this is content",
					Author: domain.Author{
						Id: 666,
					},
				}).Return(int64(12), nil)
				return authorRepo, readerRepo
			},
			article: domain.Article{
				Id:      12,
				Title:   "edit article and publish",
				Content: "fix: this is content",
				Author: domain.Author{
					Id: 666,
				},
			},
			wantId:  12,
			wantErr: nil,
		},
		{
			name: "写者库修改成功，读者库保存失败, 重试三次仍然失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)

				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      12, // Id > 0 表示是已经创建的文章
					Title:   "edit article and publish",
					Content: "fix: this is content",
					Author: domain.Author{
						Id: 666,
					},
				}).Return(int64(12), nil)

				for i := 0; i < 3; i++ {
					// 模拟查询读者库，返回 文章 id，nil ，表示文章存在
					readerRepo.EXPECT().FindById(gomock.Any(), int64(12)).Return(domain.Article{
						Id: 12,
					}, nil)

					readerRepo.EXPECT().Update(gomock.Any(), domain.Article{
						Id:      12, // 用写者库的 id
						Title:   "edit article and publish",
						Content: "fix: this is content",
						Author: domain.Author{
							Id: 666,
						},
					}).Return(int64(12), errors.New("update failed"))
				}

				return authorRepo, readerRepo
			},
			article: domain.Article{
				Id:      12,
				Title:   "edit article and publish",
				Content: "fix: this is content",
				Author: domain.Author{
					Id: 666,
				},
			},
			wantId:  12,
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			l, err := zap.NewDevelopment()
			assert.Nil(t, err)
			logger := logger.NewZapLogger(l)
			authorRepo, readerRepo := tc.mock(ctrl)
			svc := NewArticleServiceWithTwoRepo(authorRepo, readerRepo, logger)
			resId, err := svc.PublishWithTwoRepo(context.Background(), tc.article)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, resId)

		})
	}
}
