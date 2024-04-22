package grpc

import (
	"context"
	codev1 "ebook/cmd/api/proto/gen/code/v1"
	"ebook/cmd/code/service"
	"google.golang.org/grpc"
)

type CodeServiceServer struct {
	codev1.UnimplementedCodeServiceServer
	service service.CodeService
}

func NewCodeServiceServer(svc service.CodeService) *CodeServiceServer {
	return &CodeServiceServer{
		service: svc,
	}
}
func (s *CodeServiceServer) Register(server grpc.ServiceRegistrar) {
	codev1.RegisterCodeServiceServer(server, s)
}

func (s *CodeServiceServer) Send(ctx context.Context, req *codev1.CodeSendRequest) (*codev1.CodeSendResponse, error) {
	err := s.service.Send(ctx, req.Biz, req.Phone)
	return &codev1.CodeSendResponse{}, err
}

func (s *CodeServiceServer) Verify(ctx context.Context, req *codev1.VerifyRequest) (*codev1.VerifyResponse, error) {
	ans, err := s.service.Verify(ctx, req.Biz, req.Phone, req.InputCode)
	return &codev1.VerifyResponse{
		Answer: ans,
	}, err
}
