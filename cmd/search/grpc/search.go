package grpc

import (
	"context"
	searchv1 "ebook/cmd/api/proto/gen/search/v1"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/service"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

type SearchServiceServer struct {
	searchv1.UnimplementedSearchServiceServer
	searchService service.SearchService
}

func (s *SearchServiceServer) Search(ctx context.Context, request *searchv1.SearchRequest) (*searchv1.SearchResponse, error) {
	resp, err := s.searchService.Search(ctx, request.Uid, request.Expression)
	if err != nil {
		return &searchv1.SearchResponse{}, err
	}
	return &searchv1.SearchResponse{
		User: &searchv1.UserResult{
			Users: slice.Map(resp.Users, func(idx int, src domain.User) *searchv1.User {
				return userConvertToView(src)
			}),
		},
		Article: &searchv1.ArticleResult{
			Articles: slice.Map(resp.Articles, func(idx int, src domain.Article) *searchv1.Article {
				return articleConvertToView(src)
			}),
		},
		BizTags: &searchv1.BizTagsResult{
			Mutibiztags: slice.Map(resp.BizTags, func(idx int, src domain.BizTags) *searchv1.BizTags {
				return bizTagsConvertToView(src)
			}),
		},
	}, nil
}

func NewSearchServiceServer(searchService service.SearchService) *SearchServiceServer {
	return &SearchServiceServer{
		searchService: searchService,
	}
}

func (s *SearchServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterSearchServiceServer(server, s)
}

func bizTagsConvertToView(bizTags domain.BizTags) *searchv1.BizTags {
	return &searchv1.BizTags{
		Uid:   bizTags.Uid,
		Biz:   bizTags.Biz,
		Bizid: bizTags.BizId,
		Tags:  bizTags.Tags,
	}
}

func articleConvertToView(art domain.Article) *searchv1.Article {
	return &searchv1.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  art.Status,
		Content: art.Content,
	}
}

func userConvertToView(u domain.User) *searchv1.User {
	return &searchv1.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Email:    u.Email,
		Phone:    u.Phone,
	}
}
