package fixer

import (
	"context"
	"ebook/cmd/pkg/migrator"
	"ebook/cmd/pkg/migrator/events"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OverrideFixer[T migrator.Entity] struct {
	// 因为本身其实这个不涉及什么领域对象，
	// 这里操作的不是 migrator 本身的领域对象
	base    *gorm.DB
	target  *gorm.DB
	columns []string
}

func NewOverrideFixer[T migrator.Entity](base *gorm.DB,
	target *gorm.DB) (*OverrideFixer[T], error) {
	// 在这里需要查询一下数据库中究竟有哪些列
	var t T
	rows, err := base.Model(&t).Limit(1).Rows()
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	return &OverrideFixer[T]{
		base:    base,
		target:  target,
		columns: columns,
	}, nil
}

func (f *OverrideFixer[T]) Fix(ctx context.Context, id int64) error {
	var src T
	err := f.base.WithContext(ctx).Where("id = ?", id).First(&src).Error
	switch err {
	case nil:
		return f.target.WithContext(ctx).Clauses(&clause.OnConflict{
			// 需要 Entity 告诉我们，修复哪些数据
			DoUpdates: clause.AssignmentColumns(f.columns),
		}).Create(&src).Error
	case gorm.ErrRecordNotFound:
		return f.target.WithContext(ctx).Delete("id = ?", id).Error
	default:
		return err
	}
}

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
			Where("id = ?", evt.ID).Delete(&t).Error
	default:
		return err
	}
}

// FixV1 最一了百了的写法
// 不管三七二十一，我TM直接覆盖
// 把 event 当成一个触发器，不依赖的 event 的具体内容（ID 必须不可变）
// 修复这里，也改成批量？？
func (f *Fixer[T]) FixV1(ctx context.Context, evt events.InconsistentEvent) error {
	switch evt.Type {
	case events.InconsistentEventTypeTargetMissing, events.InconsistentEventTypeNEQ:
		// 这边要插入
		var t T
		err := f.base.WithContext(ctx).
			Where("id = ?", evt.ID).First(&t).Error
		switch err {
		case gorm.ErrRecordNotFound:
			return f.target.WithContext(ctx).
				Where("id = ?", evt.ID).Delete(new(T)).Error
		case nil:
			return f.base.WithContext(ctx).Clauses(clause.OnConflict{
				// 这边要更新全部列
				DoUpdates: clause.AssignmentColumns(f.columns),
			}).Create(&t).Error
		default:
			return err
		}
	case events.InconsistentEventTypeBaseMissing:
		return f.target.WithContext(ctx).
			Where("id = ?", evt.ID).Delete(new(T)).Error
	default:
		return errors.New("未知的不一致类型")
	}
}

// FixV2 最一了百了的写法
// 不管三七二十一，我TM直接覆盖
// 把 event 当成一个触发器，不依赖的 event 的具体内容（ID 必须不可变）
// 修复这里，也改成批量？？
func (f *Fixer[T]) FixV2(ctx context.Context, evt events.InconsistentEvent) error {
	switch evt.Type {
	case events.InconsistentEventTypeTargetMissing:
		// 这边要插入
		var t T
		err := f.base.WithContext(ctx).
			Where("id = ?", evt.ID).First(&t).Error
		switch err {
		case gorm.ErrRecordNotFound:
			// base 也删除了这条数据
			return nil
		case nil:
			// 就在你插入的时候，双写的程序，也插入了，你就会冲突
			return f.target.Create(&t).Error
		default:
			return err
		}
	case events.InconsistentEventTypeNEQ:
		var t T
		err := f.base.WithContext(ctx).
			Where("id = ?", evt.ID).First(&t).Error
		switch err {
		case gorm.ErrRecordNotFound:
			// target 要删除
			return f.target.WithContext(ctx).
				Where("id = ?", evt.ID).Delete(&t).Error
		case nil:
			return f.target.Updates(&t).Error
		default:
			return err
		}
	case events.InconsistentEventTypeBaseMissing:
		return f.target.WithContext(ctx).
			Where("id = ?", evt.ID).Delete(new(T)).Error
	default:
		return errors.New("未知的不一致类型")
	}
}
