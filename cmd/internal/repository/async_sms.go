package repository

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository/dao/async_sms"
	"github.com/ecodeclub/ekit/sqlx"
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
	return repo.dao.Insert(ctx, async_sms.AsyncSms{
		RetryMax: s.RetryMax,
		Config: sqlx.JsonColumn[async_sms.SmsConfig]{
			Val: async_sms.SmsConfig{
				TplId:   s.TplId,
				Args:    s.Args,
				Numbers: s.Numbers,
			},
			Valid: true,
		},
	})
}

func (repo *asyncSmsRepository) PreemptWaitingSMS(ctx context.Context) (domain.AsyncSms, error) {
	as, err := repo.dao.GetWaitingSMS(ctx)
	if err != nil {
		return domain.AsyncSms{}, err
	}
	return domain.AsyncSms{
		Id:       as.Id,
		TplId:    as.Config.Val.TplId,
		Numbers:  as.Config.Val.Numbers,
		Args:     as.Config.Val.Args,
		RetryMax: as.RetryMax,
	}, nil
}

func (repo *asyncSmsRepository) ReportScheduleResult(ctx context.Context, id int64, success bool) error {
	if success {
		return repo.dao.MarkSuccess(ctx, id)
	}
	return repo.dao.MarkFailed(ctx, id)
}
