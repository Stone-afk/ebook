package wechat

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/service/oauth2"
	"net/http"
)

var _ oauth2.Service = &Service{}

type Service struct {
	appId     string
	appSecret string
	client    *http.Client
	//cmd       redis.Cmdable
}

type Result struct{}

func NewService(appId string, appSecret string, client *http.Client) oauth2.Service {
	return &Service{
		appId:     appId,
		appSecret: appSecret,
		client:    client,
	}
}

func (s *Service) AuthURL(ctx context.Context, state string) (string, error) {
	panic("")
}

func (s *Service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	panic("")
}
