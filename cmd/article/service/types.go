package service

import (
	"context"
	"ebook/cmd/article/domain"
	"time"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/article/service/types.go -package=svcmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/article/service/mocks/article.mock.go
type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid, id int64) error
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
	List(ctx context.Context, authorId int64, offset, limit int) ([]domain.Article, error)
	// ListPub 只会取 start 七天内的数据
	ListPub(ctx context.Context, uTime time.Time, offset, limit int) ([]domain.Article, error)
	GetPublishedById(ctx context.Context, id int64, userId int64) (domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
}
