package service

import (
	"context"
	"ebook/cmd/followrelation/domain"
)

type FollowRelationService interface {
	GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error)
	GetAllFollower(ctx context.Context, followee int64) ([]domain.FollowRelation, error)
	GetFollower(ctx context.Context, followee, offset, limit int64) ([]domain.FollowRelation, error)
	GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error)
	FollowInfo(ctx context.Context,
		follower, followee int64) (domain.FollowRelation, error)
	Follow(ctx context.Context, follower, followee int64) error
	CancelFollow(ctx context.Context, follower, followee int64) error
}
