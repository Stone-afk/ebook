package repository

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository/dao/job"
	"time"
)

var ErrNoMoreJob = job.ErrNoMoreJob

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/cron_job.go.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/mocks/cron_job.mock.go
type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
}

type PreemptCronJobRepository struct {
	dao job.JobDAO
}

func NewPreemptCronJobRepository(dao job.JobDAO) CronJobRepository {
	return &PreemptCronJobRepository{dao: dao}
}

func (p *PreemptCronJobRepository) Preempt(ctx context.Context) (domain.CronJob, error) {
	panic("")
}

func (p *PreemptCronJobRepository) toEntity(j domain.CronJob) job.Job {
	return job.Job{
		Id:       j.Id,
		Name:     j.Name,
		Cfg:      j.Cfg,
		NextTime: j.NextTime.UnixMilli(),
	}
}

func (p *PreemptCronJobRepository) toDomain(j job.Job) domain.CronJob {
	return domain.CronJob{
		Id:       j.Id,
		Name:     j.Name,
		Cfg:      j.Cfg,
		NextTime: time.UnixMilli(j.NextTime),
	}
}
