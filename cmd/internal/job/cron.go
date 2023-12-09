package job

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/service"
	"ebook/cmd/pkg/logger"
)

type Executor interface {
	Name() string
	// Exec ctx 是整个任务调度的上下文
	// 当从 ctx.Done 有信号的时候，就需要考虑结束执行
	// 具体实现来控制
	Exec(ctx context.Context, j domain.CronJob) error
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, j domain.CronJob) error
}

// Scheduler 调度器
type Scheduler struct {
	svc service.CronJobService
	l   logger.Logger
}

func NewScheduler(svc service.CronJobService, l logger.Logger) *Scheduler {
	return &Scheduler{
		svc: svc,
		l:   l,
	}
}

// Start 开始调度。当被取消，或者超时的时候，就会结束调度
func (s *Scheduler) Start(ctx context.Context) error {
	panic("")
}

// CronJob 使用别名来做一个解耦
// 后续万一我们要加字段，就很方便扩展
type CronJob = domain.CronJob
