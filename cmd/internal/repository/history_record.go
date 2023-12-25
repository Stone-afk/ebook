package repository

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/pkg/logger"
)

// HistoryRecordRepository 也就是一个增删改查的事情
type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, r domain.HistoryRecord) error
}

type historyRecordRepository struct {
	l logger.Logger
}

func (repo *historyRecordRepository) AddRecord(ctx context.Context, r domain.HistoryRecord) error {
	//TODO implement me
	panic("implement me")
}
