package grpc

import (
	"context"
	searchv1 "ebook/cmd/api/proto/gen/search/v1"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/service"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

type UserSearchServiceServer struct {
	searchv1.UnimplementedUserSearchServiceServer
	userService service.UserSearchService
}

func (s *UserSearchServiceServer) SearchUser(ctx context.Context, request *searchv1.UserSearchRequest) (*searchv1.UserSearchResponse, error) {
	resp := &searchv1.UserSearchResponse{}
	res, err := s.userService.SearchUser(ctx, request.GetExpression())
	if err != nil {
		return resp, err
	}
	resp.Users = slice.Map(res.Users, func(idx int, src domain.User) *searchv1.User {
		return userConvertToView(src)
	})
	return resp, nil
}

func NewUserSearchServiceServer(userService service.UserSearchService) *UserSearchServiceServer {
	return &UserSearchServiceServer{
		userService: userService,
	}
}

func (s *UserSearchServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterUserSearchServiceServer(server, s)
}
