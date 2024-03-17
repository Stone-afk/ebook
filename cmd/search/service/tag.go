package service

import (
	"context"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/repository"
	"strings"
)

type tagSearchService struct {
	tagRepo repository.TagRepository
}

func (s *tagSearchService) SearchBizTags(ctx context.Context, uid int64, biz string, expression string) (domain.SearchBizTagsResult, error) {
	keywords := strings.Split(expression, " ")
	var res domain.SearchBizTagsResult
	tags, err := s.tagRepo.SearchTag(ctx, uid, biz, keywords)
	res.BizTags = tags
	return res, err
}

func NewTagSearchService(tagRepo repository.TagRepository) TagService {
	return &tagSearchService{tagRepo: tagRepo}
}
