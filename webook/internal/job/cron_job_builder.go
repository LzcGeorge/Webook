package job

import (
	"Webook/webook/pkg/logger"
	"github.com/robfig/cron/v3"
	"time"
)

type CronJonBuilder struct {
	logger logger.Logger
}

func NewCronJobBuilder(logger logger.Logger) *CronJonBuilder {
	return &CronJonBuilder{
		logger: logger,
	}
}

type cronJobAdapterFunc func()

func (c cronJobAdapterFunc) Run() {
	c()
}

func (c *CronJonBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return cronJobAdapterFunc(func() {
		start := time.Now()
		c.logger.Debug("start job", logger.String("name", name))
		err := job.Run()
		if err != nil {
			c.logger.Error("job failed",
				logger.String("name", name),
				logger.Error(err),
			)
		}
		duration := time.Since(start)
		c.logger.Debug("finish job",
			logger.String("name", name),
			logger.String("duration", duration.String()),
		)
	})
}
