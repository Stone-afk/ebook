package ioc

import (
	"ebook/cmd/internal/job"
	"ebook/cmd/internal/service"
	"ebook/cmd/pkg/logger"
	rlock "github.com/gotomicro/redis-lock"
	"time"
)

func InitRankingJob(svc service.RankingService,
	client *rlock.Client,
	l logger.Logger) *job.RankingJob {
	return job.NewRankingJob(svc, client, l, time.Second*30)
}
