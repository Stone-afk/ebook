package service

import (
	"context"
	"ebook/cmd/im/domain"
	"net/http"
)

type RESTUserService struct {
	// 部署 IM 时候配置的 IM Secret，默认是 openIM123
	secret string
	base   string
	client *http.Client
}

func (s *RESTUserService) Sync(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}

func NewRESTUserService(secret string, base string) UserService {
	return &RESTUserService{
		secret: secret,
		base:   base}
}
