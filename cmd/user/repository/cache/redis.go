package cache

import (
	"context"
	"ebook/cmd/user/domain"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var _ UserCache = &RedisUserCache{}

type RedisUserCache struct {
	cmd redis.Cmdable
	// 过期时间
	expiration time.Duration
}

func NewRedisUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (c *RedisUserCache) Delete(ctx context.Context, id int64) error {
	return c.cmd.Del(ctx, c.key(id)).Err()
}

func (c *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := c.key(id)
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	// 反序列化回来
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}

func (c *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, c.key(u.Id), data, c.expiration).Err()
}

func (c *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
