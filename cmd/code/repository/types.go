package repository

import (
	"context"
	"ebook/cmd/code/repository/cache"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/code/repository/types.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/code/repository/mocks/code.mock.go

type CodeRepository interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}
