package repository

import (
	"context"
	"ebook/cmd/article/domain"
	"time"
)

// repository 还是要用来操作缓存和DAO
// 事务概念应该在 DAO 这一层

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/article/repository/types.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/article/repository/mocks/article.mock.go
type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	// Sync 本身要求先保存到制作库，再同步到线上库
	Sync(ctx context.Context, art domain.Article) (int64, error)
	// SyncStatus 仅仅同步状态
	SyncStatus(ctx context.Context, uid, id int64, status domain.ArticleStatus) error
	List(ctx context.Context, author int64, offset, limit int) ([]domain.Article, error)
	ListPub(ctx context.Context, uTime time.Time, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
}

// AuthorRepository 封装user的client用于获取用户信息
type AuthorRepository interface {
	// FindAuthor id为文章id
	FindAuthor(ctx context.Context, id int64) (domain.Author, error)
}

// HistoryRecordRepository 也就是一个增删改查的事情
type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, r domain.HistoryRecord) error
}

// ArticleAuthorRepository 演示在 service 层面上分流
type ArticleAuthorRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}

type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) error
}
