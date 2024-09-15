package local

import (
	"context"
	"ebook/cmd/cronjob/domain"
	"ebook/cmd/cronjob/executor"
	"fmt"
)

var _ executor.Executor = (*Executor)(nil)

type Executor struct {
	funcMap map[string]func(ctx context.Context, j domain.CronJob) error
}

func NewExecutor() *Executor {
	return &Executor{
		funcMap: make(map[string]func(ctx context.Context, j domain.CronJob) error),
	}
}

func (e *Executor) Name() string {
	return "local"
}

func (e *Executor) RegisterFunc(name string, fn func(ctx context.Context, j domain.CronJob) error) {
	e.funcMap[name] = fn
}

func (e *Executor) Exec(ctx context.Context, j domain.CronJob) error {
	fn, ok := e.funcMap[j.Name]
	if !ok {
		return fmt.Errorf("未知任务，你是否注册了？ %s", j.Name)
	}
	return fn(ctx, j)
}
