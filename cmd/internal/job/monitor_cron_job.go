package job

import (
	"ebook/cmd/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
)

// CronJobBuilder 根据需要加各种监控
// 这种 Builder 写法是为了避开 prometheus 重复注册的问题
// 也可以用来组装不同的装饰器，比较灵活
type CronJobBuilder struct {
	vector *prometheus.SummaryVec
	l      logger.Logger
}

func (m *CronJobBuilder) Build(job Job) cron.Job {
	panic("")
}
