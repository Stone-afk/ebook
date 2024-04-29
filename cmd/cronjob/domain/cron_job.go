package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type CronJob struct {
	Id int64
	// Job 的名称，必须唯一
	Name       string
	Executor   string
	Expression string
	// 通用的任务的抽象，我们也不知道任务的具体细节，所以就搞一个 Cfg
	// 具体任务设置具体的值
	Cfg      string
	NextTime time.Time
	// 放弃抢占状态
	CancelFunc func() error
}

var parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom |
	cron.Month | cron.Dow | cron.Descriptor)

func (j CronJob) Next(t time.Time) time.Time {
	// 怎么算？要根据 cron 表达式来算
	// 可以做成包变量，因为基本不可能变
	// 这个地方 Expression 必须不能出错，这需要用户在注册的时候确保
	s, _ := parser.Parse(j.Expression)
	return s.Next(t)
}
