package job

import (
	"context"
	"ebook/cmd/internal/service"
	"ebook/cmd/pkg/logger"
	rlock "github.com/gotomicro/redis-lock"
	"time"
)

var _ Job = (*RankingJob)(nil)

type RankingJob struct {
	svc service.RankingService
	// 一次运行的超时时间
	timeout    time.Duration
	lockClient *rlock.Client
	l          logger.Logger
	key        string
}

func NewRankingJob(svc service.RankingService,
	lockClient *rlock.Client,
	l logger.Logger,
	timeout time.Duration) *RankingJob {
	return &RankingJob{
		lockClient: lockClient,
		svc:        svc,
		timeout:    timeout,
		key:        "job:ranking",
		l:          l,
	}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	panic("")
}

func (r *RankingJob) run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.RankTopN(ctx)
}
