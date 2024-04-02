package repository

import (
	"context"
	"ebook/cmd/followrelation/domain"
)

type FollowRepository interface {
	// GetFollower 获取某人的粉丝列表
	GetFollower(ctx context.Context, followee, offset, limit int64) ([]domain.FollowRelation, error)
	// GetFollowee 获取某人的关注列表
	GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error)
	// FollowInfo 查看关注人的详情
	FollowInfo(ctx context.Context, follower int64, followee int64) (domain.FollowRelation, error)
	// AddFollowRelation 创建关注关系
	AddFollowRelation(ctx context.Context, f domain.FollowRelation) error
	// InactiveFollowRelation 取消关注
	InactiveFollowRelation(ctx context.Context, follower int64, followee int64) error
	GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error)
}
