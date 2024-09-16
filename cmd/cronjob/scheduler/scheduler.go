package scheduler

import (
	"context"
	"ebook/cmd/cronjob/domain"
	"ebook/cmd/cronjob/executor"
	"ebook/cmd/cronjob/service"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/schedulerx"
	"golang.org/x/sync/semaphore"
	"time"
)

// CronJob 使用别名来做一个解耦
// 后续万一我们要加字段，就很方便扩展
type CronJob = domain.CronJob

var _ schedulerx.Scheduler = (*Scheduler)(nil)

// Scheduler 调度器
type Scheduler struct {
	l         logger.Logger
	interval  time.Duration
	dbTimeout time.Duration
	limiter   *semaphore.Weighted
	svc       service.CronJobService
	execs     map[string]executor.Executor
}

func NewScheduler(svc service.CronJobService, l logger.Logger) *Scheduler {
	return &Scheduler{
		svc:       svc,
		l:         l,
		interval:  time.Second,
		dbTimeout: time.Second,
		execs:     make(map[string]executor.Executor),
		limiter:   semaphore.NewWeighted(100),
	}
}

func (s *Scheduler) RegisterExecutor(exec executor.Executor) {
	s.execs[exec.Name()] = exec
}

func (s *Scheduler) RegisterExecutors(execs ...executor.Executor) {
	for _, exec := range execs {
		s.execs[exec.Name()] = exec
	}
}

func (s *Scheduler) RegisterJob(ctx context.Context, j CronJob) error {
	return s.svc.AddJob(ctx, j)
}

// Start 开始调度。当被取消，或者超时的时候，就会结束调度
func (s *Scheduler) Start(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			// 退出调度循环
			return ctx.Err()
		}
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			// 正常来说，只有 ctx 超时或者取消才会进来这里
			return err
		}
		// 一次调度的数据库查询时间
		dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		j, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			// 你不能 return
			// 你要继续下一轮
			s.l.Error("抢占任务失败", logger.Error(err))
			// 没有抢占到，进入下一个循环
			// 这里可以考虑睡眠一段时间
			// 你也可以进一步细分不同的错误，如果是可以容忍的错误，
			// 就继续，不然就直接 return
			time.Sleep(s.interval)
			continue
		}
		exec, ok := s.execs[j.Executor]
		if !ok {
			// DEBUG 的时候最好中断
			// 线上就继续
			s.l.Error("未找到对应的执行器, 请检查是否是框架支持的Executor方式",
				logger.String("executor", j.Executor))
			// 不支持的执行方式。
			// 比如说，这里要求的runner是调用 gRPC，我们就不支持
			err = j.CancelFunc()
			s.l.Error("释放任务失败",
				logger.Error(err),
				logger.Int64("jid", j.Id))
			continue
		}
		// 接下来就是执行
		// 怎么执行？
		go func() {
			defer func() {
				s.limiter.Release(1)
				er := j.CancelFunc()
				if er != nil {
					s.l.Error("释放任务失败",
						logger.Error(er),
						logger.Int64("jid", j.Id))
				}
			}()
			// 异步执行，不要阻塞主调度循环
			// 执行完毕之后
			// 这边要考虑超时控制，任务的超时控制
			er := exec.Exec(ctx, j)
			if er != nil {
				// 也可以考虑在这里重试
				s.l.Error("任务执行失败", logger.Error(er))
			}
			// 要不要考虑下一次调度？
			ctx, cancel := context.WithTimeout(context.Background(), s.dbTimeout)
			defer cancel()
			er = s.svc.ResetNextTime(ctx, j)
			if er != nil {
				s.l.Error("设置下一次执行时间失败", logger.Error(er))
			}
		}()
	}
}
