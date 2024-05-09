package service

import (
	"context"
	"ebook/cmd/ranking/domain"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/ranking/service/types.go -package=svcmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/ranking/service/mocks/ranking.mock.go
type RankingService interface {
	// RankTopN 计算 TopN
	RankTopN(ctx context.Context) error
	// TopN 返回业务的 ID
	TopN(ctx context.Context) ([]domain.Article, error)
}
