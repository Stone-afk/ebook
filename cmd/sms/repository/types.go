package repository

import (
	"context"
	"ebook/cmd/sms/domain"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/sms/repository/types.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/sms/repository/mocks/async_sms.mock.go
type AsyncSmsRepository interface {
	// Add 添加一个异步 SMS 记录。
	// 你叫做 Create 或者 Insert 也可以
	Add(ctx context.Context, s domain.AsyncSms) error
	PreemptWaitingSMS(ctx context.Context) (domain.AsyncSms, error)
	ReportScheduleResult(ctx context.Context, id int64, success bool) error
}
