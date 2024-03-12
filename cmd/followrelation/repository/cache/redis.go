package cache

import (
	"context"
	"ebook/cmd/followrelation/domain"
	"github.com/redis/go-redis/v9"
)

const (
	// 被多少人关注
	fieldFollowerCnt = "follower_cnt"
	// 关注了多少人
	fieldFolloweeCnt = "followee_cnt"
)

var ErrKeyNotExist = redis.Nil

type RedisFollowCache struct {
	client redis.Cmdable
}

func (r *RedisFollowCache) StaticsInfo(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisFollowCache) SetStaticsInfo(ctx context.Context, uid int64, statics domain.FollowStatics) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisFollowCache) Follow(ctx context.Context, follower, followee int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisFollowCache) CancelFollow(ctx context.Context, follower, followee int64) error {
	//TODO implement me
	panic("implement me")
}

func NewRedisFollowCache(client redis.Cmdable) FollowCache {
	return &RedisFollowCache{
		client: client,
	}
}
