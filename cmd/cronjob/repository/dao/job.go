package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

const (
	// 等待被调度，意思就是没有人正在调度
	jobStatusWaiting = iota
	// 已经被 goroutine 抢占了
	jobStatusRunning
	// 不再需要调度了，比如说被终止了，或者被删除了。
	// 我们这里没有严格区分这两种情况的必要性
	// 暂停调度
	jobStatusPaused
)

type GORMJobDAO struct {
	db *gorm.DB
}

func NewGORMJobDAO(db *gorm.DB) JobDAO {
	return &GORMJobDAO{db: db}
}

func (dao *GORMJobDAO) Insert(ctx context.Context, j Job) error {
	now := time.Now().UnixMilli()
	j.Ctime = now
	j.Utime = now
	return dao.db.WithContext(ctx).Create(&j).Error
}

func (dao *GORMJobDAO) Stop(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
		"status": jobStatusPaused,
		"utime":  time.Now().UnixMilli(),
	}).Error
}

func (dao *GORMJobDAO) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return dao.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
		"next_time": next.UnixMilli(),
	}).Error
}

func (dao *GORMJobDAO) UpdateUtime(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
		"utime": time.Now().UnixMilli(),
	}).Error
}

func (dao *GORMJobDAO) Release(ctx context.Context, id int64) error {
	// 这里有一个问题。你要不要检测 status 或者 version?
	// WHERE version = ?
	// TODO 要, 记得修改
	return dao.db.WithContext(ctx).Model(&Job{}).
		Where("id =?", id).Updates(map[string]any{
		"status": jobStatusWaiting,
		"utime":  time.Now().UnixMilli(),
	}).Error
}

func (dao *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	// 高并发情况下，大部分都是陪太子读书
	// 100 个 goroutine
	// 要转几次？ 所有 goroutine 执行的循环次数加在一起是
	// 1+2+3+4 +5 + ... + 99 + 100
	// 特定一个 goroutine，最差情况下，要循环一百次
	db := dao.db.WithContext(ctx)
	for {
		// 每一个循环都重新计算 time.Now，因为之前可能已经花了一些时间了
		now := time.Now().UnixMilli()
		var j Job
		// 分布式任务调度系统
		// 1. 一次拉一批，我一次性取出 100 条来，然后，我随机从某一条开始，向后开始抢占
		// 2. 我搞个随机偏移量，0-100 生成一个随机偏移量。兜底：第一轮没查到，偏移量回归到 0
		// 3. 我搞一个 id 取余分配，status = ? AND next_time <=? AND id%10 = ? 兜底：不加余数条件，取next_time 最老的
		err := db.Where("next_time <= ? AND status = ?",
			now, jobStatusWaiting).First(&j).Error
		if err != nil {
			// 数据库有问题
			return Job{}, err
		}
		// 然后要开始抢占
		// 两个 goroutine 都拿到 id =1 的数据
		// 能不能用 utime?
		// 乐观锁，CAS 操作，compare AND Swap
		// 有一个很常见的面试刷亮点：就是用乐观锁取代 FOR UPDATE
		// 面试套路（性能优化）：曾将用了 FOR UPDATE =>性能差，还会有死锁 => 我优化成了乐观锁
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
		if res.RowsAffected == 0 {
			// 没有抢占到，也就是同一时刻被人抢走了，那么就下一个循环
			continue
		}
		// 抢占成功
		return j, nil
	}
}
