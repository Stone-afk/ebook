package service

import (
	"context"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/repository"
)

type syncService struct {
	tagRepo     repository.TagRepository
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
	anyRepo     repository.AnyRepository
}

func (s *syncService) InputBizTags(ctx context.Context, tag domain.BizTags) error {
	return s.tagRepo.InputBizTags(ctx, tag)
}

func (s *syncService) InputArticle(ctx context.Context, article domain.Article) error {
	return s.articleRepo.InputArticle(ctx, article)
}

func (s *syncService) InputUser(ctx context.Context, user domain.User) error {
	return s.userRepo.InputUser(ctx, user)
}

func (s *syncService) InputAny(ctx context.Context, index, docID, data string) error {
	//cvt := s.converter(index)
	//data = cvt.Convert(data)
	return s.anyRepo.Input(ctx, index, docID, data)
}

func NewSyncService(
	anyRepo repository.AnyRepository,
	userRepo repository.UserRepository,
	articleRepo repository.ArticleRepository) SyncService {
	return &syncService{
		userRepo:    userRepo,
		articleRepo: articleRepo,
		anyRepo:     anyRepo,
	}
}
