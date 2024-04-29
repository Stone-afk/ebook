package service

import (
	"context"
	"ebook/cmd/cronjob/domain"
	"ebook/cmd/cronjob/repository"
	"ebook/cmd/pkg/logger"
	"time"
)

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

func (s *cronJobService) AddJob(ctx context.Context, j domain.CronJob) error {
	j.NextTime = j.Next(time.Now())
	return s.repo.AddJob(ctx, j)
}

func (s *cronJobService) Preempt(ctx context.Context) (domain.CronJob, error) {
	j, err := s.repo.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	// 你的续约呢？
	//ch := make(chan struct{})
	//go func() {
	//	ticker := time.NewTicker(p.refreshInterval)
	//	for {
	//		select {
	//		case <-ticker.C:
	//			// 在这里续约
	//			p.refresh(j.Id)
	//		case <-ch:
	//			// 结束
	//			return
	//		}
	//	}
	//}()
	ticker := time.NewTicker(s.refreshInterval)
	go func() {
		for range ticker.C {
			s.refresh(j.Id)
		}
	}()

	// 你抢占之后，你一直抢占着吗？
	// 你要考虑一个释放的问题
	j.CancelFunc = func() error {
		//close(ch)
		// 自己在这里释放掉
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		return s.repo.Release(ctx, j.Id)
	}
	return j, nil
}

func (s *cronJobService) ResetNextTime(ctx context.Context, j domain.CronJob) error {
	next := j.Next(time.Now())
	if next.IsZero() {
		// 没有下一次
		return s.repo.Stop(ctx, j.Id)
	}
	return s.repo.UpdateNextTime(ctx, j.Id, next)
}

func (s *cronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 续约怎么个续法？
	// 更新一下更新时间就可以
	// 比如说我们的续约失败逻辑就是：处于 running 状态，
	// 但是更新时间在三分钟以前, 代表发生该情况该通常是该节点挂掉了
	err := s.repo.UpdateUtime(ctx, id)
	if err != nil {
		// 可以考虑立刻重试
		s.l.Error("续约失败",
			logger.Error(err),
			logger.Int64("jid", id))
	}
}
