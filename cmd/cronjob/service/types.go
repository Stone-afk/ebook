package service

import (
	"context"
	"ebook/cmd/cronjob/domain"
	"ebook/cmd/cronjob/repository"
)

var ErrNoMoreJob = repository.ErrNoMoreJob

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/cronjob/service/types.go -package=svcmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/cronjob/service/mocks/cron_job.mock.go
type CronJobService interface {
	AddJob(ctx context.Context, j domain.CronJob) error
	Preempt(ctx context.Context) (domain.CronJob, error)
	ResetNextTime(ctx context.Context, j domain.CronJob) error
	// 返回一个释放的方法，然后调用者取调
	// PreemptV1(ctx context.Context) (domain.Job, func() error,  error)
	// Release
	//Release(ctx context.Context, id int64) error
}
