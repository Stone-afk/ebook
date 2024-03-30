package validator

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/migrator"
	"ebook/cmd/pkg/migrator/events"
	"gorm.io/gorm"
)

type CanalIncrValidator[T migrator.Entity] struct {
	baseValidator
}

func NewCanalIncrValidator[T migrator.Entity](
	base *gorm.DB,
	target *gorm.DB,
	direction string,
	l logger.Logger,
	producer events.Producer,
) *CanalIncrValidator[T] {
	return &CanalIncrValidator[T]{
		baseValidator: baseValidator{
			base:      base,
			target:    target,
			direction: direction,
			l:         l,
			producer:  producer,
		},
	}
}

// Validate 一次校验一条
func (v *CanalIncrValidator[T]) Validate(ctx context.Context, id int64) error {
	var base T
	err := v.base.WithContext(ctx).Where("id = ?", id).First(&base).Error
	switch err {
	case gorm.ErrRecordNotFound:
		var target T
		err1 := v.target.WithContext(ctx).Where("id = ?", id).First(&target).Error
		switch err1 {
		case gorm.ErrRecordNotFound:
			// 数据一致
			return nil
		case nil:
			v.notify(id, events.InconsistentEventTypeBaseMissing)
		default:
			return err
		}
	case nil:
		var target T
		err1 := v.target.WithContext(ctx).Where("id = ?", id).First(&target).Error
		switch err1 {
		case gorm.ErrRecordNotFound:
			v.notify(id, events.InconsistentEventTypeTargetMissing)
		case nil:
			if !base.CompareTo(target) {
				v.notify(id, events.InconsistentEventTypeNEQ)
			}
		default:
			return err
		}
	default:
		return err
	}
	return nil
}
