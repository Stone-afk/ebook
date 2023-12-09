package job

import (
	"context"
	"gorm.io/gorm"
	"time"
)

var ErrNoMoreJob = gorm.ErrRecordNotFound

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
}

type GORMJobDAO struct {
	db *gorm.DB
}

func NewGORMJobDAO(db *gorm.DB) JobDAO {
	return &GORMJobDAO{db: db}
}

func (dao *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	db := dao.db.WithContext(ctx)
	for {
		// 每一个循环都重新计算 time.Now，因为之前可能已经花了一些时间了
		now := time.Now().UnixMilli()
		var j Job
		err := db.Where("next_time <= ? AND status = ?",
			now, jobStatusWaiting).First(&j).Error
		if err != nil {
			// 数据库有问题
			return Job{}, err
		}
		// 然后要开始抢占
		// 这里利用 utime 来执行 CAS 操作
		// 其它一些公司可能会有一些 version 之类的字段
		res := db.Model(&Job{}).
			Where("id = ? AND version=?", j.Id, j.Version).
			Updates(map[string]any{
				"utime":   now,
				"version": j.Version + 1,
				"status":  jobStatusRunning,
			})
		if res.Error != nil {
			// 数据库错误
			return Job{}, err
		}
		// 抢占成功
		if res.RowsAffected == 1 {
			return j, nil
		}
		// 没有抢占到，也就是同一时刻被人抢走了，那么就下一个循环
	}
}
