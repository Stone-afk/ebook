package service

import (
	"context"
	"ebook/cmd/comment/domain"
	"ebook/cmd/comment/repository"
)

type commentService struct {
	repo repository.CommentRepository
}

func (s *commentService) GetCommentList(ctx context.Context, biz string, bizId, minID, limit int64) ([]domain.Comment, error) {
	list, err := s.repo.FindByBiz(ctx, biz, bizId, minID, limit)
	if err != nil {
		return nil, err
	}
	return list, err
}

func (s *commentService) DeleteComment(ctx context.Context, id int64) error {
	return s.repo.DeleteComment(ctx, domain.Comment{Id: id})
}

func (s *commentService) CreateComment(ctx context.Context, comment domain.Comment) error {
	return s.repo.CreateComment(ctx, comment)
}

func (s *commentService) GetMoreReplies(ctx context.Context, rid int64, maxID int64, limit int64) ([]domain.Comment, error) {
	return s.repo.GetMoreReplies(ctx, rid, maxID, limit)
}

func NewCommentSvc(repo repository.CommentRepository) CommentService {
	return &commentService{
		repo: repo,
	}
}
