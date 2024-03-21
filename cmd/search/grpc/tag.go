package grpc

import (
	"context"
	searchv1 "ebook/cmd/api/proto/gen/search/v1"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/service"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

type TagSearchServiceServer struct {
	searchv1.UnimplementedTagSearchServiceServer
	tagService service.TagService
}

func (s *TagSearchServiceServer) SearchBizTags(ctx context.Context, request *searchv1.BizTagsSearchRequest) (*searchv1.BizTagsSearchResponse, error) {
	resp := &searchv1.BizTagsSearchResponse{}
	res, err := s.tagService.SearchBizTags(ctx, request.GetUid(), request.GetBiz(), request.GetExpression())
	if err != nil {
		return resp, err
	}
	resp.Mutibiztags = slice.Map(res.BizTags, func(idx int, src domain.BizTags) *searchv1.BizTags {
		return bizTagsConvertToView(src)
	})
	return resp, nil
}

func NewTagSearchServiceServer(tagService service.TagService) *TagSearchServiceServer {
	return &TagSearchServiceServer{
		tagService: tagService,
	}
}

func (s *TagSearchServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterTagSearchServiceServer(server, s)
}
