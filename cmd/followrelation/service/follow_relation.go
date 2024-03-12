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
	return s.repo.GetFollowee(ctx, follower, offset, limit)
}

func (s *followRelationService) FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error) {
	return s.repo.FollowInfo(ctx, follower, followee)
}

func (s *followRelationService) Follow(ctx context.Context, follower, followee int64) error {
	return s.repo.AddFollowRelation(ctx, domain.FollowRelation{
		Followee: followee,
		Follower: follower,
	})
}

func (s *followRelationService) CancelFollow(ctx context.Context, follower, followee int64) error {
	return s.repo.InactiveFollowRelation(ctx, follower, followee)
}

func NewFollowRelationService(repo repository.FollowRepository) FollowRelationService {
	return &followRelationService{
		repo: repo,
	}
}
