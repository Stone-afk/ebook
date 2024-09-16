package ioc

import (
	"context"
	rankingv1 "ebook/cmd/api/proto/gen/ranking/v1"
	"ebook/cmd/cronjob/domain"
	"ebook/cmd/cronjob/executor"
	"ebook/cmd/cronjob/executor/local"
	"ebook/cmd/cronjob/scheduler"
	"ebook/cmd/cronjob/service"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/schedulerx"
	"time"
)

func InitScheduler(l logger.Logger,
	svc service.CronJobService,
	execs ...executor.Executor) schedulerx.Scheduler {
	res := scheduler.NewScheduler(svc, l)
	res.RegisterExecutors(execs...)
	return res
}

func InitExecutors(svc rankingv1.RankingServiceClient) []executor.Executor {
	res := local.NewExecutor()
	res.RegisterFunc("ranking", func(ctx context.Context, j domain.CronJob) error {
		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()
		_, err := svc.RankTopN(ctx, &rankingv1.RankTopNRequest{})
		return err
	})
	return []executor.Executor{res}
}
