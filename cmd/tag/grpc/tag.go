package grpc

import (
	"context"
	tagv1 "ebook/cmd/api/proto/gen/tag/v1"
	"ebook/cmd/tag/domain"
	"ebook/cmd/tag/service"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

var _ tagv1.TagServiceServer = (*TagServiceServer)(nil)

type TagServiceServer struct {
	tagv1.UnimplementedTagServiceServer
	service service.TagService
}

func (s *TagServiceServer) CreateTag(ctx context.Context, request *tagv1.CreateTagRequest) (*tagv1.CreateTagResponse, error) {
	id, err := s.service.CreateTag(ctx, request.Uid, request.Name)
	return &tagv1.CreateTagResponse{
		Tag: &tagv1.Tag{
			Id:   id,
			Uid:  request.Uid,
			Name: request.Name,
		},
	}, err
}

func (s *TagServiceServer) AttachTags(ctx context.Context, request *tagv1.AttachTagsRequest) (*tagv1.AttachTagsResponse, error) {
	err := s.service.AttachTags(ctx, request.Uid, request.Biz, request.BizId, request.Tids)
	return &tagv1.AttachTagsResponse{}, err
}

func (s *TagServiceServer) GetTags(ctx context.Context, request *tagv1.GetTagsRequest) (*tagv1.GetTagsResponse, error) {
	tags, err := s.service.GetTags(ctx, request.GetUid())
	if err != nil {
		return nil, err
	}
	return &tagv1.GetTagsResponse{
		Tag: slice.Map(tags, func(idx int, src domain.Tag) *tagv1.Tag {
			return s.toDTO(src)
		}),
	}, nil
}

func (s *TagServiceServer) GetBizTags(ctx context.Context, req *tagv1.GetBizTagsRequest) (*tagv1.GetBizTagsResponse, error) {
	res, err := s.service.GetBizTags(ctx, req.Uid, req.Biz, req.BizId)
	if err != nil {
		return nil, err
	}
	return &tagv1.GetBizTagsResponse{
		Tags: slice.Map(res, func(idx int, src domain.Tag) *tagv1.Tag {
			return s.toDTO(src)
		}),
	}, nil
}

func (s *TagServiceServer) toDTO(tag domain.Tag) *tagv1.Tag {
	return &tagv1.Tag{
		Id:   tag.Id,
		Uid:  tag.Uid,
		Name: tag.Name,
	}
}

func NewTagServiceServer(svc service.TagService) *TagServiceServer {
	return &TagServiceServer{
		service: svc,
	}
}

func (s *TagServiceServer) Register(server grpc.ServiceRegistrar) {
	tagv1.RegisterTagServiceServer(server, s)
}
