package cache

import (
	"context"
	"ebook/cmd/ranking/domain"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRankingCache struct {
	key        string
	client     redis.Cmdable
	expiration time.Duration
}

func NewRedisRankingCache(client redis.Cmdable) *RedisRankingCache {
	return &RedisRankingCache{
		key:        "ranking:article",
		client:     client,
		expiration: time.Minute * 3,
	}
}

func (r *RedisRankingCache) Set(ctx context.Context, arts []domain.Article) error {
	// 这里不会缓存内容
	for i := 0; i < len(arts); i++ {
		arts[i].Content = arts[i].Abstract()
	}
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	// 过期时间要设置得比定时计算的间隔长
	return r.client.Set(ctx, r.key, val, r.expiration).Err()
}

func (r *RedisRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	val, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal(val, &res)
	return nil, err
}
