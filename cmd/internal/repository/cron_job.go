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

func (repo *PreemptCronJobRepository) Preempt(ctx context.Context) (domain.CronJob, error) {
	j, err := repo.dao.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	return repo.toDomain(j), nil
}

func (repo *PreemptCronJobRepository) toEntity(j domain.CronJob) job.Job {
	return job.Job{
		Id:       j.Id,
		Name:     j.Name,
		Cfg:      j.Cfg,
		NextTime: j.NextTime.UnixMilli(),
	}
}

func (repo *PreemptCronJobRepository) toDomain(j job.Job) domain.CronJob {
	return domain.CronJob{
		Id:       j.Id,
		Name:     j.Name,
		Cfg:      j.Cfg,
		NextTime: time.UnixMilli(j.NextTime),
	}
}
