package service

import (
	"context"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/repository"
)

type tagSearchService struct {
	tagRepo repository.TagRepository
}

func (t tagSearchService) Search(ctx context.Context, uid int64, biz string, expression string) (domain.SearchResult, error) {
	//TODO implement me
	panic("implement me")
}

func NewTagSearchService(tagRepo repository.TagRepository) TagService {
	return &tagSearchService{tagRepo: tagRepo}
}
