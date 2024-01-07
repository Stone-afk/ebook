package fixer

import (
	"context"
	"ebook/cmd/pkg/migrator"
	"ebook/cmd/pkg/migrator/events"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Fixer[T migrator.Entity] struct {
	base    *gorm.DB
	target  *gorm.DB
	columns []string
}

// Fix 最一了百了的写法
// 不管三七二十一，我TM直接覆盖
// 把 event 当成一个触发器，不依赖的 event 的具体内容（ID 必须不可变）
// 修复这里，也改成批量？？
func (f *Fixer[T]) Fix(ctx context.Context, evt events.InconsistentEvent) error {
	var t T
	err := f.base.WithContext(ctx).Where("id = ?", evt.ID).
		First(&t).Error
	switch err {
	case nil:
		// 找到了数据
		// base 有数据
		// 修复数据的时候，可以考虑增加 WHERE base.utime >= target.utime
		// utime 用不了，就看有没有version 之类的，或者能够判定数据新老的
		return f.target.Clauses(&clause.OnConflict{
			// 需要 Entity 告诉我们，修复哪些数据
			DoUpdates: clause.AssignmentColumns(f.columns),
		}).Create(&t).Error
	case gorm.ErrRecordNotFound:
		// base 没了
		return f.target.WithContext(ctx).
			Where("id=?", evt.ID).Delete(&t).Error
	default:
		return err
	}
}

// FixV1 最一了百了的写法
// 不管三七二十一，我TM直接覆盖
// 把 event 当成一个触发器，不依赖的 event 的具体内容（ID 必须不可变）
// 修复这里，也改成批量？？
func (f *Fixer[T]) FixV1(ctx context.Context, evt events.InconsistentEvent) error {
	panic("")
}

// FixV2 最一了百了的写法
// 不管三七二十一，我TM直接覆盖
// 把 event 当成一个触发器，不依赖的 event 的具体内容（ID 必须不可变）
// 修复这里，也改成批量？？
func (f *Fixer[T]) FixV2(ctx context.Context, evt events.InconsistentEvent) error {
	panic("")
}
