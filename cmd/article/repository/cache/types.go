package cache

import (
	"context"
	"ebook/cmd/internal/domain"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/article/repository/cache/types.go -package=cachemocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/cache/mocks/article.mock.go
type ArticleCache interface {
	// GetFirstPage 只缓存第第一页的数据
	// 并且不缓存整个 Content
	GetFirstPage(ctx context.Context, authorId int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, authorId int64, arts []domain.Article) error
	DelFirstPage(ctx context.Context, authorId int64) error

	Set(ctx context.Context, art domain.Article) error
	Get(ctx context.Context, id int64) (domain.Article, error)

	// SetPub 正常来说，创作者和读者的 Redis 集群要分开，因为读者是一个核心中的核心
	SetPub(ctx context.Context, article domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
}
