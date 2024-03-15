package cache

import (
	"context"
	"ebook/cmd/tag/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisTagCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func (r *RedisTagCache) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisTagCache) Append(ctx context.Context, uid int64, tags ...domain.Tag) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisTagCache) DelTags(ctx context.Context, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func NewRedisTagCache(client redis.Cmdable) TagCache {
	return &RedisTagCache{
		client: client,
	}
}
