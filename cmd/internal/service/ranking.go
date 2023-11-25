package service

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository"
	"math"
	"time"
)

type RankingService interface {
	// RankTopN 计算 TopN
	RankTopN(ctx context.Context) error
	// TopN 返回业务的 ID
	TopN(ctx context.Context) ([]domain.Article, error)
}

// BatchRankingService 分批计算
type BatchRankingService struct {
	intrSvc InteractiveService
	artSvc  ArticleService
	// 为了测试，不得已暴露出去
	BatchSize int
	N         int
	repo      repository.RankingRepository // 将来扩展，以及支持测试
	scoreFunc func(likeCnt int64, utime time.Time) float64
}

func NewBatchRankingService(
	intrSvc InteractiveService,
	artSvc ArticleService,
	repo repository.RankingRepository) RankingService {
	res := &BatchRankingService{
		intrSvc:   intrSvc,
		artSvc:    artSvc,
		repo:      repo,
		BatchSize: 100,
		N:         100,
	}
	res.scoreFunc = res.score
	return res
}

func (svc *BatchRankingService) RankTopN(ctx context.Context) error {
	arts, err := svc.rankTopN(ctx)
	if err != nil {
		return err
	}
	// 准备放到缓存里面
	return svc.repo.ReplaceTopN(ctx, arts)
}

func (svc *BatchRankingService) rankTopN(ctx context.Context) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (svc *BatchRankingService) TopN(ctx context.Context) ([]domain.Article, error) {
	return svc.repo.GetTopN(ctx)
}

// 这里不需要提前抽象算法，因为正常一家公司的算法都是固定的，不会今天切换到这里，明天切换到那里
func (svc *BatchRankingService) score(likeCnt int64, utime time.Time) float64 {
	// 这个 factor 也可以做成一个参数
	const factor = 1.5
	return float64(likeCnt-1) /
		math.Pow(time.Since(utime).Hours()+2, factor)
}
