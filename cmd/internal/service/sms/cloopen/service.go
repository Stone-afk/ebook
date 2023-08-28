package cloopen

// 容联云短信的实现
// SDK文档:https://doc.yuntongxun.com/pe/5f029a06a80948a1006e7760

import (
	"context"
	"github.com/cloopen/go-sms-sdk/cloopen"
)

type Service struct {
	client *cloopen.SMS
	appId  string
}

func NewService(c *cloopen.SMS, addId string) *Service {
	return &Service{
		client: c,
		appId:  addId,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, data []string, numbers ...string) error {
	panic("")
}
