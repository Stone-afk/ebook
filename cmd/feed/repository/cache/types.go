package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var FolloweesNotFound = redis.Nil

const FolloweeKeyExpiration = 10 * time.Minute

type FeedEventCache interface {
	// SetFollowees follower: 关注者id， followees: 被关注者id列表
	SetFollowees(ctx context.Context, follower int64, followees []int64) error
	GetFollowees(ctx context.Context, follower int64) ([]int64, error)
}
