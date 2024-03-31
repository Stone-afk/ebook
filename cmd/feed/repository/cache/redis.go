package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type feedEventCache struct {
	client redis.Cmdable
}

func NewFeedEventCache(client redis.Cmdable) FeedEventCache {
	return &feedEventCache{
		client: client,
	}
}

func (c *feedEventCache) SetFollowees(ctx context.Context, follower int64, followees []int64) error {
	key := c.getFolloweeKey(follower)
	followeesStr, err := json.Marshal(followees)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, followeesStr, FolloweeKeyExpiration).Err()
}

func (c *feedEventCache) GetFollowees(ctx context.Context, follower int64) ([]int64, error) {
	key := c.getFolloweeKey(follower)
	res, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, FolloweesNotFound
	}
	var followees []int64
	err = json.Unmarshal([]byte(res), &followees)
	if err != nil {
		return nil, err
	}
	return followees, nil
}

func (c *feedEventCache) getFolloweeKey(follower int64) string {
	return fmt.Sprintf("feed_event:%d", follower)
}
