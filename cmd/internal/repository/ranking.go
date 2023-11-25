package repository

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository/cache"
	"github.com/ecodeclub/ekit/syncx/atomicx"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type rankingRepository struct {
	localCache *cache.RankingLocalCache
	redisCache *cache.RedisRankingCache
	// 也可以考虑将这个本地缓存塞进去 RankingCache 里面，作为一个实现
	topN atomicx.Value[[]domain.Article]
}

func NewCachedRankingRepository(redisCache *cache.RedisRankingCache,
	localCache *cache.RankingLocalCache) RankingRepository {
	return &rankingRepository{
		localCache: localCache,
		redisCache: redisCache,
	}
}

func (repo *rankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	return repo.redisCache.Set(ctx, arts)
}

func (repo *rankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}
