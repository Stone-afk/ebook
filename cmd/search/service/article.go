package service

import (
	"context"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/repository"
	"strings"
)

type articleSearchService struct {
	articleRepo repository.ArticleRepository
}

func (s *searchService) SearchArticle(ctx context.Context, uid int64, expression string) (domain.SearchArticleResult, error) {
	keywords := strings.Split(expression, " ")
	var res domain.SearchArticleResult
	arts, err := s.articleRepo.SearchArticle(ctx, uid, keywords)
	res.Articles = arts
	return res, err
}

func NewArticleSearchService(articleRepo repository.ArticleRepository) ArticleSearchService {
	return &searchService{articleRepo: articleRepo}
}
