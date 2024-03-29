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
	panic("")
}
