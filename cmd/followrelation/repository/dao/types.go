package dao

import "context"

type FollowRelationDao interface {
	// FindFollowerList 获取某人的粉丝列表
	FindFollowerList(ctx context.Context, followee, offset, limit int64) ([]FollowRelation, error)
	// FindFolloweeList 获取某人的关注列表
	FindFolloweeList(ctx context.Context, follower, offset, limit int64) ([]FollowRelation, error)
	FollowRelationDetail(ctx context.Context, follower int64, followee int64) (FollowRelation, error)
	// CreateFollowRelation 创建联系人
	CreateFollowRelation(ctx context.Context, f FollowRelation) error
	// UpdateStatus 更新状态
	UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error
	// CntFollower 统计计算关注自己的人有多少
	CntFollower(ctx context.Context, uid int64) (int64, error)
	// CntFollowee 统计自己关注了多少人
	CntFollowee(ctx context.Context, uid int64) (int64, error)
}
