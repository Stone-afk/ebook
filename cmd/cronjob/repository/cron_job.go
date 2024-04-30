package repository

import (
	"context"
	"ebook/cmd/cronjob/domain"
	"ebook/cmd/cronjob/repository/dao"
	"time"
)

type PreemptCronJobRepository struct {
	dao dao.JobDAO
}

func NewPreemptCronJobRepository(dao dao.JobDAO) CronJobRepository {
	return &PreemptCronJobRepository{dao: dao}
}

func (repo *PreemptCronJobRepository) AddJob(ctx context.Context, j domain.CronJob) error {
	return repo.dao.Insert(ctx, repo.toEntity(j))
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

func (repo *PreemptCronJobRepository) toEntity(j domain.CronJob) dao.Job {
	return dao.Job{
		Id:         j.Id,
		Name:       j.Name,
		Expression: j.Expression,
		Cfg:        j.Cfg,
		Executor:   j.Executor,
		NextTime:   j.NextTime.UnixMilli(),
	}
}

func (repo *PreemptCronJobRepository) toDomain(j dao.Job) domain.CronJob {
	return domain.CronJob{
		Id:         j.Id,
		Name:       j.Name,
		Expression: j.Expression,
		Cfg:        j.Cfg,
		Executor:   j.Executor,
		NextTime:   time.UnixMilli(j.NextTime),
	}
}
