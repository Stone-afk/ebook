package executor

import (
	"context"
	"ebook/cmd/cronjob/domain"
)

type Executor interface {
	Name() string
	// Exec ctx 是整个任务调度的上下文
	// 当从 ctx.Done 有信号的时候，就需要考虑结束执行
	// 具体实现来控制
	Exec(ctx context.Context, j domain.CronJob) error
}
