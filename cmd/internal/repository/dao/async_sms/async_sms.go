package async_sms

import (
	"context"
	"gorm.io/gorm"
)

var ErrWaitingSMSNotFound = gorm.ErrRecordNotFound

type GORMAsyncSmsDAO struct {
	db *gorm.DB
}

func NewGORMAsyncSmsDAO(db *gorm.DB) AsyncSmsDAO {
	return &GORMAsyncSmsDAO{
		db: db,
	}
}

func (dao *GORMAsyncSmsDAO) Insert(ctx context.Context, s AsyncSms) error {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMAsyncSmsDAO) GetWaitingSMS(ctx context.Context) (AsyncSms, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMAsyncSmsDAO) MarkSuccess(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMAsyncSmsDAO) MarkFailed(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}
