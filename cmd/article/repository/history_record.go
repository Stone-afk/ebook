package repository

import (
	"context"
	"ebook/cmd/article/domain"
	"ebook/cmd/pkg/logger"
)

type historyRecordRepository struct {
	l logger.Logger
}

func NewHistoryRecordRepository(l logger.Logger) HistoryRecordRepository {
	return &historyRecordRepository{
		l: l,
	}
}

func (repo *historyRecordRepository) AddRecord(ctx context.Context, r domain.HistoryRecord) error {
	//TODO implement me
	panic("implement me")
}
