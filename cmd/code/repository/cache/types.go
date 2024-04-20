package cache

import (
	"context"
	"errors"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("我也不知发生什么了，反正是跟 code 有关")
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/code/repository/cache/types.go -package=cachemocks -destination=/Users/stone/go_project/ebook/ebook/cmd/code/repository/cache/mocks/code.mock.go
type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}
