package grpc

import (
	"context"
	oauth2v1 "ebook/cmd/api/proto/gen/oauth2/v1"
	"ebook/cmd/oauth2/service/wechat"
	"google.golang.org/grpc"
)

type Oauth2ServiceServer struct {
	oauth2v1.UnimplementedOauth2ServiceServer
	service wechat.Service
}

func NewOauth2ServiceServer(svc wechat.Service) *Oauth2ServiceServer {
	return &Oauth2ServiceServer{
		service: svc,
	}
}

func (s *Oauth2ServiceServer) Register(server grpc.ServiceRegistrar) {
	oauth2v1.RegisterOauth2ServiceServer(server, s)
}

func (s *Oauth2ServiceServer) AuthURL(ctx context.Context, req *oauth2v1.AuthURLRequest) (*oauth2v1.AuthURLResponse, error) {
	url, err := s.service.AuthURL(ctx, req.State)
	if err != nil {
		return nil, err
	}
	return &oauth2v1.AuthURLResponse{
		Url: url,
	}, nil
}

func (s *Oauth2ServiceServer) VerifyCode(ctx context.Context, req *oauth2v1.VerifyCodeRequest) (*oauth2v1.VerifyCodeResponse, error) {
	info, err := s.service.VerifyCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &oauth2v1.VerifyCodeResponse{
		OpenId:  info.OpenId,
		UnionId: info.UnionId,
	}, nil
}
