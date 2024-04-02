package grpc

import (
	"context"
	followv1 "ebook/cmd/api/proto/gen/followrelation/v1"
	"ebook/cmd/followrelation/domain"
	"ebook/cmd/followrelation/service"
	"google.golang.org/grpc"
)

type FollowServiceServer struct {
	followv1.UnimplementedFollowServiceServer
	svc service.FollowRelationService
}

func (s *FollowServiceServer) convertToView(relation domain.FollowRelation) *followv1.FollowRelation {
	return &followv1.FollowRelation{
		Followee: relation.Followee,
		Follower: relation.Follower,
	}
}

func (s *FollowServiceServer) GetFollower(ctx context.Context, request *followv1.GetFollowerRequest) (*followv1.GetFollowerResponse, error) {
	// TODO 先用 GetAllFollower
	relationList, err := s.svc.GetAllFollower(ctx, request.Followee)
	if err != nil {
		return nil, err
	}
	res := make([]*followv1.FollowRelation, 0, len(relationList))
	for _, relation := range relationList {
		res = append(res, s.convertToView(relation))
	}
	return &followv1.GetFollowerResponse{
		FollowRelations: res,
	}, nil
}

func (s *FollowServiceServer) GetFollowStatic(ctx context.Context, request *followv1.GetFollowStaticRequest) (*followv1.GetFollowStaticResponse, error) {
	res, err := s.svc.GetFollowStatics(ctx, request.Followee)
	if err != nil {
		return nil, err
	}
	return &followv1.GetFollowStaticResponse{
		FollowStatic: &followv1.FollowStatic{
			Followers: res.Followers,
			Followees: res.Followees,
		},
	}, nil
}

func (s *FollowServiceServer) GetFollowee(ctx context.Context, request *followv1.GetFolloweeRequest) (*followv1.GetFolloweeResponse, error) {
	relationList, err := s.svc.GetFollowee(ctx, request.Follower, request.Offset, request.Limit)
	if err != nil {
		return nil, err
	}
	res := make([]*followv1.FollowRelation, 0, len(relationList))
	for _, relation := range relationList {
		res = append(res, s.convertToView(relation))
	}
	return &followv1.GetFolloweeResponse{
		FollowRelations: res,
	}, nil
}

func (s *FollowServiceServer) FollowInfo(ctx context.Context, request *followv1.FollowInfoRequest) (*followv1.FollowInfoResponse, error) {
	info, err := s.svc.FollowInfo(ctx, request.Follower, request.Followee)
	if err != nil {
		return nil, err
	}
	return &followv1.FollowInfoResponse{
		FollowRelation: s.convertToView(info),
	}, nil
}

func (s *FollowServiceServer) Follow(ctx context.Context, request *followv1.FollowRequest) (*followv1.FollowResponse, error) {
	err := s.svc.Follow(ctx, request.Follower, request.Followee)
	return &followv1.FollowResponse{}, err
}

func (s *FollowServiceServer) CancelFollow(ctx context.Context, request *followv1.CancelFollowRequest) (*followv1.CancelFollowResponse, error) {
	err := s.svc.CancelFollow(ctx, request.Follower, request.Followee)
	return &followv1.CancelFollowResponse{}, err
}

func (s *FollowServiceServer) Register(server grpc.ServiceRegistrar) {
	followv1.RegisterFollowServiceServer(server, s)
}

func NewFollowRelationServiceServer(svc service.FollowRelationService) *FollowServiceServer {
	return &FollowServiceServer{
		svc: svc,
	}
}
