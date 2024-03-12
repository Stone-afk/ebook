package service

import (
	"context"
	"ebook/cmd/followrelation/domain"
	"ebook/cmd/followrelation/repository"
)

type followRelationService struct {
	repo repository.FollowRepository
}

func (s *followRelationService) GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error) {
	//TODO implement me
	panic("implement me")
}

func (s *followRelationService) FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error) {
	//TODO implement me
	panic("implement me")
}

func (s *followRelationService) Follow(ctx context.Context, follower, followee int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *followRelationService) CancelFollow(ctx context.Context, follower, followee int64) error {
	//TODO implement me
	panic("implement me")
}

func NewFollowRelationService(repo repository.FollowRepository) FollowRelationService {
	return &followRelationService{
		repo: repo,
	}
}
