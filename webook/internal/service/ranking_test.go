package service

import (
	"Webook/webook/internal/domain"
	svcmocks "Webook/webook/internal/service/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestRankingTopN(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (ArticleService, InteractiveService)
		wantErr  error
		wantArts []domain.Article
	}{
		{
			name: "计算成功",
			mock: func(ctrl *gomock.Controller) (ArticleService, InteractiveService) {
				artSvc := svcmocks.NewMockArticleService(ctrl)
				intrSvc := svcmocks.NewMockInteractiveService(ctrl)
				artSvc.EXPECT().PublicList(gomock.Any(), gomock.Any(), 0, 2).Return([]domain.Article{
					{Id: 1, Ctime: now, Utime: now},
					{Id: 2, Ctime: now, Utime: now},
				}, nil)
				intrSvc.EXPECT().GetInterMapByBizIds(gomock.Any(), "article", []int64{1, 2}).
					Return(map[int64]domain.Interactive{
						1: {BizId: 1, LikeCnt: 1},
						2: {BizId: 2, LikeCnt: 2},
					}, nil)
				artSvc.EXPECT().PublicList(gomock.Any(), gomock.Any(), 2, 2).Return([]domain.Article{
					{Id: 3, Ctime: now, Utime: now},
					{Id: 4, Ctime: now, Utime: now},
				}, nil)
				intrSvc.EXPECT().GetInterMapByBizIds(gomock.Any(), "article", []int64{3, 4}).
					Return(map[int64]domain.Interactive{
						3: {BizId: 3, LikeCnt: 3},
						4: {BizId: 4, LikeCnt: 4},
					}, nil)

				artSvc.EXPECT().PublicList(gomock.Any(), gomock.Any(), 4, 2).Return([]domain.Article{}, nil)
				intrSvc.EXPECT().GetInterMapByBizIds(gomock.Any(), "article", []int64{}).Return(map[int64]domain.Interactive{}, nil)
				return artSvc, intrSvc
			},
			wantArts: []domain.Article{
				{Id: 4, Ctime: now, Utime: now},
				{Id: 3, Ctime: now, Utime: now},
				{Id: 2, Ctime: now, Utime: now},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			artSvc, interSvc := tc.mock(ctrl)
			svc := &BatchRankingService{
				artSvc:    artSvc,
				interSvc:  interSvc,
				n:         3,
				batchSize: 2,
				scoreFunc: func(likeCnt int64, utime time.Time) float64 {
					return float64(likeCnt + 2)
				},
			}

			arts, err := svc.GetTop100(context.Background())
			t.Log("arts: ", arts)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantArts, arts)
		})
	}
}
