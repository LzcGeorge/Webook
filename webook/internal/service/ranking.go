package service

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/repository"
	"context"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"math"
	"time"
)

type RankingService interface {
	SetTop100(ctx context.Context) error
	GetTop100(ctx context.Context) ([]domain.Article, error)
	GetFromCache(ctx context.Context) ([]domain.Article, error)
}

type BatchRankingService struct {
	artSvc   ArticleService
	interSvc InteractiveService

	n         int
	batchSize int
	scoreFunc func(likeCnt int64, utime time.Time) float64

	repo repository.RankingRepository
}

func NewRankingService(artSvc ArticleService, interSvc InteractiveService, repo repository.RankingRepository) RankingService {
	return &BatchRankingService{
		artSvc:    artSvc,
		interSvc:  interSvc,
		n:         100,
		batchSize: 100,
		scoreFunc: func(likeCnt int64, utime time.Time) float64 {
			dur := time.Since(utime).Seconds()
			return float64(likeCnt-1) / math.Pow(dur+2, 1.5)
		},
		repo: repo,
	}
}
func (svc *BatchRankingService) SetTop100(ctx context.Context) error {
	articles, err := svc.GetTop100(ctx)
	if err != nil {
		return err
	}

	// 存到缓存中
	return svc.repo.ReplaceTop100(ctx, articles)
}

func (svc *BatchRankingService) GetFromCache(ctx context.Context) ([]domain.Article, error) {
	return svc.repo.GetTop100(ctx)
}

// GetTop100 获取排行榜前 100 的文章
func (svc *BatchRankingService) GetTop100(ctx context.Context) ([]domain.Article, error) {
	offset := 0
	now := time.Now()

	type Score struct {
		score float64
		art   domain.Article
	}

	// 用 小根堆 来维护 Score 前 100 的文章。
	pq := queue.NewPriorityQueue[Score](svc.n, func(a, b Score) int {
		if a.score > b.score {
			return 1
		} else if a.score == b.score {
			return 0
		} else {
			return -1
		}
	})

	// 批次取数据
	for {
		arts, err := svc.artSvc.PublicList(ctx, now, offset, svc.batchSize)
		if err != nil {
			return nil, err
		}

		bizIds := slice.Map(arts, func(idx int, art domain.Article) int64 {
			return art.Id
		})

		// 获取文章点赞数
		interMap, err := svc.interSvc.GetInterMapByBizIds(ctx, "article", bizIds)
		if err != nil {
			return nil, err
		}
		for _, art := range arts {
			inter := interMap[art.Id]
			score := svc.scoreFunc(inter.LikeCnt, art.Utime)
			v := Score{
				score: score,
				art:   art,
			}

			// 如果队列已满，则比较两者 score，保留较大的
			if err := pq.Enqueue(v); err == queue.ErrOutOfCapacity {
				minV, _ := pq.Dequeue()
				if minV.score < v.score {
					pq.Enqueue(v)
				} else {
					pq.Enqueue(minV)
				}
			}
		}

		// 有可能 len(arts) < svc.batchSize，说明已经取完了
		offset = offset + len(arts)
		if len(arts) < svc.batchSize {
			break
		}
	}

	// 封装结果
	res := make([]domain.Article, pq.Len())
	for i := pq.Len() - 1; i >= 0; i-- {
		v, _ := pq.Dequeue()
		res[i] = v.art
	}
	return res, nil
}
