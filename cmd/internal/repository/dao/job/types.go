package job

import (
	"context"
	"gorm.io/gorm"
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
	panic("")
}
