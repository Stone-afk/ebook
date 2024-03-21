package grpc

import (
	"context"
	searchv1 "ebook/cmd/api/proto/gen/search/v1"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/service"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchServiceServer struct {
	searchv1.UnimplementedSearchServiceServer
	searchv1.UnimplementedArticleSearchServiceServer
	searchv1.UnimplementedTagSearchServiceServer
	searchv1.UnimplementedUserSearchServiceServer
	searchService  service.SearchService
	userService    service.UserSearchService
	articleService service.ArticleSearchService
	tagService     service.TagService
}

func (s *SearchServiceServer) SearchUser(ctx context.Context, request *searchv1.UserSearchRequest) (*searchv1.UserSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchUser not implemented")
}

func (s *SearchServiceServer) SearchBizTags(ctx context.Context, request *searchv1.BizTagsSearchRequest) (*searchv1.BizTagsSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchBizTags not implemented")
}

func (s *SearchServiceServer) SearchArticle(ctx context.Context, request *searchv1.ArticleSearchRequest) (*searchv1.ArticleSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchArticle not implemented")
}

func (s *SearchServiceServer) Search(ctx context.Context, request *searchv1.SearchRequest) (*searchv1.SearchResponse, error) {
	resp, err := s.searchService.Search(ctx, request.Uid, request.Expression)
	if err != nil {
		return nil, err
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
	}, nil
}

func NewSearchService(searchService service.SearchService,
	userService service.UserSearchService,
	articleService service.ArticleSearchService,
	tagService service.TagService) *SearchServiceServer {
	return &SearchServiceServer{
		searchService:  searchService,
		userService:    userService,
		articleService: articleService,
		tagService:     tagService,
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
