package service

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/tag/domain"
	"ebook/cmd/tag/events"
	"ebook/cmd/tag/repository"
)

type tagService struct {
	repo     repository.TagRepository
	logger   logger.Logger
	producer events.Producer
}

func (s *tagService) CreateTag(ctx context.Context, uid int64, name string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *tagService) AttachTags(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *tagService) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (s *tagService) GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func NewTagService(repo repository.TagRepository, l logger.Logger) TagService {
	return &tagService{
		repo:   repo,
		logger: l,
	}
}
