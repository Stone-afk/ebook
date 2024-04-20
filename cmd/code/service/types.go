package service

import (
	"context"
	"ebook/cmd/code/repository"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

const codeTplId = "1877556"

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/code/service/types.go -package=svcmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/code/service/mocks/code.mock.go
type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}
