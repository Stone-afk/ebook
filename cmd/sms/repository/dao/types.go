package dao

import "context"

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/sms/repository/dao/types.go -package=daomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/sms/repository/dao/mocks/async_sms.mock.go
type AsyncSmsDAO interface {
	Insert(ctx context.Context, s AsyncSms) error
	GetWaitingSMS(ctx context.Context) (AsyncSms, error)
	MarkSuccess(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64) error
}
