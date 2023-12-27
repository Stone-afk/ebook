package ioc

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/job"
	"ebook/cmd/internal/service"
	"ebook/cmd/pkg/logger"
	rlock "github.com/gotomicro/redis-lock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"time"
)

func InitRankingJob(svc service.RankingService,
	client *rlock.Client,
	l logger.Logger) *job.RankingJob {
	return job.NewRankingJob(svc, client, l, time.Second*30)
}

func InitJobs(l logger.Logger, rankingJob *job.RankingJob) *cron.Cron {
	bd := job.NewCronJobBuilder(l, prometheus.SummaryOpts{
		Namespace: "stone",
		Subsystem: "ebook",
		Name:      "cron_job",
		Help:      "定时任务",
	})
	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("@every 1m", bd.Build(rankingJob))
	if err != nil {
		panic(err)
	}
	return expr
}

func InitScheduler(l logger.Logger,
	local *job.LocalFuncExecutor,
	svc service.CronJobService) *job.Scheduler {
	res := job.NewScheduler(svc, l)
	res.RegisterExecutor(local)
	return res
}

func InitLocalFuncExecutor(svc service.RankingService) *job.LocalFuncExecutor {
	res := job.NewLocalFuncExecutor()
	res.RegisterFunc("ranking", func(ctx context.Context, j domain.CronJob) error {
		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()
		return svc.RankTopN(ctx)
	})
	return res
}
