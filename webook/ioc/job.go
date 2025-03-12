package ioc

import (
	"Webook/webook/internal/job"
	"Webook/webook/internal/service"
	"Webook/webook/pkg/logger"
	"github.com/robfig/cron/v3"
	"time"
)

func InitRankingJob(svc service.RankingService) *job.RankingJob {
	return job.NewRankingJob(svc, time.Second*30)
}
func InitJobs(logger logger.Logger, rankJob *job.RankingJob) *cron.Cron {
	cronJobBuilder := job.NewCronJobBuilder(logger)
	cornn := cron.New(cron.WithSeconds())
	_, err := cornn.AddJob("@every 1m", cronJobBuilder.Build(rankJob))
	if err != nil {
		panic(err)
	}
	return cornn
}
