package domain

import "time"

type CronJob struct {
	Id int64
	// Job 的名称，必须唯一
	Name     string
	Cfg      string
	NextTime time.Time
}
