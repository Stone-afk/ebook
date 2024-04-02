package service

import (
	"context"
	"ebook/cmd/followrelation/domain"
	"ebook/cmd/followrelation/repository"
)

type followRelationService struct {
	repo repository.FollowRepository
}

func (s *followRelationService) GetAllFollower(ctx context.Context, followee int64) ([]domain.FollowRelation, error) {
	offset := 0
	limit := 1000
	relations := make([]domain.FollowRelation, 0, limit)
	for {
		res, err := s.repo.GetFollower(ctx, followee, int64(offset), int64(limit))
		if err != nil {
			return nil, err
		}
		relations = append(relations, res...)
		if len(res) < limit {
			break
		}
		offset += limit
	}
	return relations, nil
}

func (s *followRelationService) GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	return s.repo.GetFollowStatics(ctx, uid)
}

func (s *followRelationService) GetFollower(ctx context.Context, followee, offset, limit int64) ([]domain.FollowRelation, error) {
	return s.repo.GetFollower(ctx, followee, offset, limit)
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
