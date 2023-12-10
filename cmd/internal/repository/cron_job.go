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
	Release(ctx context.Context, id int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	Stop(ctx context.Context, id int64) error
}

type PreemptCronJobRepository struct {
	dao job.JobDAO
}

func NewPreemptCronJobRepository(dao job.JobDAO) CronJobRepository {
	return &PreemptCronJobRepository{dao: dao}
}

func (repo *PreemptCronJobRepository) UpdateUtime(ctx context.Context, id int64) error {
	return repo.dao.UpdateUtime(ctx, id)
}

func (repo *PreemptCronJobRepository) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return repo.dao.UpdateNextTime(ctx, id, next)
}

func (repo *PreemptCronJobRepository) Stop(ctx context.Context, id int64) error {
	return repo.dao.Stop(ctx, id)
}

func (repo *PreemptCronJobRepository) Release(ctx context.Context, id int64) error {
	return repo.dao.Release(ctx, id)
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
		Id:   j.Id,
		Name: j.Name,
		Cfg:  j.Cfg,
	}
}

func (repo *PreemptCronJobRepository) toDomain(j job.Job) domain.CronJob {
	return domain.CronJob{
		Id:   j.Id,
		Name: j.Name,
		Cfg:  j.Cfg,
	}
}
