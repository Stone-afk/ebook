package repository

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository/dao/async_sms"
)

var ErrWaitingSMSNotFound = async_sms.ErrWaitingSMSNotFound

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/async_sms.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/mocks/async_sms.mock.go
type AsyncSmsRepository interface {
	// Add 添加一个异步 SMS 记录。
	// 你叫做 Create 或者 Insert 也可以
	Add(ctx context.Context, s domain.AsyncSms) error
	PreemptWaitingSMS(ctx context.Context) (domain.AsyncSms, error)
	ReportScheduleResult(ctx context.Context, id int64, success bool) error
}

type asyncSmsRepository struct {
	dao async_sms.AsyncSmsDAO
}

func NewAsyncSMSRepository(dao async_sms.AsyncSmsDAO) AsyncSmsRepository {
	return &asyncSmsRepository{
		dao: dao,
	}
}

func (repo *asyncSmsRepository) Add(ctx context.Context, s domain.AsyncSms) error {
	//TODO implement me
	panic("implement me")
}

func (repo *asyncSmsRepository) PreemptWaitingSMS(ctx context.Context) (domain.AsyncSms, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *asyncSmsRepository) ReportScheduleResult(ctx context.Context, id int64, success bool) error {
	//TODO implement me
	panic("implement me")
}
