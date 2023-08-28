package aliyun

import (
	"context"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type Service struct {
	client   *dysmsapi.Client
	signName string
}

func NewService(c *dysmsapi.Client, signName string) *Service {
	return &Service{
		client:   c,
		signName: signName,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	panic("")
}
