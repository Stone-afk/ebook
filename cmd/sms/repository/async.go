package repository

import (
	"context"
	"ebook/cmd/internal/repository/dao/async_sms"
	"ebook/cmd/sms/domain"
	"ebook/cmd/sms/repository/dao"
	"github.com/ecodeclub/ekit/sqlx"
)

var ErrWaitingSMSNotFound = dao.ErrWaitingSMSNotFound

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
