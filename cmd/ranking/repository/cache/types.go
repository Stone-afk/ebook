package cache

import (
	"context"
	"ebook/cmd/ranking/domain"
)

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}
