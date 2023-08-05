package cache

import (
	"context"
	"ebook/cmd/internal/domain"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// ErrKeyNotExist 因为我们目前还是只有一个实现，所以可以保持用别名
var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	cmd redis.Cmdable
	// 过期时间
	expiration time.Duration
}

func NewRedisUserCache(cmd redis.Cmdable) *RedisUserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	panic("")
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	panic("")
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
