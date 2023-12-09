package service

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository"
	"ebook/cmd/pkg/logger"
	"time"
)

var ErrNoMoreJob = repository.ErrNoMoreJob

type CronJobService interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
}

type cronJobService struct {
	repo            repository.CronJobRepository
	l               logger.Logger
	refreshInterval time.Duration
}

func NewCronJobService(
	repo repository.CronJobRepository,
	l logger.Logger) CronJobService {
	return &cronJobService{
		repo:            repo,
		l:               l,
		refreshInterval: time.Second * 10,
	}
}

func (s *cronJobService) Preempt(ctx context.Context) (domain.CronJob, error) {
	j, err := s.repo.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	return j, nil
}
