package repository

import (
	"context"
	"ebook/cmd/ranking/domain"
	"ebook/cmd/ranking/repository/cache"
	"github.com/ecodeclub/ekit/syncx/atomicx"
)

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
	// 这一步必然不会出错
	_ = repo.localCache.Set(ctx, arts)
	return repo.redisCache.Set(ctx, arts)
}

func (repo *rankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	arts, err := repo.localCache.Get(ctx)
	if err == nil {
		return arts, nil
	}
	arts, err = repo.redisCache.Get(ctx)
	if err != nil {
		// 这里，没有进一步区分是什么原因导致的 Redis 错误
		return repo.localCache.ForceGet(ctx)
	}
	// 回写本地缓存
	_ = repo.localCache.Set(ctx, arts)
	return arts, err
}
