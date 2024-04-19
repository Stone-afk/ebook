package retryable

import (
	"context"
	"ebook/cmd/sms/service"
	"errors"
)

type Service struct {
	svc service.Service
	// 重试
	retryMax int
}

func NewService(svc service.Service, retryMax int) service.Service {
	return &Service{
		svc:      svc,
		retryMax: retryMax,
	}
}

func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, numbers...)
	cnt := 1
	if err != nil && cnt < s.retryMax {
		err = s.svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		cnt++
	}
	return errors.New("重试都失败了")
}
