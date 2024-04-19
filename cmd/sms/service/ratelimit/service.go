package ratelimit

import (
	"context"
	"ebook/cmd/pkg/ratelimit"
	"ebook/cmd/sms/service"
	"fmt"
)

var errLimited = fmt.Errorf("触发了限流")

type Service struct {
	key     string // 服务商地址
	svc     service.Service
	limiter ratelimit.Limiter
}

func NewService(key string, svc service.Service, limiter ratelimit.Limiter) service.Service {
	return &Service{
		key:     key,
		svc:     svc,
		limiter: limiter,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	key := fmt.Sprintf("sms:%s", s.key)
	// 可以限流：保守策略，你的下游很坑的时候，
	// 可以不限：你的下游很强，业务可用性要求很高，尽量容错策略
	limited, err := s.limiter.Limit(ctx, key)
	if err != nil {
		// 系统错误,
		// 保守策略，返回错误
		// 非保守策略，跳过错误，直接发
		// 包一下这个错误
		return fmt.Errorf("短信服务限流出现问题，%w", err)
	}
	if limited {
		return errLimited
	}

	// 你这里加一些代码，新特性
	// err = s.svc.Send(ctx, tpl, args, numbers...)
	// 你在这里也可以加一些代码，新特性
	return s.svc.Send(ctx, tpl, args, numbers...)
}
