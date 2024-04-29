package repository

import (
	"context"
	"ebook/cmd/cronjob/domain"
	"ebook/cmd/cronjob/repository/dao"
	"time"
)

var ErrNoMoreJob = dao.ErrNoMoreJob

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/cronjob/repository/types.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/cronjob/repository/mocks/cron_job.mock.go
type CronJobRepository interface {
	AddJob(ctx context.Context, j domain.CronJob) error
	Preempt(ctx context.Context) (domain.CronJob, error)
	Release(ctx context.Context, id int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	Stop(ctx context.Context, id int64) error
}
