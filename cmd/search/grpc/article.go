package grpc

import (
	"context"
	searchv1 "ebook/cmd/api/proto/gen/search/v1"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/service"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

type ArticleSearchServiceServer struct {
	searchv1.UnimplementedArticleSearchServiceServer
	articleService service.ArticleSearchService
}

func (s *ArticleSearchServiceServer) SearchArticle(ctx context.Context, request *searchv1.ArticleSearchRequest) (*searchv1.ArticleSearchResponse, error) {
	resp := &searchv1.ArticleSearchResponse{}
	res, err := s.articleService.SearchArticle(ctx, request.GetUid(), request.GetExpression())
	if err != nil {
		return resp, err
	}
	resp.Articles = slice.Map(res.Articles, func(idx int, src domain.Article) *searchv1.Article {
		return articleConvertToView(src)
	})
	return resp, nil
}

func NewArticleSearchServiceServer(articleService service.ArticleSearchService) *ArticleSearchServiceServer {
	return &ArticleSearchServiceServer{
		articleService: articleService,
	}
}

func (s *ArticleSearchServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterArticleSearchServiceServer(server, s)
}
