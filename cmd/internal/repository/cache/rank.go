package cache

import (
	"context"
	"ebook/cmd/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}

type RedisRankingCache struct {
	client     redis.Cmdable
	key        string
	expiration time.Duration
}

func (r *RedisRankingCache) Set(ctx context.Context, arts []domain.Article) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}
